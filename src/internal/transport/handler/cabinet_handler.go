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

type CabinetHandler struct{ catalog service.ICabinetService }

func NewCabinetHandler(i do.Injector) (*CabinetHandler, error) {
	return &CabinetHandler{catalog: do.MustInvoke[service.ICabinetService](i)}, nil
}

func (h *CabinetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/cabinets")
	group.GET("", h.List)
	group.POST("", h.Create)
	group.PATCH("/:cabinetId", h.Rename)
	group.DELETE("/:cabinetId", h.Delete)
}

func (h *CabinetHandler) List(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	input, ok := pageRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListCabinetsCmd(userID, input.ID, input.Limit, input.Offset)
	page, err := h.catalog.List(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, page)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *CabinetHandler) Create(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	var input req.CabinetWriteDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewCreateCabinetCmd(userID, input.Name)
	item, err := h.catalog.Create(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Location", "/api/v2/cabinets?id="+item.ID.String())
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *CabinetHandler) Rename(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	path, ok := pathIDRequest(c, "cabinetId")
	if !ok {
		return
	}
	var input req.CabinetWriteDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewRenameCabinetCmd(userID, path.ID, input.Name)
	item, err := h.catalog.Rename(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *CabinetHandler) Delete(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	path, ok := pathIDRequest(c, "cabinetId")
	if !ok {
		return
	}
	cmd := command.NewDeleteCabinetCmd(userID, path.ID)
	if err := h.catalog.Delete(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
