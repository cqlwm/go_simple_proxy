package module

import (
	"encoding/json"
	"fmt"
	"http_forwarder_go/util"
)

const HttpAuthenticatorModule string = "HttpAuthenticatorModule"

type Authenticator interface {
	Authenticate(credentials string) bool
}

type HttpAuthenticatorConfig struct {
	url            string
	method         string
	header         map[string][]string
	query          map[string]string
	body           interface{}
	varsAssign     []string
	checkCondition []string
	connectTimeout int32
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
	name       string
	httpConfig HttpAuthenticatorConfig
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
	cf := h.httpConfig
	var requestBody []byte
	if cf.body == nil {
		requestBody = []byte{}
	} else {
		bs, err := json.Marshal(cf.body)
		if err != nil {
			panic("The body cannot be serialized")
		}
		requestBody = bs
	}
	response := util.DoRequest(cf.method, cf.url, cf.header, requestBody)

	// todo varsAssign && checkCondition parse
	return response.State == 200
}
