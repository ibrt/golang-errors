package errorz

var (
	_ Option = Args{}
)

// Args describes a list of args used for formatting an error message.
type Args []interface{}

// Apply implements the Option interface.
func (a Args) Apply(_ error) {
	// intentionally empty
}

// A is a shorthand builder for args.
func A(a ...interface{}) Args {
	return a
}
