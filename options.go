package goap

import "strings"

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
		opt := NewOption(OPTION_URI_PATH, p)
		opts = append(opts, opt)
	}

	return opts
}

/*
	OPTION_IF_MATCH       OptionCode = 1
	OPTION_URI_HOST       OptionCode = 3
	OPTION_ETAG           OptionCode = 4
	OPTION_IF_NONE_MATCH  OptionCode = 5
	OPTION_URI_PORT       OptionCode = 7
	OPTION_LOCATION_PATH  OptionCode = 8
	OPTION_URI_PATH       OptionCode = 11
	OPTION_CONTENT_FORMAT OptionCode = 12
	OPTION_MAX_AGE        OptionCode = 14
	OPTION_URI_QUERY      OptionCode = 15
	OPTION_ACCEPT         OptionCode = 17
	OPTION_LOCATION_QUERY OptionCode = 20
	OPTION_PROXY_URI      OptionCode = 35
	OPTION_PROXY_SCHEME   OptionCode = 39
	OPTION_SIZE1          OptionCode = 60
 */

