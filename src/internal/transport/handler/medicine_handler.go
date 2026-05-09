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
func (h *MedicineHandler) Select(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
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
