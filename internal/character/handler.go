package character

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
)

type Handler struct {
	characterService *Service
}

func NewHandler(characterService *Service) *Handler {
	return &Handler{characterService: characterService}
}

// CreateCharacter godoc
// @Summary      Criar um novo personagem
// @Description  Cria um novo personagem com os dados informados no corpo da requisição.
// @Tags         Character
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        character  body      CreateInput  true  "Dados do Personagem"
// @Success      201        {object}  Response
// @Router       /characters [post]
func (h *Handler) CreateCharacter(c *gin.Context) {
	httpx.LogRequest(c, "Character", "CreateCharacter")

	var input CreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
		return
	}

	character := models.Character{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Type:        input.Type,
	}

	if err := h.characterService.CreateCharacter(c.Request.Context(), &character); err != nil {
		if errors.Is(err, ErrExists) {
			c.JSON(http.StatusConflict, httpx.ErrorResponse{Error: err.Error()})
			return
		}
		httpx.InternalError(c, "Character", "CreateCharacter", err)
		return
	}

	c.JSON(http.StatusCreated, toCharacterResponse(&character))
}

// GetCharacterById godoc
// @Summary      Buscar personagem por ID
// @Tags         Character
// @Produce      json
// @Param        id   path      int  true  "ID do personagem" minimum(1)
// @Success      200  {object}  Response
// @Router       /characters/{id} [get]
func (h *Handler) GetCharacterById(c *gin.Context) {
	httpx.LogRequest(c, "Character", "GetCharacterById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	character, err := h.characterService.GetCharacterById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "personagem não encontrado"})
			return
		}
		httpx.InternalError(c, "Character", "GetCharacterById", err)
		return
	}

	c.JSON(http.StatusOK, toCharacterResponse(character))
}

// DeleteCharacterById godoc
// @Summary      Deletar personagem por ID
// @Tags         Character
// @Produce      json
// @Security     BearerAuth
// @Param        id   path        int  true  "ID do personagem" minimum(1)
// @Success      200  {object}    map[string]string
// @Router       /characters/{id} [delete]
func (h *Handler) DeleteCharacterById(c *gin.Context) {
	httpx.LogRequest(c, "Character", "DeleteCharacterById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.characterService.DeleteCharacterById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "personagem não encontrado"})
			return
		}
		if errors.Is(err, ErrHasGuides) {
			c.JSON(http.StatusConflict, httpx.ErrorResponse{Error: "personagem possui guias"})
			return
		}
		httpx.InternalError(c, "Character", "DeleteCharacterById", err)
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
// @Param        body  body      UpdateInput true  "Dados para atualização"
// @Success      200   {object}  Response
// @Router       /characters/{id} [patch]
func (h *Handler) UpdateCharacter(c *gin.Context) {
	httpx.LogRequest(c, "Character", "UpdateCharacter")

	var uriInput IDInput
	var bodyInput UpdateInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "corpo da requisição inválido"})
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
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "personagem não encontrado"})
		case errors.Is(err, apperr.ErrNoUpdateData):
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
		case errors.Is(err, ErrExists):
			c.JSON(http.StatusConflict, httpx.ErrorResponse{Error: err.Error()})
		default:
			httpx.InternalError(c, "Character", "UpdateCharacter", err)
		}
		return
	}

	c.JSON(http.StatusOK, toCharacterResponse(updatedCharacter))
}

// GetCharacters godoc
// @Summary      Listar personagens
// @Description  Lista todos os personagens. Nome (busca parcial) e tipo são filtros opcionais.
// @Tags         Character
// @Produce      json
// @Param        name  query     string  false  "Filtrar por nome"
// @Param        type  query     string  false  "Filtrar por tipo"
// @Success      200   {array}   Response
// @Failure      400   {object}  httpx.ErrorResponse
// @Router       /characters [get]
func (h *Handler) GetCharacters(c *gin.Context) {
	httpx.LogRequest(c, "Character", "GetCharacters")

	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "parâmetros de consulta inválidos"})
		return
	}

	characters, err := h.characterService.GetCharacters(
		c.Request.Context(),
		strings.TrimSpace(query.Name),
		strings.TrimSpace(query.Type),
	)
	if err != nil {
		httpx.InternalError(c, "Character", "GetCharacters", err)
		return
	}

	response := make([]Response, 0, len(characters))
	for i := range characters {
		response = append(response, toCharacterResponse(&characters[i]))
	}

	c.JSON(http.StatusOK, response)
}

func toCharacterResponse(character *models.Character) Response {
	return Response{
		CharacterID: character.CharacterID,
		Name:        character.Name,
		Description: character.Description,
		Type:        character.Type,
		ImageURL:    character.ImageURL,
	}
}
