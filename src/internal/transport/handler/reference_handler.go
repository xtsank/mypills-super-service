package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type ReferenceHandler struct{ dictionaries service.IDictionaryService }

func NewReferenceHandler(i do.Injector) (*ReferenceHandler, error) {
	return &ReferenceHandler{dictionaries: do.MustInvoke[service.IDictionaryService](i)}, nil
}

func (h *ReferenceHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/reference", h.Index)
	rg.GET("/reference/:kind", h.List)
}

func (h *ReferenceHandler) Index(c *gin.Context) {
	cmd := command.NewGetReferenceIndexCmd()
	result, err := h.dictionaries.ReferenceIndex(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *ReferenceHandler) List(c *gin.Context) {
	input, ok := referenceListRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListReferenceCmd(input.Kind, input.Limit, input.Offset)
	result, err := h.dictionaries.ReferenceList(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
