package canopus

import (
	"strings"
)

// Represents an Option for a CoAP Message
type Option struct {
	Code  OptionCode
	Value interface{}
}

func (o *Option) Name() string {
	return "Name of option"
}

// Determines if an option is elective
func (o *Option) IsElective() bool {
	if (int(o.Code) % 2) != 0 {
		return false
	}
	return true
}

// Determines if an option is critical
func (o *Option) IsCritical() bool {
	if (int(o.Code) % 2) != 0 {
		return true
	}
	return false
}

// Returns the string value of an option
func (o *Option) StringValue() string {
	return o.Value.(string)
}

func (o *Option) IntValue() int {
	return o.Value.(int)
}

// Instantiates a New Option
func NewOption(optionNumber OptionCode, optionValue interface{}) *Option {
	return &Option{
		Code:  optionNumber,
		Value: optionValue,
	}
}

// Creates an array of options decomposed from a given path
func NewPathOptions(path string) []*Option {
	opts := []*Option{}
	ps := strings.Split(path, "/")

	for _, p := range ps {
		if p != "" {
			opt := NewOption(OptionURIPath, p)
			opts = append(opts, opt)
		}
	}
	return opts
}

// Checks if an option is repeatable
func IsRepeatableOption(opt *Option) bool {
	switch opt.Code {

	case OptionIfMatch, OptionEtag, OptionURIPort, OptionLocationPath, OptionURIPath, OptionURIQuery, OptionLocationQuery,
		OptionBlock2, OptionBlock1:
		return true

	default:
		return false
	}
}

// Checks if an option/option code is recognizable/valid
func IsValidOption(opt *Option) bool {
	switch opt.Code {

	case OptionIfNoneMatch, OptionURIHost,
		OptionEtag, OptionIfMatch, OptionObserve, OptionURIPort, OptionLocationPath,
		OptionURIPath, OptionContentFormat, OptionMaxAge, OptionURIQuery, OptionAccept,
		OptionLocationQuery, OptionBlock2, OptionBlock1, OptionProxyURI, OptionProxyScheme, OptionSize1:
		return true

	default:
		return false
	}
}

// Determines if an option is elective
func IsElectiveOption(opt *Option) bool {
	i := int(opt.Code)

	if (i & 1) == 1 {
		return false
	}
	return true
}

// Determines if an option is critical
func IsCriticalOption(opt *Option) bool {
	return !IsElectiveOption(opt)
}
