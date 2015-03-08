package goap

func NewCoreAttribute(key string, value interface{}) *CoreAttribute {
	return &CoreAttribute{
		Key:   key,
		Value: value,
	}
}

type CoreAttribute struct {
	Key   string
	Value interface{}
}

func NewCoreResource() *CoreResource {
	c := &CoreResource{}

	return c
}

type CoreAttributes []*CoreAttribute
type CoreResource struct {
	Target     string
	Attributes CoreAttributes
}

func (c *CoreResource) AddAttribute(key string, value interface{}) {
	c.Attributes = append(c.Attributes, NewCoreAttribute(key, value))
}
