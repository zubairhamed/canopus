package canopus

type CoreAttributes []*CoreAttribute

type CoreResource struct {
	Target     string
	Attributes CoreAttributes
}

type CoreAttribute struct {
	Key   string
	Value interface{}
}

// Adds an attribute (key/value) for a given core resource
func (c *CoreResource) AddAttribute(key string, value interface{}) {
	c.Attributes = append(c.Attributes, NewCoreAttribute(key, value))
}

// Gets an attribute for a core resource
func (c *CoreResource) GetAttribute(key string) *CoreAttribute {
	for _, attr := range c.Attributes {
		if attr.Key == key {
			return attr
		}
	}
	return nil
}
