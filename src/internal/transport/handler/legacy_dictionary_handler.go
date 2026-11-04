package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type LegacyDictionaryHandler struct {
	dictionaryService service.IDictionaryService
}

func NewLegacyDictionaryHandler(i do.Injector) (*LegacyDictionaryHandler, error) {
	service := do.MustInvoke[service.IDictionaryService](i)
	return &LegacyDictionaryHandler{dictionaryService: service}, nil
}

func (h *LegacyDictionaryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	dict := rg.Group("/dictionary")
	{
		dict.GET("/illnesses", h.ListIllnesses)
		dict.GET("/substances", h.ListSubstances)
		dict.GET("/forms", h.ListForms)
		dict.GET("/units", h.ListUnits)
		dict.GET("/medicines", h.ListMedicines)
	}
}

func (h *LegacyDictionaryHandler) ListIllnesses(c *gin.Context) {
	items, err := h.dictionaryService.ListIllnesses(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *LegacyDictionaryHandler) ListSubstances(c *gin.Context) {
	items, err := h.dictionaryService.ListSubstances(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *LegacyDictionaryHandler) ListForms(c *gin.Context) {
	items, err := h.dictionaryService.ListForms(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *LegacyDictionaryHandler) ListUnits(c *gin.Context) {
	items, err := h.dictionaryService.ListUnits(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}

func (h *LegacyDictionaryHandler) ListMedicines(c *gin.Context) {
	items, err := h.dictionaryService.ListMedicines(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, items)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
