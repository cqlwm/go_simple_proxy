package module

type BaseModule interface {
	ModuleName() string
	Init(parameter map[string]interface{}) (err error)
}

var AllModule = map[string]BaseModule{}

func t() {
	ha := AllModule["HttpAuthenticator"]
	m  := map[string]interface {}{}
	ha.Init(m)
}
