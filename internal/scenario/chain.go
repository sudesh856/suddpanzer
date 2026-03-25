package scenario

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ChainStore struct {
	data map[string]string
}

func NewChainStore() *ChainStore {
	return &ChainStore{data: make(map[string]string)}
}

func (c *ChainStore) Store(endpointName string, body []byte, extract map[string]string) {
	if len(extract) == 0{
		return
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return
	}

	for varName, jsonKey := range extract {
		key := strings.TrimPrefix(jsonKey, "$.")
		if val, ok := parsed[key]; ok {
			c.data[endpointName+"."+varName] = fmt.Sprintf("%v", val)
		}
	}
}

func (c *ChainStore) Get(endpointName, varName string) (string, bool) {
	val, ok := c.data[endpointName+"."+varName]
	return val, ok
}

func (c *ChainStore) ToVars() map[string]string {
	out := make(map[string]string)
	for k,v := range c.data {
		out[k] = v
	}
	return out
}
