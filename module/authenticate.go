package module

import "fmt"

const HttpAuthenticatorModule string = "HttpAuthenticatorModule"

type Authenticator interface {
	Authenticate(credentials string) bool
}

type HttpAuthenticatorConfig struct {
	url    string
	method string
	header map[string][]string
	query  map[string]string
	body   interface{}
}

func (c *HttpAuthenticatorConfig) Load(parameter map[string]interface{}) error {
	var ok bool
	if c.url, ok = parameter["url"].(string); !ok {
		return fmt.Errorf("url must be a string")
	}
	if c.method, ok = parameter["method"].(string); !ok {
		return fmt.Errorf("method must be a string")
	}
	if c.header, ok = parameter["header"].(map[string][]string); !ok {
		return fmt.Errorf("header must be a map[string][]string")
	}
	if c.query, ok = parameter["query"].(map[string]string); !ok {
		return fmt.Errorf("query must be a map[string]string")
	}
	c.body = parameter["body"]
	return nil
}

type HttpAuthenticator struct {
	name string
	httpConfig  HttpAuthenticatorConfig
}

func (h *HttpAuthenticator) ModuleName() string {
	return h.name
}

func (h *HttpAuthenticator) Init(parameter map[string]interface{}) (err error) {
	h.name = HttpAuthenticatorModule
	h.httpConfig = HttpAuthenticatorConfig{}
	err = h.httpConfig.Load(parameter)


	return err
}

func (h *HttpAuthenticator) Authenticate(credentials string) bool {

	return false
}
