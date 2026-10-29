package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/apperr"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/middleware"
)

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{userService: userService}
}

// CreateUser godoc
// @Summary      Criar um novo usuário
// @Description  Cria um usuário com os dados informados no corpo da requisição. Todo usuário novo começa como "iniciante".
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      CreateInput  true  "Dados do usuário"
// @Success      201   {object}  Response
// @Failure      400   {object}  httpx.ErrorResponse
// @Failure      409   {object}  httpx.ErrorResponse
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	httpx.LogRequest(c, "User", "CreateUser")

	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	// Role nunca vem do body: o default do banco ("iniciante") é aplicado
	user := models.User{
		Nickname:  input.Nickname,
		Email:     input.Email,
		DiscordID: input.DiscordID,
		Password:  input.Password, // o service faz o hash
	}

	if err := h.userService.CreateUser(c.Request.Context(), &user); err != nil {
		if errors.Is(err, ErrEmailTaken) {
			c.JSON(http.StatusConflict, httpx.ErrorResponse{Error: "email já cadastrado"})
			return
		}
		httpx.InternalError(c, "User", "CreateUser", err)
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(&user))
}

// GetUserById godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna o perfil público do usuário (sem email).
// @Tags         User
// @Produce      json
// @Param        id   path      int  true  "ID do usuário" minimum(1)
// @Success      200  {object}  PublicResponse
// @Failure      400  {object}  httpx.ErrorResponse
// @Failure      404  {object}  httpx.ErrorResponse
// @Router       /users/{id} [get]
func (h *Handler) GetUserById(c *gin.Context) {
	httpx.LogRequest(c, "User", "GetUserById")

	var input IDInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	user, err := h.userService.GetUserById(c.Request.Context(), input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		httpx.InternalError(c, "User", "GetUserById", err)
		return
	}

	c.JSON(http.StatusOK, toPublicUserResponse(user))
}

// GetUserByEmail godoc
// @Summary      Buscar usuário por email (admin)
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Param        email  query     string  true  "Email do usuário"
// @Success      200    {object}  Response
// @Failure      400    {object}  httpx.ErrorResponse
// @Router       /users [get]
func (h *Handler) GetUserByEmail(c *gin.Context) {
	httpx.LogRequest(c, "User", "GetUserByEmail")

	var input EmailInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "email inválido"})
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		httpx.InternalError(c, "User", "GetUserByEmail", err)
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

// DeleteUserById godoc
// @Summary      Deletar usuário por ID (admin)
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID do usuário" minimum(1)
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httpx.ErrorResponse
// @Router       /users/{id} [delete]
func (h *Handler) DeleteUserById(c *gin.Context) {
	httpx.LogRequest(c, "User", "DeleteUserById")

	var input IDInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.userService.DeleteUserById(c.Request.Context(), input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		httpx.InternalError(c, "User", "DeleteUserById", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuário deletado com sucesso"})
}

// UpdateUser godoc
// @Summary      Atualizar dados do usuário
// @Description  O próprio usuário pode editar a sua conta; admin pode editar qualquer conta.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                  true  "ID do usuário" minimum(1)
// @Param        body  body      UpdateInput  true  "Dados para atualização"
// @Success      200   {object}  Response
// @Failure      400   {object}  httpx.ErrorResponse
// @Router       /users/{id} [patch]
func (h *Handler) UpdateUser(c *gin.Context) {
	httpx.LogRequest(c, "User", "UpdateUser")

	var uriInput IDInput
	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "id inválido"})
		return
	}

	// Só o próprio usuário ou um admin pode editar a conta
	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "não autenticado"})
		return
	}
	if models.Role(claims.Role) != models.RoleAdmin && uint(claims.UserID) != uint(uriInput.ID) {
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Error: "você só pode editar sua própria conta"})
		return
	}

	var bodyInput UpdateInput
	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	userModel := models.User{
		Nickname:  bodyInput.Nickname,
		Email:     bodyInput.Email,
		DiscordID: bodyInput.DiscordID,
		Password:  bodyInput.Password, // o service faz o hash se vier preenchida
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), uriInput.ID, &userModel)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, httpx.ErrorResponse{Error: "usuário não encontrado"})
		case errors.Is(err, apperr.ErrNoUpdateData):
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: err.Error()})
		case errors.Is(err, ErrEmailTaken):
			c.JSON(http.StatusConflict, httpx.ErrorResponse{Error: "email já cadastrado"})
		default:
			httpx.InternalError(c, "User", "UpdateUser", err)
		}
		return
	}

	c.JSON(http.StatusOK, toUserResponse(updatedUser))
}

// toUserResponse: resposta completa (com email), para o próprio usuário e admin.
func toUserResponse(u *models.User) Response {
	return Response{
		UserID:           u.UserID,
		Nickname:         u.Nickname,
		Email:            u.Email,
		DiscordID:        u.DiscordID,
		RegistrationDate: u.CreatedAt,
		EventCount:       u.EventCount,
		GuideCount:       u.GuideCount,
	}
}

// toPublicUserResponse: perfil público, sem email.
func toPublicUserResponse(u *models.User) PublicResponse {
	return PublicResponse{
		UserID:           u.UserID,
		Nickname:         u.Nickname,
		DiscordID:        u.DiscordID,
		RegistrationDate: u.CreatedAt,
		EventCount:       u.EventCount,
		GuideCount:       u.GuideCount,
	}
}
