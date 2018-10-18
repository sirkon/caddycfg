package caddycfg

type argumentAccess interface {
	appendData(items []string)
	Arguments() []string
}

// Args provides access to positional parameters, this is to be used to consume arguments that predate a block
// I mean the following; let we have a config
//     head arg₁ arg₂ … argₙ {
//         …
//     }
// Head is used (in case of structs) to choose a field with its type. In case if there are arguments, i.e. n > 0
// followed by block (there can be no block, in this case regular []int, []string, etc is enough) the type must
// implement argumentAccess – this is equivalent for having type Args embedded or being type Args itself. Although
// it is possible to use Args itself, it will not work well enough. So, use it for embedding into your own types
type Args struct {
	data []string
}

func (a *Args) appendData(items []string) {
	if len(items) == 0 {
		a.data = nil
		return
	}
	a.data = make([]string, len(items))
	copy(a.data, items)
}

func (a *Args) Arguments() []string {
	return a.data
}

// ArgumentCollector special interface whose implementations can be used for taking arguments with additional control
// over the content, e.g. they can keep context to provide valuable error diagnostic
type ArgumentCollector interface {
	// AppendArgument is taking arguments in a right order from left to right. It is user responsibility to incorporate
	// positional information into error message. Helpers functions caddycfg.TokenError and caddycfg.TokenErrorf can do
	// this
	AppendArgument(arg Token) error
	// Arguments is considered to return collected arguments
	Arguments() []string
}
