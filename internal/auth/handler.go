package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
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

// dummyHash is compared when the user does not exist, so the response time
// is similar in both cases (prevents discovering registered emails by timing).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)

// Login godoc
// @Summary      Login
// @Description  Authenticates with email and password and returns an access and a refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Credentials"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  httpx.ErrorResponse
// @Failure      401   {object}  httpx.ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	httpx.LogRequest(c, "Auth", "Login")

	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "email and password are required"})
		return
	}

	user, err := h.userService.GetUserByEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(input.Password))
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "invalid credentials"})
			return
		}
		// real error (database down, etc.): log it, but don't expose details
		zap.L().Error("[AuthHandler] Login: failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Error: "internal error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "invalid credentials"})
		return
	}

	resp, err := issueToken(user)
	if err != nil {
		zap.L().Error("[AuthHandler] Login: failed to generate tokens", zap.Error(err))
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Error: "internal error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Exchanges a valid refresh token for a new access and refresh token pair.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshRequest  true  "Refresh token"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  httpx.ErrorResponse
// @Failure      401   {object}  httpx.ErrorResponse
// @Router       /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	httpx.LogRequest(c, "Auth", "Refresh")

	var input RefreshRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Error: "refresh_token is required"})
		return
	}

	claims, err := jwt.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		msg := "invalid refresh token"
		if errors.Is(err, jwt.ErrExpiredToken) {
			msg = "refresh token expired, please log in again"
		}
		c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: msg})
		return
	}

	u, err := h.userService.GetUserById(c.Request.Context(), uint(claims.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Error: "invalid refresh token"})
			return
		}
		httpx.InternalError(c, "Auth", "Refresh", err)
		return
	}

	resp, err := issueToken(u)
	if err != nil {
		httpx.InternalError(c, "Auth", "Refresh", err)
		return
	}

	c.JSON(http.StatusOK, resp)

}

func issueToken(u *models.User) (TokenResponse, error) {
	access, err := jwt.GenerateAccessToken(int64(u.UserID), u.Email, u.Nickname, string(u.Role))
	if err != nil {
		return TokenResponse{}, err
	}
	refresh, err := jwt.GenerateRefreshToken(int64(u.UserID))
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(jwt.AccessTokenDuration.Seconds()),
	}, nil
}
