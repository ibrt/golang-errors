package errorz

var (
	_ Option = ID("")
)

// ID describes an error id.
type ID string

// String implements the fmt.Stringer interface.
func (id ID) String() string {
	return string(id)
}

// Apply implements the Option interface.
func (id ID) Apply(err error) {
	if e, ok := err.(*wrappedError); ok {
		e.id = id
	}
}

// GetID gets the id from the error, or an empty id if not set.
func GetID(err error) ID {
	if e, ok := err.(*wrappedError); ok {
		return e.id
	}
	return ""
}
