package rest

import (
	"fmt"
	"reflect"
	"strings"
)

// Name of the struct tag used in examples.
const tagName = "validate"

// Generic data validator.
type Validator interface {
	// Validate method performs validation and returns result and optional error.
	Validate(interface{}) (bool, error)
}

// DefaultValidator does not perform any validations.
type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

// StringValidator validates string presence and/or its length.
type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	input, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("should be string than %v", reflect.TypeOf(val))
	}

	l := len(input)

	if l == 0 {
		return false, fmt.Errorf("cannot be blank")
	}

	return true, nil
}

// NumberValidator performs numerical value validation.
// Its limited to int type for simplicity.
type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num, ok := val.(int)
	if !ok {
		return false, fmt.Errorf("should be int than %v", reflect.TypeOf(val))
	}

	if num < v.Min {
		return false, fmt.Errorf("should be greater than %v", v.Min)
	}

	if v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("should be less than %v", v.Max)
	}

	return true, nil
}

// Returns validator struct corresponding to validation type
func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	}

	return DefaultValidator{}
}

// Performs actual data validation using validator definitions on the struct
func ValidateStruct(s interface{}) error {
	// ValueOf returns a Value representing the run-time data
	v := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value

		tag := v.Type().Field(i).Tag.Get(tagName)
		// Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		// Get a validator that corresponds to a tag
		validator := getValidatorFromTag(tag)
		// Perform validation
		valid, err := validator.Validate(v.Field(i).Interface())

		// Append error to results
		if !valid && err != nil {
			return fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error())
		}
	}

	return nil
}
