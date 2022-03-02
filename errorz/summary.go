package errorz

// Summary provides a serializable summary of an error and its metadata.
type Summary struct {
	ID         ID                     `json:"id,omitempty" yaml:"id,omitempty"`
	Status     Status                 `json:"statusCode,omitempty" yaml:"id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" yaml:"id,omitempty"`
	Message    string                 `json:"message,omitempty" yaml:"id,omitempty"`
	StackTrace []string               `json:"stackTrace,omitempty" yaml:"id,omitempty"`
}

// ToSummary converts an error to Summary.
func ToSummary(err error) *Summary {
	return &Summary{
		ID:         GetID(err),
		Status:     GetStatus(err),
		Metadata:   GetMetadata(err),
		Message:    err.Error(),
		StackTrace: FormatStackTrace(getCallersInternal(err, 1)),
	}
}
