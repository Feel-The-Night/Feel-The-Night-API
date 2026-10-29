package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

type CharacterHandler struct {
	characterService *services.CharacterService
}

func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

// CreateCharacter godoc
// @Summary      Criar um novo personagem
// @Description  Cria um novo personagem com os dados informados no corpo da requisição.
// @Tags         Character
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        character  body      dto.CreateCharacterInput  true  "Dados do Personagem"
// @Success      201        {object}  dto.CharacterResponse
// @Router       /character [post]
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	logRequest(c, "Character", "CreateCharacter")

	var input dto.CreateCharacterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	character := models.Character{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Type:        input.Type,
	}

	if err := h.characterService.CreateCharacter(c.Request.Context(), &character); err != nil {
		if errors.Is(err, services.ErrCharacterExists) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		internalError(c, "Character", "CreateCharacter", err)
		return
	}

	c.JSON(http.StatusCreated, toCharacterResponse(&character))
}

// GetCharacterById godoc
// @Summary      Buscar personagem por ID
// @Tags         Character
// @Produce      json
// @Param        id   path      int  true  "ID do personagem" minimum(1)
// @Success      200  {object}  dto.CharacterResponse
// @Router       /character/{id} [get]
func (h *CharacterHandler) GetCharacterById(c *gin.Context) {
	logRequest(c, "Character", "GetCharacterById")

	var input dto.CharacterIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	character, err := h.characterService.GetCharacterById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "personagem não encontrado"})
			return
		}
		internalError(c, "Character", "GetCharacterById", err)
		return
	}

	c.JSON(http.StatusOK, toCharacterResponse(character))
}

// DeleteCharacterById godoc
// @Summary      Deletar personagem por ID
// @Tags         Character
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do personagem" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /character/{id} [delete]
func (h *CharacterHandler) DeleteCharacterById(c *gin.Context) {
	logRequest(c, "Character", "DeleteCharacterById")

	var input dto.CharacterIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.characterService.DeleteCharacterById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "personagem não encontrado"})
			return
		}
		internalError(c, "Character", "DeleteCharacterById", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "personagem deletado com sucesso"})
}

// UpdateCharacter godoc
// @Summary      Atualizar dados do personagem
// @Tags         Character
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                      true  "ID do personagem" minimum(1)
// @Param        body  body      dto.UpdateCharacterInput true  "Dados para atualização"
// @Success      200   {object}  dto.CharacterResponse
// @Router       /character/{id} [patch]
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	logRequest(c, "Character", "UpdateCharacter")

	var uriInput dto.CharacterIDInput
	var bodyInput dto.UpdateCharacterInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	characterModel := models.Character{
		Name:        bodyInput.Name,
		Description: bodyInput.Description,
		Type:        bodyInput.Type,
		ImageURL:    bodyInput.ImageURL,
	}

	updatedCharacter, err := h.characterService.UpdateCharacterById(c.Request.Context(), uriInput.ID, &characterModel)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "personagem não encontrado"})
		case errors.Is(err, services.ErrNoUpdateData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, services.ErrCharacterExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		default:
			internalError(c, "Character", "UpdateCharacter", err)
		}
		return
	}

	c.JSON(http.StatusOK, toCharacterResponse(updatedCharacter))
}

// GetCharacters godoc
// @Summary      Listar personagens
// @Description  Lista personagens com paginação. O Nome é um filtro opcional (busca parcial).
// @Tags         Character
// @Produce      json
// @Param        name   query     string  false  "Filtrar por Nome"
// @Success      200    {object}  dto.CharacterResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Router       /character [get]
func (h *CharacterHandler) GetCharacters(c *gin.Context) {
	logRequest(c, "Character", "GetCharacters")

	var query dto.CharacterListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "parâmetros de consulta inválidos"})
		return
	}

	characters, total, err := h.characterService.GetCharacters(c.Request.Context(), strings.TrimSpace(query.Name), query.Page, query.Limit)
	if err != nil {
		internalError(c, "Character", "GetCharacters", err)
		return
	}

	data := make([]dto.CharacterResponse, len(characters))
	for i, character := range characters {
		data[i] = toCharacterResponse(&character)
	}

	// c.JSON(http.StatusOK, response)
	c.JSON(http.StatusOK, dto.PageResponse[dto.CharacterResponse]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	})
}

func toCharacterResponse(character *models.Character) dto.CharacterResponse {
	return dto.CharacterResponse{
		CharacterID: character.CharacterID,
		Name:        character.Name,
		Description: character.Description,
		Type:        character.Type,
		ImageURL:    character.ImageURL,
	}
}
