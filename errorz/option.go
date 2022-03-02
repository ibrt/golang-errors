package errorz

var (
	_ Option = OptionFunc(nil)
)

// Option describes an option which can be applied to an error.
type Option interface {
	Apply(err error)
}

// OptionFunc describes an option which can be applied to an error.
type OptionFunc func(err error)

// Apply implements the Option interface.
func (f OptionFunc) Apply(err error) {
	f(err)
}
