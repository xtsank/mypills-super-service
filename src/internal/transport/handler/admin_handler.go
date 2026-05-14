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
	_ "github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
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

// AddMedicine godoc
// @Summary      Добавление лекарства
// @Description  Создает новое лекарство
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.AddMedicineDto   true  "Данные лекарства"
// @Success      201    {object}  res.AdminResDto      "Лекарство создано"
// @Failure      400    {object}  errors.AppError      "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError      "Недостаточно прав"
// @Router       /admin/medicine [post]
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

// UpdateMedicine godoc
// @Summary      Обновление лекарства
// @Description  Обновляет данные лекарства
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.UpdateMedicineDto  true  "Поля для обновления"
// @Success      200    {object}  res.AdminResDto        "Лекарство обновлено"
// @Failure      400    {object}  errors.AppError        "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError        "Недостаточно прав"
// @Failure      404    {object}  errors.AppError        "Лекарство не найдено"
// @Router       /admin/medicine [patch]
func (h *AdminHandler) UpdateMedicine(c *gin.Context) {
	if !isAdmin(c) {
		return
	}

	var input req.UpdateMedicineDto
	if err := c.BindJSON(&input); err != nil {
		_ = c.Error(errors.ErrInvalidInput.WithError(err))
		return
	}

	cmd := &command.UpdateMedicineCmd{
		ID:                  input.ID,
		ExpireTime:          input.ExpireTime,
		IsPrescription:      input.IsPrescription,
		MethodOfApplication: input.MethodOfApplication,
		EffectOnPregnant:    input.EffectOnPregnant,
		EffectOnDriver:      input.EffectOnDriver,
		FormID:              input.FormID,
		UnitID:              input.UnitID,
	}

	result, err := h.adminService.UpdateMedicine(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

// RemoveMedicine godoc
// @Summary      Удаление лекарства
// @Description  Удаляет лекарство по идентификатору
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.RemoveMedicineDto  true  "Идентификатор лекарства"
// @Success      200    {object}  res.SuccessResDTO      "Лекарство удалено"
// @Failure      400    {object}  errors.AppError        "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError        "Недостаточно прав"
// @Router       /admin/medicine [delete]
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

// UpdateIndications godoc
// @Summary      Обновление показаний
// @Description  Обновляет список показаний для лекарства
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.UpdateLinksDto  true  "Список показаний"
// @Success      200    {object}  res.SuccessResDTO   "Показания обновлены"
// @Failure      400    {object}  errors.AppError     "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError     "Недостаточно прав"
// @Router       /admin/medicine/indications [patch]
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

// UpdateContraindications godoc
// @Summary      Обновление противопоказаний
// @Description  Обновляет список противопоказаний для лекарства
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.UpdateLinksDto  true  "Список противопоказаний"
// @Success      200    {object}  res.SuccessResDTO   "Противопоказания обновлены"
// @Failure      400    {object}  errors.AppError     "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError     "Недостаточно прав"
// @Router       /admin/medicine/contraindications [patch]
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

// UpdateComposition godoc
// @Summary      Обновление состава
// @Description  Обновляет список действующих веществ
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.UpdateCompositionDto  true  "Состав лекарства"
// @Success      200    {object}  res.SuccessResDTO         "Состав обновлен"
// @Failure      400    {object}  errors.AppError           "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError           "Недостаточно прав"
// @Router       /admin/medicine/composition [patch]
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

// AddDosageRule godoc
// @Summary      Добавление правила дозировки
// @Description  Добавляет правило дозировки к лекарству
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.AddDosageRuleDto  true  "Правило дозировки"
// @Success      200    {object}  res.SuccessResDTO     "Правило дозировки добавлено"
// @Failure      400    {object}  errors.AppError       "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError       "Недостаточно прав"
// @Router       /admin/medicine/dosage [post]
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

// RemoveDosageRule godoc
// @Summary      Удаление правила дозировки
// @Description  Удаляет правило дозировки по идентификатору
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        input  body      req.RemoveDosageRuleDto  true  "Идентификатор правила"
// @Success      200    {object}  res.SuccessResDTO        "Правило дозировки удалено"
// @Failure      400    {object}  errors.AppError          "Невалидные входные данные"
// @Failure      401    {object}  errors.AppError          "Недостаточно прав"
// @Router       /admin/medicine/dosage [delete]
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
		_ = c.Error(errors.ErrUnauthorized.WithSource())
		return false
	}
	isAdmin, ok := isAdminValue.(bool)
	if !ok || !isAdmin {
		_ = c.Error(errors.ErrUnauthorized.WithSource())
		return false
	}
	return true
}
