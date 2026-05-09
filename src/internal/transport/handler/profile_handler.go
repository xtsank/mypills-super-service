package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type ProfileHandler struct {
	profileService service.IProfileService
}

func NewProfileHandler(i do.Injector) (*ProfileHandler, error) {
	profileService := do.MustInvoke[service.IProfileService](i)

	return &ProfileHandler{profileService: profileService}, nil
}

func (h *ProfileHandler) RegisterRoutes(rg *gin.RouterGroup) {
	cabinet := rg.Group("/profile")
	{
		cabinet.PATCH("/me", h.UpdateProfile)
	}
}

// UpdateProfile godoc
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
		return
	}
	userID := userIDValue.(uuid.UUID)

	var input req.UpdateProfileDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewUpdateProfileCmd(
		userID,
		input.Sex,
		input.Weight,
		input.Age,
		input.IsPregnant,
		input.IsDriver,
		input.Illnesses,
		input.Allergies,
	)

	result, err := h.profileService.UpdateProfile(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
