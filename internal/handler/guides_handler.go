package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dto "github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/middleware"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

type GuideHandler struct {
	guideService *services.GuideService
}

func NewGuideHandler(guideService *services.GuideService) *GuideHandler {
	return &GuideHandler{guideService: guideService}
}

// CreateGuide godoc
// @Summary      Criar um novo guia
// @Description  Cria um guia com os dados informados no corpo da requisição.
// @Tags         Guide
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        guide  body      dto.CreateGuideInput  true  "Dados do guia"
// @Success      201    {object}  dto.GuideResponse
// @Router       /guides [post]
func (h *GuideHandler) CreateGuide(c *gin.Context) {
	logRequest(c, "Guide", "CreateGuide")

	var input dto.CreateGuideInput
	claims := middleware.GetClaims(c)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	guide := models.Guide{
		Title:       input.Title,
		BannerURL:   input.BannerURL,
		Type:        input.Type,
		Description: input.Description,
		Link:        input.Link,
		CharacterID: input.CharacterID,
		UserID:      uint(claims.UserID), // Supondo que você tenha um middleware que adiciona o ID do usuário ao contexto
	}

	if err := h.guideService.CreateGuide(c.Request.Context(), &guide); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "personagem não existe"})
			return
		}
		internalError(c, "Guide", "CreateGuide", err)
		return
	}

	c.JSON(http.StatusCreated, toGuideResponse(&guide))
}

// GetGuideById godoc
// @Summary      Buscar guia por ID
// @Tags         Guide
// @Produce      json
// @Param        id   path      int  true  "ID do guia" minimum(1)
// @Success      200  {object}  dto.GuideResponse
// @Router       /guides/{id} [get]
func (h *GuideHandler) GetGuideById(c *gin.Context) {
	logRequest(c, "Guide", "GetGuideById")

	var input dto.GuideIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	guide, err := h.guideService.GetGuideById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		internalError(c, "Guide", "GetGuideById", err)
		return
	}

	c.JSON(http.StatusOK, toGuideResponse(guide))
}

// DeleteGuideById godoc
// @Summary      Deletar guia por ID
// @Tags         Guide
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do guia" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /guides/{id} [delete]
func (h *GuideHandler) DeleteGuideById(c *gin.Context) {
	logRequest(c, "Guide", "DeleteGuideById")

	var input dto.GuideIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.guideService.DeleteGuideById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		internalError(c, "Guide", "DeleteGuideById", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guia deletado com sucesso"})
}

// UpdateGuide godoc
// @Summary      Atualizar dados do guia
// @Tags         Guide
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                   true  "ID do guia" minimum(1)
// @Param        body  body      dto.UpdateGuideInput  true  "Dados para atualização"
// @Success      200   {object}  dto.GuideResponse
// @Router       /guides/{id} [patch]
func (h *GuideHandler) UpdateGuide(c *gin.Context) {
	logRequest(c, "Guide", "UpdateGuide")

	var uriInput dto.GuideIDInput
	var bodyInput dto.UpdateGuideInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Não autenticado"})
		return
	}
	if models.Role(claims.Role) != models.RoleAdmin && uint(claims.UserID) != uint(uriInput.ID) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Você so pode editar seu proprio guia"})
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	guideModel := models.Guide{
		Title:       bodyInput.Title,
		BannerURL:   bodyInput.BannerURL,
		Type:        bodyInput.Type,
		Description: bodyInput.Description,
		Link:        bodyInput.Link,
	}

	updatedGuide, err := h.guideService.UpdateGuide(c.Request.Context(), uriInput.ID, &guideModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		if err.Error() == "nenhum dado recebido para atualização" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		internalError(c, "Guide", "UpdateGuide", err)
		return
	}

	c.JSON(http.StatusOK, toGuideResponse(updatedGuide))
}

// GetGuides godoc
// @Summary      Listar guias
// @Description  Lista guias com paginação. O título é um filtro opcional (busca parcial).
// @Tags         Guide
// @Produce      json
// @Param        title  query     string  false  "Filtrar por título"
// @Param        page   query     int     false  "Página"           default(1)  minimum(1)
// @Param        limit  query     int     false  "Itens por página"  default(20) minimum(1) maximum(100)
// @Success      200    {object}  dto.PageResponse[dto.GuideResponse]
// @Failure      400    {object}  dto.ErrorResponse
// @Router       /guides [get]
func (h *GuideHandler) GetGuides(c *gin.Context) {
	logRequest(c, "Guide", "GetGuides")

	var query dto.GuideListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "parâmetros inválidos"})
		return
	}

	guides, total, err := h.guideService.GetGuides(c.Request.Context(), strings.TrimSpace(query.Title), query.Page, query.Limit)
	if err != nil {
		internalError(c, "Guide", "GetGuides", err)
		return
	}

	data := make([]dto.GuideResponse, 0, len(guides))
	for i := range guides {
		data = append(data, toGuideResponse(&guides[i]))
	}

	c.JSON(http.StatusOK, dto.PageResponse[dto.GuideResponse]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	})
}

func toGuideResponse(g *models.Guide) dto.GuideResponse {
	return dto.GuideResponse{
		GuideID:      g.GuideID,
		Title:        g.Title,
		BannerURL:    g.BannerURL,
		Type:         g.Type,
		Description:  g.Description,
		Link:         g.Link,
		CreationDate: g.CreatedAt,
		Likes:        g.Likes,
		Dislikes:     g.Dislikes,
		UserID:       g.UserID,
		CharacterID:  g.CharacterID,
	}
}
