package models

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type AzureSandboxList struct {
	sandboxInstances map[uuid.UUID]SandboxDetails
	sync.RWMutex
}

func InitAzureSandboxList() *AzureSandboxList {
	return &AzureSandboxList{
		sandboxInstances: make(map[uuid.UUID]SandboxDetails),
	}
}

func (s *AzureSandboxList) Put(sandbox SandboxDetails) error {

	uuid, err := uuid.Parse(sandbox.UUID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.sandboxInstances[uuid] = sandbox

	return nil
}

func (s *AzureSandboxList) Remove(id string) (SandboxDetails, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	s.Lock()
	defer s.Unlock()

	sandbox, ok := s.sandboxInstances[uuid]
	if !ok {
		return SandboxDetails{}, errors.New(error_not_found)
	}

	delete(s.sandboxInstances, uuid)

	return sandbox, nil
}

func (s *AzureSandboxList) ListAll() []SandboxDetails {

	s.RLock()
	defer s.RUnlock()

	sandboxes := []SandboxDetails{}

	for _, sandbox := range s.sandboxInstances {
		sandboxes = append(sandboxes, sandbox)
	}

	return sandboxes
}

func (s *AzureSandboxList) GetByName(name string) []SandboxDetails {

	sandboxes := []SandboxDetails{}

	s.RLock()
	defer s.RUnlock()

	for _, sandbox := range s.sandboxInstances {
		if sandbox.Name == name {
			sandboxes = append(sandboxes, sandbox)
		}
	}

	return sandboxes
}

func (s *AzureSandboxList) GetByUUID(id string) (SandboxDetails, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	s.RLock()
	defer s.RUnlock()

	sandbox, ok := s.sandboxInstances[uuid]
	if !ok {
		return SandboxDetails{}, errors.New(error_not_found)
	}

	return sandbox, nil
}
