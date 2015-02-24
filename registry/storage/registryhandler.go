package storage

import (
	"fmt"

	"github.com/docker/distribution"
)

type RegistryHandlerInitFunc func(registry distribution.Registry, options map[string]interface{}) (distribution.Registry, error)

var registryHandlers map[string]RegistryHandlerInitFunc

func RegisterRegistryHandler(name string, initFunc RegistryHandlerInitFunc) error {
	if registryHandlers == nil {
		registryHandlers = make(map[string]RegistryHandlerInitFunc)
	}
	if _, exists := registryHandlers[name]; exists {
		return fmt.Errorf("name already registered: %s", name)
	}

	registryHandlers[name] = initFunc

	return nil
}

func GetRegistryHandler(name string, registry distribution.Registry, options map[string]interface{}) (distribution.Registry, error) {
	if registryHandlers != nil {
		if initFunc, exists := registryHandlers[name]; exists {
			return initFunc(registry, options)
		}
	}

	return nil, fmt.Errorf("no registry handler registered with name: %s", name)
}
