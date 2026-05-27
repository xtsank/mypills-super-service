package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type DictionaryHandler struct {
	dictionaryService service.IDictionaryService
}

func NewDictionaryHandler(i do.Injector) (*DictionaryHandler, error) {
	service := do.MustInvoke[service.IDictionaryService](i)
	return &DictionaryHandler{dictionaryService: service}, nil
}

func (h *DictionaryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dict := rg.Group("/dictionary")
	{
		dict.GET("/illnesses", h.ListIllnesses)
		dict.GET("/substances", h.ListSubstances)
		dict.GET("/forms", h.ListForms)
		dict.GET("/units", h.ListUnits)
		dict.GET("/medicines", h.ListMedicines)
	}
}

func (h *DictionaryHandler) ListIllnesses(c *gin.Context) {
	items, err := h.dictionaryService.ListIllnesses(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *DictionaryHandler) ListSubstances(c *gin.Context) {
	items, err := h.dictionaryService.ListSubstances(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *DictionaryHandler) ListForms(c *gin.Context) {
	items, err := h.dictionaryService.ListForms(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *DictionaryHandler) ListUnits(c *gin.Context) {
	items, err := h.dictionaryService.ListUnits(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *DictionaryHandler) ListMedicines(c *gin.Context) {
	items, err := h.dictionaryService.ListMedicines(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

