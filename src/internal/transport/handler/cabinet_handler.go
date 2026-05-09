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
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type CabinetHandler struct {
	cabinetService service.ICabinetService
}

func NewCabinetHandler(i do.Injector) (*CabinetHandler, error) {
	cabinetService := do.MustInvoke[service.ICabinetService](i)

	return &CabinetHandler{cabinetService: cabinetService}, nil
}

func (h *CabinetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	cabinet := rg.Group("/cabinet")
	{
		cabinet.POST("/items", h.AddItem)
		cabinet.DELETE("/items/", h.RemoveItem)
		cabinet.PATCH("/items/", h.UpdateQty)
	}
}

// AddItem godoc
func (h *CabinetHandler) AddItem(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
		return
	}
	userID := userIDValue.(uuid.UUID)

	var input req.AddItemDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewAddItemCmd(
		userID,
		input.MedicineID,
		input.DateOfManufacture,
		input.Quantity,
	)

	result, err := h.cabinetService.AddItem(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

// RemoveItem godoc
func (h *CabinetHandler) RemoveItem(c *gin.Context) {
	_, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
		return
	}

	var input req.RemoveItemDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewRemoveItemCmd(
		input.ID,
	)

	err = h.cabinetService.RemoveItem(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

// UpdateQty godoc
func (h *CabinetHandler) UpdateQty(c *gin.Context) {
	_, exists := c.Get(middleware.UserIDKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
		return
	}

	var input req.UpdateQtyDto

	err := c.BindJSON(&input)
	if err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewUpdateQtyCmd(
		input.ID,
		input.Qty,
	)

	result, err := h.cabinetService.UpdateQty(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
