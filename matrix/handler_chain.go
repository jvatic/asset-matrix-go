package matrix

type HandlerChain struct {
	Handler       Handler
	embededChains []*HandlerChain
}

func NewHandlerChain(file *File, handlers []HandlerConstructor) *HandlerChain {
	chain := &HandlerChain{embededChains: make([]*HandlerChain, 0)}
	chain.evaluate(file, handlers)

	return chain
}

func (chain *HandlerChain) AddHandlerChain(otherChain *HandlerChain) {
	if !chain.containsEmbededChain(otherChain) {
		chain.embededChains = append(chain.embededChains, otherChain)
	}
}

func (chain *HandlerChain) Evaluate() {
	for _, directive := range chain.Handler.InputFile().Directives {
		for _, file := range directive.Files() {
			if file.HandlerChain != nil {
				chain.AddHandlerChain(file.HandlerChain)
			}
		}
	}
}

func (chain *HandlerChain) evaluate(file *File, handlers []HandlerConstructor) {
	for _, handlerConstructor := range handlers {
		if handler, canHandle := handlerConstructor(file); canHandle {
			chain.Handler = handler
			break
		}
	}

	if chain.Handler == nil {
		handler, _ := NewDefaultHandler(file)
		chain.Handler = handler
	}
}

func (chain *HandlerChain) containsEmbededChain(otherChain *HandlerChain) bool {
	for _, c := range chain.embededChains {
		if c == otherChain {
			return true
		}
	}
	return false
}
