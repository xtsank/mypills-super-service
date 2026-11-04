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

type MedicineCatalogHandler struct {
	service service.IMedicineCatalogService
}

func NewMedicineCatalogHandler(i do.Injector) (*MedicineCatalogHandler, error) {
	return &MedicineCatalogHandler{service: do.MustInvoke[service.IMedicineCatalogService](i)}, nil
}

func (h *MedicineCatalogHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/medicines", h.List)
}

func (h *MedicineCatalogHandler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	rg.POST("/medicines", h.Create)
	rg.PATCH("/medicines/:id", h.Update)
	rg.DELETE("/medicines/:id", h.Delete)
}

func (h *MedicineCatalogHandler) List(c *gin.Context) {
	input, ok := medicineListRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListMedicinesCmd(input.ID, input.Name, input.Limit, input.Offset)
	page, err := h.service.List(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, page)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *MedicineCatalogHandler) Create(c *gin.Context) {
	var input req.MedicineWriteDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewMedicineWriteCmd(input)
	result, err := h.service.Create(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Location", "/api/v2/medicines?id="+result.ID.String())
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *MedicineCatalogHandler) Update(c *gin.Context) {
	path, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	var input req.MedicinePatchDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewMedicinePatchCmd(path.ID, input)
	result, err := h.service.Update(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *MedicineCatalogHandler) Delete(c *gin.Context) {
	path, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	cmd := command.NewDeleteMedicineCmd(path.ID)
	if err := h.service.Delete(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
