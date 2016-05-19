package canopus

import (
	"math"
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

func NewBlock1Option(bs BlockSize, more bool) *Block1Option {
	opt := &Block1Option{}

	val := uint32(bs)

	if more {
		val |= (1 << 3)
	}

	opt.Value = val

	/*
		BLockSize
		val := o.Value.(uint32)
		exp := val & 0x07

		return math.Exp2(float64(exp + 4))

		More
		val := o.Value.(uint32)

		return ((val >> 3) & 0x01) == 1

	*/

	return opt
}

func Block1OptionFromOption(opt *Option) *Block1Option {
	blockOpt := &Block1Option{}

	blockOpt.Value = opt.Value
	blockOpt.Code = opt.Code

	return blockOpt
}

type Block1Option struct {
	Option
}

func (o *Block1Option) Sequence() uint32 {
	val := o.Value.(uint32)

	return val >> 4
}

func (o *Block1Option) Exponent() uint32 {
	val := o.Value.(uint32)

	return val & 0x07
}

func (o *Block1Option) Size() float64 {
	val := o.Value.(uint32)
	exp := val & 0x07

	return math.Exp2(float64(exp + 4))
}

func (o *Block1Option) HasMore() bool {
	val := o.Value.(uint32)

	return ((val >> 3) & 0x01) == 1
}
