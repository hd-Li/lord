package v1alpha1

type StringMatch struct {
	// Specified exactly one of the fields below.

	// exact string match
	Exact string `json:"exact,omitempty"`

	// prefix-based match
	Prefix string `json:"prefix,omitempty"`

	// suffix-based match.
	Suffix string `json:"suffix,omitempty"`

	// ECMAscript style regex-based match
	Regex string `json:"regex,omitempty"`
}