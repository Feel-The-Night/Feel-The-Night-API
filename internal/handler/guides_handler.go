package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	dto "github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
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
// @Param        guide  body      dto.CreateGuideInput  true  "Dados do guia"
// @Success      201    {object}  dto.GuideResponse
// @Router       /guides [post]
func (h *GuideHandler) CreateGuide(c *gin.Context) {
	zap.L().Info("[GuideHandler] CreateGuide",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CreateGuideInput

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
		UserID:      input.UserID,
		CharacterID: input.CharacterID,
	}

	if err := h.guideService.CreateGuide(&guide); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.GuideResponse{
		GuideID:      guide.GuideID,
		Title:        guide.Title,
		BannerURL:    guide.BannerURL,
		Type:         guide.Type,
		Description:  guide.Description,
		Link:         guide.Link,
		CreationDate: guide.CreatedAt,
		Likes:        guide.Likes,
		Dislikes:     guide.Dislikes,
		UserID:       guide.UserID,
		CharacterID:  guide.CharacterID,
	}

	c.JSON(http.StatusCreated, response)
}

// GetGuideById godoc
// @Summary      Buscar guia por ID
// @Tags         Guide
// @Produce      json
// @Param        id   path      int  true  "ID do guia" minimum(1)
// @Success      200  {object}  dto.GuideResponse
// @Router       /guides/{id} [get]
func (h *GuideHandler) GetGuideById(c *gin.Context) {
	zap.L().Info("[GuideHandler] GetGuideById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.GuideIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	guide, err := h.guideService.GetGuideById(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.GuideResponse{
		GuideID:      guide.GuideID,
		Title:        guide.Title,
		BannerURL:    guide.BannerURL,
		Type:         guide.Type,
		Description:  guide.Description,
		Link:         guide.Link,
		CreationDate: guide.CreatedAt,
		Likes:        guide.Likes,
		Dislikes:     guide.Dislikes,
		UserID:       guide.UserID,
		CharacterID:  guide.CharacterID,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteGuideById godoc
// @Summary      Deletar guia por ID
// @Tags         Guide
// @Produce      json
// @Param        id   path      int  true  "ID do guia" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /guides/{id} [delete]
func (h *GuideHandler) DeleteGuideById(c *gin.Context) {
	zap.L().Info("[GuideHandler] DeleteGuideById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.GuideIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.guideService.DeleteGuideById(input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guia deletado com sucesso"})
}

// UpdateGuide godoc
// @Summary      Atualizar dados do guia
// @Tags         Guide
// @Accept       json
// @Produce      json
// @Param        id    path      int                   true  "ID do guia" minimum(1)
// @Param        body  body      dto.UpdateGuideInput  true  "Dados para atualização"
// @Success      200   {object}  dto.GuideResponse
// @Router       /guides/{id} [patch]
func (h *GuideHandler) UpdateGuide(c *gin.Context) {
	zap.L().Info("[GuideHandler] UpdateGuide",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var uriInput dto.GuideIDInput
	var bodyInput dto.UpdateGuideInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
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

	updatedGuide, err := h.guideService.UpdateGuide(uriInput.ID, &guideModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "guia não encontrado"})
			return
		}
		if err.Error() == "nenhum dado recebido para atualização" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.GuideResponse{
		GuideID:      updatedGuide.GuideID,
		Title:        updatedGuide.Title,
		BannerURL:    updatedGuide.BannerURL,
		Type:         updatedGuide.Type,
		Description:  updatedGuide.Description,
		Link:         updatedGuide.Link,
		CreationDate: updatedGuide.CreatedAt,
		Likes:        updatedGuide.Likes,
		Dislikes:     updatedGuide.Dislikes,
		UserID:       updatedGuide.UserID,
		CharacterID:  updatedGuide.CharacterID,
	}

	c.JSON(http.StatusOK, response)
}
