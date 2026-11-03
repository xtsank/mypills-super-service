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

type CabinetItemsHandler struct{ items service.ICabinetItemsService }

func NewCabinetItemsHandler(i do.Injector) (*CabinetItemsHandler, error) {
	return &CabinetItemsHandler{items: do.MustInvoke[service.ICabinetItemsService](i)}, nil
}

func (h *CabinetItemsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/cabinets/:cabinetId/items")
	group.GET("", h.List)
	group.POST("", h.Create)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

func (h *CabinetItemsHandler) ids(c *gin.Context) (req.CabinetItemsPathDto, bool) {
	userID, ok := requestUserID(c)
	if !ok {
		return req.CabinetItemsPathDto{}, false
	}
	path, ok := pathIDRequest(c, "cabinetId")
	return req.CabinetItemsPathDto{UserID: userID, CabinetID: path.ID}, ok
}

func (h *CabinetItemsHandler) List(c *gin.Context) {
	path, ok := h.ids(c)
	if !ok {
		return
	}
	input, ok := cabinetItemsListRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListCabinetItemsCmd(path.UserID, path.CabinetID, input.ID,
		input.LowStock, input.Expired, input.ExpiresWithinDays, input.Limit, input.Offset)
	page, err := h.items.List(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, page)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *CabinetItemsHandler) Create(c *gin.Context) {
	path, ok := h.ids(c)
	if !ok {
		return
	}
	var input req.CabinetItemCreateDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewCreateCabinetItemCmd(path.UserID, path.CabinetID, input)
	item, err := h.items.Create(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Location", "/api/v2/cabinets/"+path.CabinetID.String()+"/items?id="+item.ID.String())
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *CabinetItemsHandler) Update(c *gin.Context) {
	path, ok := h.ids(c)
	if !ok {
		return
	}
	itemPath, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	var input req.CabinetItemPatchDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewUpdateCabinetItemCmd(path.UserID, path.CabinetID, itemPath.ID, input)
	item, err := h.items.Update(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *CabinetItemsHandler) Delete(c *gin.Context) {
	path, ok := h.ids(c)
	if !ok {
		return
	}
	itemPath, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	cmd := command.NewDeleteCabinetItemCmd(path.UserID, path.CabinetID, itemPath.ID)
	if err := h.items.Delete(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
