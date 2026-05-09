package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type AdminHandler struct {
	adminService service.IAdminService
}

func NewAdminHandler(i do.Injector) (*AdminHandler, error) {
	adminService := do.MustInvoke[service.IAdminService](i)

	return &AdminHandler{adminService: adminService}, nil
}

func (h *AdminHandler) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		admin.POST("/medicine", h.AddMedicine)
		admin.PATCH("/medicine", h.UpdateMedicine)
		admin.DELETE("/medicine", h.RemoveMedicine)
		admin.PATCH("/medicine/indications", h.UpdateIndications)
		admin.PATCH("/medicine/contraindications", h.UpdateContraindications)
		admin.PATCH("/medicine/composition", h.UpdateComposition)
		admin.POST("/medicine/dosage", h.AddDosageRule)
		admin.DELETE("/medicine/dosage", h.RemoveDosageRule)
	}
}

func (h *AdminHandler) AddMedicine(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.AddMedicineDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd, err := command.NewAddMedicineCmd(
		input.Name,
		input.ExpireTime,
		input.IsPrescription,
		input.MethodOfApplication,
		input.EffectOnPregnant,
		input.EffectOnDriver,
		input.Form,
		input.Unit,
		input.Recommendations,
		input.Contraindications,
		input.Substances,
		input.Dosages,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	result, err := h.adminService.AddMedicine(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *AdminHandler) UpdateMedicine(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.UpdateMedicineDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewUpdateMedicineCmd(
		input.ID,
		input.ExpireTime,
		input.IsPrescription,
		input.MethodOfApplication,
		input.EffectOnPregnant,
		input.EffectOnDriver,
		input.FormID,
		input.UnitID,
	)

	result, err := h.adminService.UpdateMedicine(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) RemoveMedicine(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.RemoveMedicineDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewRemoveMedicineCmd(input.ID)
	if err := h.adminService.RemoveMedicine(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) UpdateIndications(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.UpdateLinksDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewUpdateLinksCmd(input.MedicineID, input.IDs)
	if err := h.adminService.UpdateIndications(c.Request.Context(), *cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) UpdateContraindications(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.UpdateLinksDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewUpdateLinksCmd(input.MedicineID, input.IDs)
	if err := h.adminService.UpdateContraindications(c.Request.Context(), *cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) UpdateComposition(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.UpdateCompositionDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd, err := command.NewUpdateCompositionCmd(input.MedicineID, input.Substances)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.adminService.UpdateComposition(c.Request.Context(), *cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) AddDosageRule(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.AddDosageRuleDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd, err := command.NewAddDosageRuleCmd(input.MedicineID, input.Dosage)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.adminService.AddDosageRule(c.Request.Context(), *cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *AdminHandler) RemoveDosageRule(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.RemoveDosageRuleDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := command.NewRemoveDosageRuleCmd(input.RuleID)
	if err := h.adminService.DeleteDosageRule(c.Request.Context(), *cmd); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, res.SuccessResDTO{Status: "success"})
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func isAdmin(c *gin.Context) bool {
	isAdminValue, exists := c.Get(middleware.IsAdminKey)
	if !exists {
		_ = c.Error(errors.ErrUnauthorized)
		return false
	}
	isAdmin, ok := isAdminValue.(bool)
	if !ok || !isAdmin {
		_ = c.Error(errors.ErrUnauthorized)
		return false
	}
	return true
}
