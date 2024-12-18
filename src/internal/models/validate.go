package models

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"log"
)

type CustomValidation struct {
	TagName    string
	Validation func(fl validator.FieldLevel) bool
}

type BaseValidator struct{}

func (r *BaseValidator) Validate(model interface{}, validations ...CustomValidation) error {
	validate := validator.New()

	if err := r.RegisterValidate(validate, validations...); err != nil {
		log.Fatalln(err)
	}

	err := validate.Struct(model)
	if err != nil {
		return err
	}

	return nil
}

func (r *BaseValidator) RegisterValidate(validate *validator.Validate, validations ...CustomValidation) error {
	for _, customValidation := range validations {
		err := validate.RegisterValidation(customValidation.TagName, customValidation.Validation)
		if err != nil {
			return fmt.Errorf("failed to register validation for tag '%s': %w", customValidation.TagName, err)
		}
	}
	return nil
}

func (r *BaseValidator) MapToModel(requestModel interface{}, mainModel interface{}) error {
	if err := copier.CopyWithOption(mainModel, requestModel, copier.Option{IgnoreEmpty: true}); err != nil {
		return err
	}
	return nil
}
