package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/middleware"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

const dayLayout = "2006-01-02"

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
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
// @Param        event  body      dto.CreateEventInput  true  "Dados do Evento"
// @Success      201    {object}  dto.EventResponse
// @Router       /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	logRequest(c, "Event", "CreateEvent")

	var input dto.CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	day, err := parseDay(input.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "data inválida, use AAAA-MM-DD"})
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
		internalError(c, "Event", "CreateEvent", err)
		return
	}

	c.JSON(http.StatusCreated, toEventResponse(&event))
}

// GetEventById godoc
// @Summary      Buscar evento por ID
// @Tags         Event
// @Produce      json
// @Param        id   path      int  true  "ID do evento" minimum(1)
// @Success      200  {object}  dto.EventResponse
// @Router       /events/{id} [get]
func (h *EventHandler) GetEventById(c *gin.Context) {
	logRequest(c, "Event", "GetEventById")

	var input dto.EventIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	event, err := h.eventService.GetEventById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		internalError(c, "Event", "GetEventById", err)
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
func (h *EventHandler) DeleteEventById(c *gin.Context) {
	logRequest(c, "Event", "DeleteEventById")

	var input dto.EventIDInput

	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.eventService.DeleteEventById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		internalError(c, "Event", "DeleteEventById", err)
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
// @Param        body  body      dto.UpdateEventInput true  "Dados para atualização"
// @Success      200   {object}  dto.EventResponse
// @Router       /events/{id} [patch]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	logRequest(c, "Event", "UpdateEvent")

	var uriInput dto.EventIDInput
	var bodyInput dto.UpdateEventInput
	var day time.Time

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	if bodyInput.Day != "" {
		parsed, err := parseDay(bodyInput.Day)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "data inválida, use AAAA-MM-DD"})
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
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "evento não encontrado"})
			return
		}
		if errors.Is(err, services.ErrNoUpdateData) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		internalError(c, "Event", "UpdateEvent", err)
		return
	}

	c.JSON(http.StatusOK, toEventResponse(updatedEvent))
}

// GetEvents godoc
// @Summary      listar eventos
// @Description  Lista eventos com filtros opcionais de título e paginação.
// @Tags         Event
// @Produce      json
// @Param        title  query     string  false  "Filtrar por título"
// @Param        page   query     int     false  "Página"           default(1)  minimum(1)
// @Param        limit  query     int     false  "Itens por página"  default(20) minimum(1) maximum(100)
// @Success      200    {object}  dto.PageResponse[dto.EventResponse]
// @Router       /events [get]
func (h *EventHandler) GetEvents(c *gin.Context) {
	logRequest(c, "Event", "GetEvents")

	var query dto.EventListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "parâmetros de query inválidos"})
		return
	}

	events, total, err := h.eventService.GetEvents(c.Request.Context(), query.Title, query.Upcoming, query.Page, query.Limit)
	if err != nil {
		internalError(c, "Event", "GetEvents", err)
		return
	}

	data := make([]dto.EventResponse, len(events))
	for i, e := range events {
		data[i] = toEventResponse(&e)
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"total":  total,
	// 	"events": eventResponses,
	// })

	c.JSON(http.StatusOK, dto.PageResponse[dto.EventResponse]{
		Data:  data,
		Page:  query.Page,
		Limit: query.Limit,
		Total: total,
	})
}

func toEventResponse(e *models.Event) dto.EventResponse {
	return dto.EventResponse{
		EventID:     e.EventID,
		Title:       e.Title,
		Description: e.Description,
		BannerURL:   e.BannerURL,
		Day:         e.Day.Format(dayLayout),
		UserID:      e.UserID,
	}
}
