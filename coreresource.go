package goap

type CoreAttributes map[string] string

func NewCoreResource() *CoreResource {
    c := &CoreResource{}

    c.Attributes = make(CoreAttributes)

    return c
}

type CoreResource struct {
    Target          string
    Attributes      CoreAttributes
}

func (c *CoreResource) AddAttribute(key string, value string) {
    c.Attributes[key] = value
}





