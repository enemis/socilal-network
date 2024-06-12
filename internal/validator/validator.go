package validator

import (
	"fmt"

	"github.com/golodash/galidator"
)

type CustomValidatorRule interface {
	Rule() string
	Message() string
	Field() string
	IsSplittableField() bool
	Callback() func(interface{}) bool
}

func NewValidator(input any, externalValidators ...CustomValidatorRule) galidator.Validator {
	customValidators := make(map[string]func(interface{}) bool, 10)
	customMessages := galidator.Messages{}

	for _, rule := range externalValidators {
		ruleName := rule.Rule()
		if rule.IsSplittableField() {
			ruleName = fmt.Sprintf("%s_%s", rule.Field(), rule.Rule())
		}
		customMessages[ruleName] = rule.Message()
		customValidators[ruleName] = rule.Callback()
	}

	g := galidator.G().CustomMessages(customMessages).CustomValidators(customValidators)
	validator := g.Validator(input)

	return validator
}
