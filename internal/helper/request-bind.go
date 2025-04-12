package helper

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"scaf-gin/internal/core"
)

// BindJSON binds the JSON body to the provided request struct and handles validation errors.
func BindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errs := extractValidationErrors(req, validationErrors)
			return core.NewValidationError(errs)
		}
		return core.ErrBadRequest
	}
	return nil
}

// BindQuery binds the query parameters to the provided request struct and handles validation errors.
func BindQuery(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindQuery(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errs := extractValidationErrors(req, validationErrors)
			return core.NewValidationError(errs)
		}
		return core.ErrBadRequest
	}
	return nil
}

// BindUri binds the URI parameters to the provided request struct and handles validation errors.
func BindUri(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindUri(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errs := extractValidationErrors(req, validationErrors)
			return core.NewValidationError(errs)
		}
		return core.ErrBadRequest
	}
	return nil
}

func extractValidationErrors(req interface{}, verr validator.ValidationErrors) []map[string]interface{} {
	var errs []map[string]interface{}
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

		errs = append(errs, map[string]interface{}{
			"field":   jsonTag,
			"message": message,
			"tag":     fe.Tag(),
			"param":   fe.Param(),
		})
	}
	return errs
}