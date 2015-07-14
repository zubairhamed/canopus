package canopus

import (
	"strings"
)

/* Option */
type Option struct {
	Code  OptionCode
	Value interface{}
}

func (o *Option) Name() string {
	return "Name of option"
}

func (o *Option) IsElective() bool {
	if (int(o.Code) % 2) != 0 {
		return false
	}
	return true
}

func (o *Option) IsCritical() bool {
	if (int(o.Code) % 2) != 0 {
		return true
	}
	return false
}

func (o *Option) StringValue() string {
	return o.Value.(string)
}

////////////////////////////////////////

func NewOption(optionNumber OptionCode, optionValue interface{}) *Option {
	return &Option{
		Code:  optionNumber,
		Value: optionValue,
	}
}

func NewPathOptions(path string) []*Option {
	opts := []*Option{}
	ps := strings.Split(path, "/")

	for _, p := range ps {
		if p != "" {
			opt := NewOption(OPTION_URI_PATH, p)
			opts = append(opts, opt)
		}
	}
	return opts
}

func RepeatableOption(opt *Option) bool {
	switch opt.Code {

	case OPTION_IF_MATCH, OPTION_ETAG, OPTION_URI_PORT, OPTION_LOCATION_PATH, OPTION_URI_PATH, OPTION_URI_QUERY, OPTION_LOCATION_QUERY,
		OPTION_BLOCK2, OPTION_BLOCK1:
		return true

	default:
		return false
	}
}
