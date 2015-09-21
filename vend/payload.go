package vend

// Payload contains resource data and versioning info.
// This is the default format returned by 2.0 endpoints.
type Payload struct {
	Data    map[string]Body  `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// Body ...
type Body struct {
	Data []byte
}
