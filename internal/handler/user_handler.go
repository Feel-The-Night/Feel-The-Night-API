package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dto "github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/middleware"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Criar um novo usuário
// @Description  Cria um usuário com os dados informados no corpo da requisição. Todo usuário novo começa como "iniciante".
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserInput  true  "Dados do usuário"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	logRequest(c, "User", "CreateUser")

	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	// Role nunca vem do body: o default do banco ("iniciante") é aplicado
	user := models.User{
		Nickname:  input.Nickname,
		Email:     input.Email,
		DiscordID: input.DiscordID,
		Password:  input.Password, // o service faz o hash
	}

	if err := h.userService.CreateUser(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "email já cadastrado"})
			return
		}
		internalError(c, "User", "CreateUser", err)
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
// @Success      200  {object}  dto.PublicUserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserById(c *gin.Context) {
	logRequest(c, "User", "GetUserById")

	var input dto.UserIDInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	user, err := h.userService.GetUserById(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		internalError(c, "User", "GetUserById", err)
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
// @Success      200    {object}  dto.UserResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	logRequest(c, "User", "GetUserByEmail")

	var input dto.UserEmailInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "email inválido"})
		return
	}

	user, err := h.userService.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		internalError(c, "User", "GetUserByEmail", err)
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
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUserById(c *gin.Context) {
	logRequest(c, "User", "DeleteUserById")

	var input dto.UserIDInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := h.userService.DeleteUserById(input.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		internalError(c, "User", "DeleteUserById", err)
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
// @Param        body  body      dto.UpdateUserInput  true  "Dados para atualização"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	logRequest(c, "User", "UpdateUser")

	var uriInput dto.UserIDInput
	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	// Só o próprio usuário ou um admin pode editar a conta
	claims := middleware.GetClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "não autenticado"})
		return
	}
	if models.Role(claims.Role) != models.RoleAdmin && uint(claims.UserID) != uint(uriInput.ID) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "você só pode editar sua própria conta"})
		return
	}

	var bodyInput dto.UpdateUserInput
	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	userModel := models.User{
		Nickname:  bodyInput.Nickname,
		Email:     bodyInput.Email,
		DiscordID: bodyInput.DiscordID,
		Password:  bodyInput.Password, // o service faz o hash se vier preenchida
	}

	updatedUser, err := h.userService.UpdateUser(uriInput.ID, &userModel)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "usuário não encontrado"})
		case errors.Is(err, services.ErrNoUpdateData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, gorm.ErrDuplicatedKey):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "email já cadastrado"})
		default:
			internalError(c, "User", "UpdateUser", err)
		}
		return
	}

	c.JSON(http.StatusOK, toUserResponse(updatedUser))
}

// toUserResponse: resposta completa (com email), para o próprio usuário e admin.
func toUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
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
func toPublicUserResponse(u *models.User) dto.PublicUserResponse {
	return dto.PublicUserResponse{
		UserID:           u.UserID,
		Nickname:         u.Nickname,
		DiscordID:        u.DiscordID,
		RegistrationDate: u.CreatedAt,
		EventCount:       u.EventCount,
		GuideCount:       u.GuideCount,
	}
}
