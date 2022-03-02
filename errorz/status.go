package errorz

var (
	_ Option = Status(0)
)

// Status describes an HTTP status code.
type Status int

// Int returns the status code as "int".
func (s Status) Int() int {
	return int(s)
}

// Apply implements the Option interface.
func (s Status) Apply(err error) {
	if e, ok := err.(*wrappedError); ok {
		e.status = s
	}
}

// GetStatus gets the status code from the error, or 0 if not set.
func GetStatus(err error) Status {
	if e, ok := err.(*wrappedError); ok {
		return e.status
	}
	return 0
}
