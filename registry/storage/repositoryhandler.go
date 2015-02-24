package storage

import (
	"fmt"

	"github.com/docker/distribution"
)

type RepositoryHandlerInitFunc func(repository distribution.Repository, options map[string]interface{}) (distribution.Repository, error)

var repositoryHandlers map[string]RepositoryHandlerInitFunc

func RegisterRepositoryHandler(name string, initFunc RepositoryHandlerInitFunc) error {
	if repositoryHandlers == nil {
		repositoryHandlers = make(map[string]RepositoryHandlerInitFunc)
	}
	if _, exists := repositoryHandlers[name]; exists {
		return fmt.Errorf("name already registered: %s", name)
	}

	repositoryHandlers[name] = initFunc

	return nil
}

func GetRepositoryHandler(name string, repository distribution.Repository, options map[string]interface{}) (distribution.Repository, error) {
	if repositoryHandlers != nil {
		if initFunc, exists := repositoryHandlers[name]; exists {
			return initFunc(repository, options)
		}
	}

	return nil, fmt.Errorf("no repository handler registered with name: %s", name)
}
