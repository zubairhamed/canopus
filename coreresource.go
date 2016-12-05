package canopus

// Instantiates a new core-attribute with a given key/value
func NewCoreAttribute(key string, value interface{}) *CoreAttribute {
	return &CoreAttribute{
		Key:   key,
		Value: value,
	}
}

// Instantiates a new Core Resource Object
func NewCoreResource() *CoreResource {
	c := &CoreResource{}

	return c
}
