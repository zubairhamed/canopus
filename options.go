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

func IsRepeatableOption(opt *Option) bool {
	switch opt.Code {

	case OPTION_IF_MATCH, OPTION_ETAG, OPTION_URI_PORT, OPTION_LOCATION_PATH, OPTION_URI_PATH, OPTION_URI_QUERY, OPTION_LOCATION_QUERY,
		OPTION_BLOCK2, OPTION_BLOCK1:
		return true

	default:
		return false
	}
}

func IsValidOption(opt *Option) bool {
	switch opt.Code {

	case OPTION_IF_MATCH , OPTION_URI_HOST,
		 OPTION_ETAG, OPTION_IF_NONE_MATCH, OPTION_OBSERVE, OPTION_URI_PORT, OPTION_LOCATION_PATH,
		 OPTION_URI_PATH, OPTION_CONTENT_FORMAT, OPTION_MAX_AGE, OPTION_URI_QUERY, OPTION_ACCEPT,
		 OPTION_LOCATION_QUERY, OPTION_BLOCK2, OPTION_BLOCK1, OPTION_PROXY_URI, OPTION_PROXY_SCHEME, OPTION_SIZE1:
		return true

	default:
		return false
	}
}

func IsElectiveOption(opt *Option) bool {
	i := int(opt.Code)

	if (i & 1) == 1 {
		return false
	}
	return true
}

func IsCriticalOption(opt *Option) bool {
	return !IsElectiveOption(opt)
}
