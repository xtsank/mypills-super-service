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

type ProfileHandler struct {
	profile service.IProfileService
}

func NewProfileHandler(i do.Injector) (*ProfileHandler, error) {
	return &ProfileHandler{profile: do.MustInvoke[service.IProfileService](i)}, nil
}

func (h *ProfileHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/users/me", h.GetProfile)
	rg.PATCH("/users/me", h.UpdateProfile)
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	cmd := command.NewGetProfileCmd(userID)
	result, err := h.profile.GetProfile(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	var input req.ProfilePatchDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewUpdateProfileCmd(
		userID, input.Email, input.Sex, input.Weight, input.Age,
		input.IsPregnant, input.IsDriver, input.Illnesses, input.Allergies)
	result, err := h.profile.UpdateProfile(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
