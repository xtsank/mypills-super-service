package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

func pageRequest(c *gin.Context) (req.PageDto, bool) {
	limit, offset, ok := pageParams(c)
	if !ok {
		return req.PageDto{}, false
	}
	input := req.PageDto{Limit: limit, Offset: offset}
	if raw := c.Query("id"); raw != "" {
		id, valid := requestUUID(c, raw)
		if !valid {
			return req.PageDto{}, false
		}
		input.ID = &id
	}
	return input, true
}

func medicineListRequest(c *gin.Context) (req.MedicineListDto, bool) {
	page, ok := pageRequest(c)
	if !ok {
		return req.MedicineListDto{}, false
	}
	input := req.MedicineListDto{ID: page.ID, Limit: page.Limit, Offset: page.Offset}
	if name, exists := c.GetQuery("q"); exists {
		input.Name = &name
	}
	return input, true
}

func referenceListRequest(c *gin.Context) (req.ReferenceListDto, bool) {
	limit, offset, ok := pageParams(c)
	return req.ReferenceListDto{Kind: c.Param("kind"), Limit: limit, Offset: offset}, ok
}

func recommendationRequest(c *gin.Context) (req.RecommendationListDto, bool) {
	if c.Query("illness_id") == "" {
		_ = c.Error(apperrors.ErrInvalidInput.WithSource())
		return req.RecommendationListDto{}, false
	}
	illnessID, ok := requestUUID(c, c.Query("illness_id"))
	if !ok {
		return req.RecommendationListDto{}, false
	}
	input := req.RecommendationListDto{IllnessID: illnessID}
	if raw := c.Query("cabinet_id"); raw != "" {
		cabinetID, valid := requestUUID(c, raw)
		if !valid {
			return req.RecommendationListDto{}, false
		}
		input.CabinetID = &cabinetID
	}
	return input, true
}

func cabinetItemsListRequest(c *gin.Context) (req.CabinetItemsListDto, bool) {
	page, ok := pageRequest(c)
	if !ok {
		return req.CabinetItemsListDto{}, false
	}
	input := req.CabinetItemsListDto{ID: page.ID, Limit: page.Limit, Offset: page.Offset}
	for _, field := range []string{"low_stock", "expired"} {
		if raw, exists := c.GetQuery(field); exists {
			value, err := strconv.ParseBool(raw)
			if err != nil {
				_ = c.Error(apperrors.ErrInvalidInput.WithError(err))
				return req.CabinetItemsListDto{}, false
			}
			if field == "low_stock" {
				input.LowStock = value
			} else {
				input.Expired = value
			}
		}
	}
	if raw, exists := c.GetQuery("expires_within_days"); exists {
		days, err := strconv.Atoi(raw)
		if err != nil {
			_ = c.Error(apperrors.ErrInvalidInput.WithError(err))
			return req.CabinetItemsListDto{}, false
		}
		input.ExpiresWithinDays = &days
	}
	return input, true
}
