package event

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/middleware"
)

const dayLayout = "2006-01-02"

type Handler struct {
	eventService *Service
}

func NewHandler(eventService *Service) *Handler {
	return &Handler{eventService: eventService}
}

func parseDay(s string) (time.Time, error) {
	return time.Parse(dayLayout, s)
}

// CreateEvent godoc
// @Summary      Criar um novo evento
// @Description  Cria um novo evento com os dados informados no corpo da requisição.
// @Tags         Event
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event  body      CreateInput  true  "Dados do Evento"
// @Success      201    {object}  Response
// @Router       /events [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	httpx.LogRequest(c, "Event", "CreateEvent")

	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	day, err := parseDay(input.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "data inválida, use AAAA-MM-DD"})
		return
	}

	claims := middleware.GetClaims(c)

	event := models.Event{
		Title:       input.Title,
		Description: input.Description,
		BannerURL:   input.BannerURL,
		Day:         day,
		UserID:      uint(claims.UserID), // vem do token
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), &event); err != nil {
		httpx.InternalError(c, "Event", "CreateEvent", err)
		return
	}

	c.JSON(http.StatusCreated, toEventResponse(&event))
}

// GetEventById godoc
// @Summary      Buscar evento por ID
// @Tags         Event
// @Produce      json
// @Param        id   path      int  true  "ID do evento" minimum(1)
// @Success      200  {object}  Response
// @Router       /events/{id} [get]
func (h *Handler) GetEventById(c *gin.Context) {
	httpx.LogRequest(c, "Event", "GetEventById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	event, err := h.eventService.GetEventById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		httpx.InternalError(c, "Event", "GetEventById", err)
		return
	}

	c.JSON(http.StatusOK, toEventResponse(event))
}

// DeleteEventById godoc
// @Summary      Deletar evento por ID
// @Tags         Event
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do evento" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /events/{id} [delete]
func (h *Handler) DeleteEventById(c *gin.Context) {
	httpx.LogRequest(c, "Event", "DeleteEventById")

	var input IDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.eventService.DeleteEventById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		httpx.InternalError(c, "Event", "DeleteEventById", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "evento deletado com sucesso"})
}

// UpdateEvent godoc
// @Summary      Atualizar dados do evento
// @Tags         Event
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                      true  "ID do evento" minimum(1)
// @Param        body  body      UpdateInput true  "Dados para atualização"
// @Success      200   {object}  Response
// @Router       /events/{id} [patch]
func (h *Handler) UpdateEvent(c *gin.Context) {
	httpx.LogRequest(c, "Event", "UpdateEvent")

	var uriInput IDInput
	var bodyInput UpdateInput
	var day time.Time

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	if bodyInput.Day != "" {
		parsed, err := parseDay(bodyInput.Day)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "data inválida, use AAAA-MM-DD"})
			return
		}
		day = parsed
	}

	eventModel := models.Event{
		Title:       bodyInput.Title,
		Description: bodyInput.Description,
		BannerURL:   bodyInput.BannerURL,
		Day:         day,
	}

	updatedEvent, err := h.eventService.UpdateEventById(c.Request.Context(), uriInput.ID, &eventModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		if errors.Is(err, apperr.ErrNoUpdateData) {
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
			return
		}
		httpx.InternalError(c, "Event", "UpdateEvent", err)
		return
	}

	c.JSON(http.StatusOK, toEventResponse(updatedEvent))
}

// GetEvents godoc
// @Summary      listar eventos
// @Description  Lista eventos com filtros opcionais de título e paginação.
// @Tags         Event
// @Produce      json
// @Param        title    query     string  false  "Filtrar por título"
// @Param        upcoming query     bool    false  "Só eventos futuros"
// @Param        page     query     int     false  "Página"           default(1)  minimum(1)
// @Param        limit    query     int     false  "Itens por página"  default(20) minimum(1) maximum(100)
// @Success      200      {object}  httpx.PageResponse[Response]
// @Router       /events  [get]
func (h *Handler) GetEvents(c *gin.Context) {
	httpx.LogRequest(c, "Event", "GetEvents")

	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "parâmetros de query inválidos"})
		return
	}

	events, total, err := h.eventService.GetEvents(c.Request.Context(), query.Title, query.Upcoming, query.Page, query.Limit)
	if err != nil {
		httpx.InternalError(c, "Event", "GetEvents", err)
		return
	}

	data := make([]Response, len(events))
	for i, e := range events {
		data[i] = toEventResponse(&e)
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"total":  total,
	// 	"events": eventResponses,
	// })

	c.JSON(http.StatusOK, httpx.PageResponse[Response]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	})
}

func toEventResponse(e *models.Event) Response {
	return Response{
		EventID:     e.EventID,
		Title:       e.Title,
		Description: e.Description,
		BannerURL:   e.BannerURL,
		Day:         e.Day.Format(dayLayout),
		UserID:      e.UserID,
	}
}
