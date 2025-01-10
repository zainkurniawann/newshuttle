package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CustomPhoneValidator(fl validator.FieldLevel) bool {
	phoneRegex := `^(\+?)([0-9]{12,15})$`
	value := fl.Field().String()
	return regexp.MustCompile(phoneRegex).MatchString(value)
}

func CustomUsernameValidator(fl validator.FieldLevel) bool {
	usernameRegex := `^[a-zA-Z0-9_-]{3,}$`
	value := fl.Field().String()
	return regexp.MustCompile(usernameRegex).MatchString(value)
}

func CustomRoleValidator(fl validator.FieldLevel) bool {
	roleRegex := `^(superadmin|schooladmin|driver|parent|Superadmin|Schooladmin|Driver|Parent)$`
	value := fl.Field().String()
	return regexp.MustCompile(roleRegex).MatchString(value)
}

func CustomGenderValidator(fl validator.FieldLevel) bool {
	genderRegex := `^(male|female|Male|Female)$`
	value := fl.Field().String()
	return regexp.MustCompile(genderRegex).MatchString(value)
}

func ValidateStruct(c *fiber.Ctx, v interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("phone", CustomPhoneValidator)
	validate.RegisterValidation("username", CustomUsernameValidator)
	validate.RegisterValidation("role", CustomRoleValidator)
	validate.RegisterValidation("gender", CustomGenderValidator)

	if err := validate.Struct(v); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return fmt.Errorf("the %s field is required", err.Field())
			case "username":
				return fmt.Errorf("the %s field must be at least 3 characters and can only contain letters, numbers, underscores, and hyphens", err.Field())
			case "gender":
				return fmt.Errorf("the %s field gender must be either male or female", err.Field())
			case "email":
				return fmt.Errorf("the %s field must be a valid email address", err.Field())
			case "phone":
				return fmt.Errorf("the %s field must be a valid phone number", err.Field())
			case "min":
				return fmt.Errorf("the %s field must be at least %s characters", err.Field(), err.Param())
			case "max":
				return fmt.Errorf("the %s field must be at most %s characters", err.Field(), err.Param())
			case "role":
				return fmt.Errorf("the %s field must be either superadmin, schooladmin, driver, or parent", err.Field())
			}
		}
	}
	return nil
}