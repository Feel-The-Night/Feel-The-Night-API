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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary      Criar um novo usuário
// @Description  Cria um usuário com os dados informados no corpo da requisição.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserInput  true  "Dados do usuário"
// @Success      201   {object}  dto.UserResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	zap.L().Info("[UserHandler] CreateUser",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user := models.User{
		Nickname:  input.Nickname,
		Email:     input.Email,
		DiscordID: input.DiscordID,
		Password:  input.Password,
	}

	if err := h.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		Email:            user.Email,
		DiscordID:        user.DiscordID,
		RegistrationDate: user.RegistrationDate,
		EventCount:       user.EventCount,
		GuideCount:       user.GuideCount,
	}

	c.JSON(http.StatusCreated, response)
}

// GetUserById godoc
// @Summary      Buscar usuário por ID
// @Tags         User
// @Produce      json
// @Param        id   path      dto.UserIDInput  true  "ID do usuário" minimum(1)
// @Success      200  {object}  dto.UserResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserById(c *gin.Context) {
	zap.L().Info("[UserHandler] GetUserById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		Email:            user.Email,
		DiscordID:        user.DiscordID,
		RegistrationDate: user.RegistrationDate,
		EventCount:       user.EventCount,
		GuideCount:       user.GuideCount,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByEmail godoc
// @Summary      Buscar usuário por email
// @Tags         User
// @Produce      json
// @Param        email  query     string  true  "Email do usuário"
// @Success      200    {object}  dto.UserResponse
// @Failure      400    {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	zap.L().Info("[UserHandler] GetUserByEmail",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		Email:            user.Email,
		DiscordID:        user.DiscordID,
		RegistrationDate: user.RegistrationDate,
		EventCount:       user.EventCount,
		GuideCount:       user.GuideCount,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUserById godoc
// @Summary      Deletar usuário por ID
// @Tags         User
// @Produce      json
// @Param        id   path      int  true  "ID do usuário" minimum(1)
// @Success      200  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUserById(c *gin.Context) {
	zap.L().Info("[UserHandler] DeleteUserById",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuário deletado com sucesso"})
}

// UpdateUser godoc
// @Summary      Atualizar dados do usuário
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id    path      int                  true  "ID do usuário" minimum(1)
// @Param        body  body      dto.UpdateUserInput  true  "Dados para atualização"
// @Success      200   {object}  dto.UserResponse
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	zap.L().Info("[UserHandler] UpdateUser",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var uriInput dto.UserIDInput
	var bodyInput dto.UpdateUserInput

	if err := c.ShouldBindUri(&uriInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := c.ShouldBindJSON(&bodyInput); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	userModel := models.User{
		Nickname:  bodyInput.Nickname,
		Email:     bodyInput.Email,
		DiscordID: bodyInput.DiscordID,
		Password:  bodyInput.Password,
	}

	updatedUser, err := h.userService.UpdateUser(uriInput.ID, &userModel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "usuário não encontrado"})
			return
		}
		if err.Error() == "nenhum dado recebido para atualização" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.UserResponse{
		UserID:           updatedUser.UserID,
		Nickname:         updatedUser.Nickname,
		Email:            updatedUser.Email,
		DiscordID:        updatedUser.DiscordID,
		RegistrationDate: updatedUser.RegistrationDate,
		EventCount:       updatedUser.EventCount,
		GuideCount:       updatedUser.GuideCount,
	}

	c.JSON(http.StatusOK, response)
}
