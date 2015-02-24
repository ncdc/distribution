package storage

import (
	"fmt"

	"github.com/docker/distribution"
	"github.com/docker/distribution/registry/storage/driver"
)

type NewManifestServiceFunc func(name string, options map[string]interface{}, storageDriver driver.StorageDriver) distribution.ManifestService

var manifestHandlers map[string]NewManifestServiceFunc

func RegisterManifestHandler(name string, initFunc NewManifestServiceFunc) error {
	if manifestHandlers == nil {
		manifestHandlers = make(map[string]NewManifestServiceFunc)
	}
	if _, exists := manifestHandlers[name]; exists {
		return fmt.Errorf("name already registered: %s", name)
	}

	manifestHandlers[name] = initFunc

	return nil
}

type ManifestServiceCreator struct {
	options       map[string]interface{}
	storageDriver driver.StorageDriver
	initFunc      NewManifestServiceFunc
}

func (c ManifestServiceCreator) CreateManifestService(name string) distribution.ManifestService {
	return c.initFunc(name, c.options, c.storageDriver)
}

func GetManifestHandler(name string, options map[string]interface{}, storageDriver driver.StorageDriver) (*ManifestServiceCreator, error) {
	if manifestHandlers != nil {
		if initFunc, exists := manifestHandlers[name]; exists {
			return &ManifestServiceCreator{options, storageDriver, initFunc}, nil
		}
	}

	return nil, fmt.Errorf("no manifest handler registered with name: %s", name)
}
