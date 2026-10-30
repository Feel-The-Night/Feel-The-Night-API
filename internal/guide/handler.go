package guide

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/middleware"
)

type Handler struct {
	guideService *Service
}

func NewHandler(guideService *Service) *Handler {
	return &Handler{guideService: guideService}
}

// CreateGuide godoc
// @Summary      Create a new guide
// @Description  Creates a guide from the request body.
// @Tags         Guide
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        guide  body      CreateInput  true  "Guide data"
// @Success      201    {object}  Response
// @Router       /guides [post]
func (h *Handler) CreateGuide(c *gin.Context) {
	httpx.LogRequest(c, "Guide", "CreateGuide")

	var input CreateInput
	claims := middleware.GetClaims(c)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
		return
	}

	guide := models.Guide{
		Title:       input.Title,
		BannerURL:   input.BannerURL,
		Type:        input.Type,
		Description: input.Description,
		Link:        input.Link,
		CharacterID: input.CharacterID,
		UserID:      uint(claims.UserID),
	}

	if err := h.guideService.CreateGuide(c.Request.Context(), &guide); err != nil {
		switch {
		case errors.Is(err, ErrCharacterNotFound),
			errors.Is(err, gorm.ErrForeignKeyViolated):
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "character not found"})
		default:
			httpx.InternalError(c, "Guide", "CreateGuide", err)
		}
		return
	}

	c.JSON(http.StatusCreated, toGuideResponse(&guide))
}

// GetGuideById godoc
// @Summary      Get guide by ID
// @Tags         Guide
// @Produce      json
// @Param        id   path      int  true  "Guide ID" minimum(1)
// @Success      200  {object}  Response
// @Router       /guides/{id} [get]
func (h *Handler) GetGuideById(c *gin.Context) {
	httpx.LogRequest(c, "Guide", "GetGuideById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "invalid id"})
		return
	}

	guide, err := h.guideService.GetGuideById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "guide not found"})
			return
		}
		httpx.InternalError(c, "Guide", "GetGuideById", err)
		return
	}

	c.JSON(http.StatusOK, toGuideResponse(guide))
}

// DeleteGuideById godoc
// @Summary      Delete guide by ID
// @Tags         Guide
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Guide ID" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /guides/{id} [delete]
func (h *Handler) DeleteGuideById(c *gin.Context) {
	httpx.LogRequest(c, "Guide", "DeleteGuideById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.guideService.DeleteGuideById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "guide not found"})
			return
		}
		httpx.InternalError(c, "Guide", "DeleteGuideById", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "guide deleted successfully"})
}

// UpdateGuide godoc
// @Summary      Update guide data
// @Tags         Guide
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                   true  "Guide ID" minimum(1)
// @Param        body  body      UpdateInput  true  "Update data"
// @Success      200   {object}  Response
// @Router       /guides/{id} [patch]
func (h *Handler) UpdateGuide(c *gin.Context) {
	httpx.LogRequest(c, "Guide", "UpdateGuide")

	var uriInput IDInput
	var bodyInput UpdateInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "invalid id"})
		return
	}

	claims := middleware.GetClaims(c)

	guide, err := h.guideService.GetGuideById(c.Request.Context(), uriInput.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "guide not found"})
			return
		}
		httpx.InternalError(c, "Guide", "UpdateGuide", err)
		return
	}

	if models.Role(claims.Role) != models.RoleAdmin && guide.UserID != uint(claims.UserID) {
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Error: "you can only edit your own guides"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "invalid request body"})
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
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "guide not found"})
		case errors.Is(err, apperr.ErrNoUpdateData):
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
		default:
			httpx.InternalError(c, "Guides", "UpdateGuide", err)
		}
		return
	}

	c.JSON(http.StatusOK, toGuideResponse(updatedGuide))
}

// GetGuides godoc
// @Summary      List guides
// @Description  Lists guides with pagination. Title is an optional filter (partial match).
// @Tags         Guide
// @Produce      json
// @Param        title  query     string  false  "Filter by title"
// @Param        page   query     int     false  "Page"            default(1)  minimum(1)
// @Param        limit  query     int     false  "Items per page"  default(20) minimum(1) maximum(100)
// @Success      200    {object}  httpx.PageResponse[Response]
// @Failure      400    {object}  httpx.ErrorResponse
// @Router       /guides [get]
func (h *Handler) GetGuides(c *gin.Context) {
	httpx.LogRequest(c, "Guide", "GetGuides")

	var query ListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "invalid parameters"})
		return
	}

	guides, total, err := h.guideService.GetGuides(c.Request.Context(), strings.TrimSpace(query.Title), query.Page, query.Limit)
	if err != nil {
		httpx.InternalError(c, "Guide", "GetGuides", err)
		return
	}

	data := make([]Response, 0, len(guides))
	for i := range guides {
		data = append(data, toGuideResponse(&guides[i]))
	}

	c.JSON(http.StatusOK, httpx.PageResponse[Response]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	})
}

func toGuideResponse(g *models.Guide) Response {
	return Response{
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
