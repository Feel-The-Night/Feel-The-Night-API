package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/auth"
	dto "github.com/kisalto/Feel-The-Night-Swagger/internal/dto"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// hash falso usado quando o usuário não existe, para o tempo de resposta
// ser parecido nos dois casos (evita descobrir emails cadastrados pelo tempo).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

// Login godoc
// @Summary      Login
// @Description  Autentica com email e senha e retorna access e refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginInput  true  "Credenciais"
// @Success      200   {object}  dto.TokenResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	zap.L().Info("[AuthHandler] Login",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "email e senha são obrigatórios"})
		return
	}

	user, err := h.userService.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// usuário não existe: roda o bcrypt mesmo assim e falha com a MESMA mensagem
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(input.Password))
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "credenciais inválidas"})
			return
		}
		// erro de verdade (banco fora, etc.): loga, mas não expõe detalhes
		zap.L().Error("[AuthHandler] Login: erro ao buscar usuário", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "erro interno"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "credenciais inválidas"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(int64(user.UserID), user.Email, user.Nickname)
	if err != nil {
		zap.L().Error("[AuthHandler] Login: erro ao gerar access token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "erro interno"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(int64(user.UserID))
	if err != nil {
		zap.L().Error("[AuthHandler] Login: erro ao gerar refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "erro interno"})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(auth.AccessTokenDuration.Seconds()),
	})
}
