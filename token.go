package caddycfg

// Token config token. Gives token location (Col is not supported currently)
type Token struct {
	File  string
	Value string
	Lin   int
	Col   int
}

// String ...
func (t Token) String() string {
	return t.Value
}
