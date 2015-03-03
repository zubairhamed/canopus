package goap

/* Option */
type Option struct {
    Code    OptionCode
    Value   interface{}
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

func NewOption(optionNumber OptionCode, optionValue interface{}) *Option{
    return &Option{
        Code: optionNumber,
        Value: optionValue,
    }
}
