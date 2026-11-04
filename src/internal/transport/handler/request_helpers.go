package handler

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/middleware"
)

func requestUserID(c *gin.Context) (uuid.UUID, bool) {
	value, ok := c.Get(middleware.UserIDKey)
	if !ok {
		_ = c.Error(apperrors.ErrUnauthorized.WithSource())
		return uuid.Nil, false
	}
	id, ok := value.(uuid.UUID)
	if !ok {
		_ = c.Error(apperrors.ErrUnauthorized.WithSource())
		return uuid.Nil, false
	}
	return id, true
}

func requestUUID(c *gin.Context, raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		_ = c.Error(apperrors.ErrInvalidInput.WithError(err))
		return uuid.Nil, false
	}
	return id, true
}

func pathIDRequest(c *gin.Context, name string) (req.IDDto, bool) {
	id, ok := requestUUID(c, c.Param(name))
	return req.IDDto{ID: id}, ok
}

func pageParams(c *gin.Context) (int, int, bool) {
	limit, offset := 20, 0
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 || parsed > 100 {
			_ = c.Error(apperrors.ErrInvalidInput.WithSource())
			return 0, 0, false
		}
		limit = parsed
	}
	if raw := c.Query("offset"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			_ = c.Error(apperrors.ErrInvalidInput.WithSource())
			return 0, 0, false
		}
		offset = parsed
	}
	return limit, offset, true
}

func decodeBody(c *gin.Context, dest any) bool {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		_ = c.Error(apperrors.ErrInvalidInput.WithError(err))
		return false
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		_ = c.Error(apperrors.ErrInvalidInput.WithSource())
		return false
	}
	return true
}

func decodeFields(body map[string]json.RawMessage, dest any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func bindBody(c *gin.Context, dest any) bool {
	var body map[string]json.RawMessage
	if !decodeBody(c, &body) {
		return false
	}
	if len(body) == 0 {
		_ = c.Error(apperrors.ErrInvalidInput.WithSource())
		return false
	}
	typ := reflect.TypeOf(dest).Elem()
	allowed := make(map[string]bool, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		allowed[name] = true
		if field.Tag.Get("required") == "true" {
			if _, exists := body[name]; !exists {
				_ = c.Error(apperrors.ErrInvalidInput.WithSource())
				return false
			}
		}
	}
	for name, value := range body {
		if !allowed[name] || string(value) == "null" {
			_ = c.Error(apperrors.ErrInvalidInput.WithSource())
			return false
		}
	}
	if err := decodeFields(body, dest); err != nil {
		_ = c.Error(apperrors.ErrInvalidInput.WithError(err))
		return false
	}
	return true
}
