package caddycfg

// Validator unmarshal will call method Err of input value if input type implements this interface
type Validator interface {
	Err(head Token) error
}
