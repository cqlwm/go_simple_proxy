package module

type BaseModule interface {
	ModuleName() string
	Init(parameter map[string]interface{}) (err error)
}

var AllModule = map[string]BaseModule{}
