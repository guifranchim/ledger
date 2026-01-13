package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, validate *validator.Validate, data interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, "error decoding json: "+err.Error())
		return true
	}
	if err := validate.Struct(data); err != nil {
		ValidationErrorResponse(w, r, err)
		return true
	}

	return false
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{
		"error": message,
	})
}

func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)

	validationErrors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				validationErrors[field] = "This field is required"
			case "min":
				validationErrors[field] = "Minimum value: " + e.Param()
			case "max":
				validationErrors[field] = "Maximum value: " + e.Param()
			case "gte":
				validationErrors[field] = "Must be greater than or equal to " + e.Param()
			case "gt":
				validationErrors[field] = "Must be greater than " + e.Param()
			case "email":
				validationErrors[field] = "Invalid email"
			default:
				validationErrors[field] = "Validation failed: " + e.Tag()
			}
		}
	}

	render.JSON(w, r, map[string]interface{}{
		"error":  "Validation failed",
		"fields": validationErrors,
	})
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}
