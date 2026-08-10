package handler

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"scaf-gin/internal/core"
)

func BindJSON(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return core.NewValidationError(extractValidationErrors(req, validationErrors))
		}
		return core.ErrBadRequest
	}
	return nil
}

func BindQuery(c *gin.Context, req any) error {
	if err := c.ShouldBindQuery(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return core.NewValidationError(extractValidationErrors(req, validationErrors))
		}
		return core.ErrBadRequest
	}
	return nil
}

func BindURI(c *gin.Context, req any) error {
	if err := c.ShouldBindUri(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return core.NewValidationError(extractValidationErrors(req, validationErrors))
		}
		return core.ErrBadRequest
	}
	return nil
}

func extractValidationErrors(req any, verr validator.ValidationErrors) []map[string]any {
	errs := make([]map[string]any, 0, len(verr))
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for _, fe := range verr {
		field, _ := t.FieldByName(fe.Field())
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fe.Field()
		}
		message := "Invalid value"
		if fe.Tag() == "required" {
			message = "Required field"
		}

		errs = append(errs, map[string]any{
			"field":   jsonTag,
			"message": message,
			"tag":     fe.Tag(),
			"param":   fe.Param(),
		})
	}
	return errs
}
