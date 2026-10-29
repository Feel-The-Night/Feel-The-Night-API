package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/httpx"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/jwt"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{userService: userService}
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
// @Param        body  body      LoginRequest  true  "Credenciais"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  httpx.ErrorResponse
// @Failure      401   {object}  httpx.ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	httpx.LogRequest(c, "Auth", "Login")

	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "email e senha são obrigatórios"})
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(input.Password))
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "credenciais inválidas"})
			return
		}
		// erro de verdade (banco fora, etc.): loga, mas não expõe detalhes
		zap.L().Error("[AuthHandler] Login: erro ao buscar usuário", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Error: "erro interno"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "credenciais inválidas"})
		return
	}

	accessToken, err := jwt.GenerateAccessToken(int64(user.UserID), user.Email, user.Nickname, string(user.Role))
	if err != nil {
		zap.L().Error("[AuthHandler] Login: erro ao gerar access token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Error: "erro interno"})
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(int64(user.UserID))
	if err != nil {
		zap.L().Error("[AuthHandler] Login: erro ao gerar refresh token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Error: "erro interno"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(jwt.AccessTokenDuration.Seconds()),
	})
}
