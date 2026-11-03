package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type AuthHandler struct{ auth service.IAuthService }

func NewAuthHandler(i do.Injector) (*AuthHandler, error) {
	return &AuthHandler{auth: do.MustInvoke[service.IAuthService](i)}, nil
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users", h.Register)
	rg.POST("/sessions", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input req.RegisterUserDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewCreateUserCmd(
		input.Login, input.Email, input.Password, input.Sex, input.Weight, input.Age,
		input.IsPregnant, input.IsDriver, input.Illnesses, input.Allergies)
	result, err := h.auth.Register(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Location", "/api/v2/users/me")
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input req.SessionDto
	if !bindBody(c, &input) {
		return
	}
	cmd, err := command.NewLoginUserCmd(input.Login, input.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	result, err := h.auth.Login(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
