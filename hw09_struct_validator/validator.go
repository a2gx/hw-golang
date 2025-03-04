package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var str strings.Builder
	for _, err := range v {
		str.WriteString(fmt.Sprintf("%s: %v\n", err.Field, err.Err))
	}
	return str.String()
}

type InternalError struct {
	Err error
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error: %v", e.Err)
}

func Validate(v interface{}) error {
	valueOf, typeOf := reflect.ValueOf(v), reflect.TypeOf(v)

	if valueOf.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct")
	}

	var result ValidationErrors
	for i := 0; i < typeOf.NumField(); i++ {
		fieldValue, fieldType := valueOf.Field(i), typeOf.Field(i)
		validateTag := fieldType.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, "|")
		for _, rule := range rules {
			if err := validCheckField(fieldValue, rule); err != nil {
				var internalError InternalError
				if errors.As(err, &internalError) {
					return err
				}
				result = append(result, ValidationError{
					Field: fieldType.Name,
					Err:   err,
				})
			}
		}

	}

	if len(result) > 0 {
		return result
	}

	return nil
}

// ... check the rules

var mapRulesFn = map[string]func(field reflect.Value, value string) error{
	"len": func(field reflect.Value, value string) error {
		if field.Kind() == reflect.String {
			check, err := strconv.Atoi(value)
			if err != nil {
				return InternalError{err}
			}
			if len(field.String()) != check {
				return fmt.Errorf("length must be exactly %d", check)
			}
		}
		return nil
	},
}

func validCheckField(field reflect.Value, rule string) error {
	if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			if err := validCheckRules(field.Index(i), rule); err != nil {
				return err
			}
		}
		return nil
	}

	return validCheckRules(field, rule)
}

func validCheckRules(field reflect.Value, rule string) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return InternalError{fmt.Errorf("invalid rule format: %s", rule)}
	}

	name, value := parts[0], parts[1]
	rulesFn, ok := mapRulesFn[name]
	if !ok {
		return InternalError{fmt.Errorf("unknown validation rule: %s", name)}
	}

	return rulesFn(field, value)
}
