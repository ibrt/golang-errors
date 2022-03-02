package errorz

var (
	_ Option = Metadata{}
)

// Metadata describes metadata which can be attached to errors.
type Metadata map[string]interface{}

// Apply implements the Option interface.
func (m Metadata) Apply(err error) {
	if e, ok := err.(*wrappedError); ok {
		for k, v := range m {
			e.metadata[k] = v
		}
	}
}

// Get gets a value from Metadata, nil if not found or Metadata is nil.
func (m Metadata) Get(k string) interface{} {
	if m == nil {
		return nil
	}
	return m[k]
}

// GetString gets a value from Metadata as string, empty if not found or of different type.
func (m Metadata) GetString(k string) string {
	if s, ok := m.Get(k).(string); ok {
		return s
	}
	return ""
}

// M is a shorthand for providing metadata to errors.
func M(k string, v interface{}) OptionFunc {
	return func(err error) {
		if e, ok := err.(*wrappedError); ok {
			e.metadata[k] = v
		}
	}
}

// GetMetadata gets the metadata from the error, nil if not found.
func GetMetadata(err error) Metadata {
	if e, ok := err.(*wrappedError); ok {
		return e.metadata
	}
	return Metadata{}
}
