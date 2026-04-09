package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xtsank/mypills-super-service/internal/service/command"
	"github.com/xtsank/mypills-super-service/internal/transport/dto"
	"github.com/xtsank/mypills-super-service/internal/transport/middleware"

	"github.com/xtsank/mypills-super-service/internal/service"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.CreateUserDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	cmd := command.CreateUserCmd{
		Login:      input.Login,
		Password:   input.Password,
		Sex:        input.Sex,
		Weight:     input.Weight,
		Age:        input.Age,
		IsPregnant: input.IsPregnant,
		IsDriver:   input.IsDriver,
		Illnesses:  input.Illnesses,
		Allergies:  input.Allergies,
	}

	createdUser, err := h.authService.Register(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, gin.H{"id": createdUser.ID})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
