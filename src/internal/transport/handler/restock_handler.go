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

type RestockHandler struct {
	restockService service.IRestockService
}

func NewRestockHandler(i do.Injector) (*RestockHandler, error) {
	restockService := do.MustInvoke[service.IRestockService](i)
	return &RestockHandler{restockService: restockService}, nil
}

func (h *RestockHandler) RegisterRoutes(rg *gin.RouterGroup) {
	restock := rg.Group("/restock-list")
	{
		restock.GET("", h.List)
		restock.POST("", h.Create)
		restock.PATCH("/:id", h.Update)
		restock.DELETE("/:id", h.Delete)
	}
}

func (h *RestockHandler) List(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	input, ok := pageRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListRestockCmd(userID, input.ID, input.Limit, input.Offset)
	page, err := h.restockService.List(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, page)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *RestockHandler) Create(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	var input req.RestockItemCreateDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewCreateRestockCmd(userID, input.MedicineID, input.DesiredQuantity)
	item, err := h.restockService.Create(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Location", "/api/v2/restock-list?id="+item.ID.String())
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusCreated)
}

func (h *RestockHandler) Update(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	path, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	var input req.RestockItemPatchDto
	if !bindBody(c, &input) {
		return
	}
	cmd := command.NewUpdateRestockCmd(userID, path.ID, input.DesiredQuantity)
	item, err := h.restockService.Update(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, item)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *RestockHandler) Delete(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	path, ok := pathIDRequest(c, "id")
	if !ok {
		return
	}
	cmd := command.NewDeleteRestockCmd(userID, path.ID)
	if err := h.restockService.Delete(c.Request.Context(), cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
