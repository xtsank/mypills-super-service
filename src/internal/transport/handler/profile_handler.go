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
	_ "github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
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
// @Summary      Обновление данных профиля
// @Description  Позволяет изменить параметры пользователя (вес, возраст, болезни и т.д.).
// @Description  Передавайте только те поля, которые нужно изменить.
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.UpdateProfileDto  true  "Новые данные профиля"
// @Success      200    {object}  res.ProfileResDto     "Профиль успешно обновлен"
// @Failure      400    {object}  errors.AppError       "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError       "Пользователь не авторизован"
// @Router       /profile/me [patch]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized.WithSource())
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
