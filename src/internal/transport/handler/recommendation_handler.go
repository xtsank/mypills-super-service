package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/service"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

type RecommendationHandler struct {
	service service.IRecommendationService
}

func NewRecommendationHandler(i do.Injector) (*RecommendationHandler, error) {
	return &RecommendationHandler{service: do.MustInvoke[service.IRecommendationService](i)}, nil
}

func (h *RecommendationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/recommendations", h.List)
}

func (h *RecommendationHandler) List(c *gin.Context) {
	userID, ok := requestUserID(c)
	if !ok {
		return
	}
	input, ok := recommendationRequest(c)
	if !ok {
		return
	}
	cmd := command.NewListRecommendationsCmd(userID, input.IllnessID, input.CabinetID)
	result, err := h.service.List(c.Request.Context(), cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set(middleware.ResponsePayloadKey, result)
	c.Set(middleware.ResponseStatusKey, http.StatusOK)
}
