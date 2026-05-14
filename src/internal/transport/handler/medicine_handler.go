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

type MedicineHandler struct {
	medicineService service.IMedicineService
}

func NewMedicineHandler(i do.Injector) (*MedicineHandler, error) {
	medicineService := do.MustInvoke[service.IMedicineService](i)

	return &MedicineHandler{medicineService}, nil
}

func (h *MedicineHandler) RegisterRoutes(rg *gin.RouterGroup) {
	cabinet := rg.Group("/medicine")
	{
		cabinet.POST("/select", h.Select)
	}
}

// Select godoc
// @Summary      Подбор лекарств
// @Description  Возвращает рекомендации по лекарствам для болезни
// @Tags         Medicine
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.SelectMedicineDto  true  "Данные для подбора"
// @Success      200    {object}  res.MedicineResDto     "Рекомендации"
// @Failure      400    {object}  errors.AppError        "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError        "Пользователь не авторизован"
// @Failure      404    {object}  errors.AppError        "Пользователь не найден"
// @Router       /medicine/select [post]
func (h *MedicineHandler) Select(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized.WithSource())
		return
	}
	userID := userIDValue.(uuid.UUID)

	var input req.SelectMedicineDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewSelectMedicineCmd(
		userID,
		input.IllnessID,
	)

	result, err := h.medicineService.Select(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
