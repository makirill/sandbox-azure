package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	error_not_found = "not found"
)

type AzureSandboxInstance struct {
	details SandboxDetails
}

type AzureSubscription struct {
	subscriptionID string
	instances      map[uuid.UUID]AzureSandboxInstance
	sync.RWMutex
}

func InitAzureSubscription(subscriptionID string) *AzureSubscription {
	return &AzureSubscription{
		subscriptionID: subscriptionID,
		//		instances:      []AzureSandboxInstance{},
		instances: make(map[uuid.UUID]AzureSandboxInstance),
	}
}

func (s *AzureSubscription) Add(name string) SandboxDetails {

	c := make(chan SandboxDetails)

	go func() {
		uuid := uuid.New()

		azureDetails := AzureSandboxInstance{
			details: SandboxDetails{
				Name:      name,
				UUID:      uuid.String(), // TODO: not sure if this is the correct type for UUID (like open-fs94 ??)
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // TODO: default expiration time is 7 days
				Status:    Pending,
			},
		}

		c <- azureDetails.details

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox creation

		azureDetails.details.Status = Running
		azureDetails.details.UpdatedAt = time.Now()

		s.Lock()
		s.instances[uuid] = azureDetails
		s.Unlock()
	}()

	return <-c
}

func (s *AzureSubscription) Remove(id string) (SandboxDetails, error) {

	c := make(chan SandboxDetails)
	e := make(chan error)

	go func() {
		uuid, err := uuid.Parse(id)
		if err != nil {
			e <- err
			return
		}

		s.Lock()
		azureSandboxInstance := s.instances[uuid]
		if azureSandboxInstance == (AzureSandboxInstance{}) {
			e <- errors.New(error_not_found)
			return
		}

		delete(s.instances, uuid)
		s.Unlock()

		c <- azureSandboxInstance.details

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox deletion

		azureSandboxInstance.details.Status = Deleted
		azureSandboxInstance.details.UpdatedAt = time.Now()

		s.Lock()
		s.instances[uuid] = azureSandboxInstance
		s.Unlock()

		fmt.Println("Sandbox deleted", azureSandboxInstance.details.UUID)
	}()

	select {
	case sandboxDetails := <-c:
		return sandboxDetails, nil
	case err := <-e:
		return SandboxDetails{}, err
	}

}

func (s *AzureSubscription) ListAll() []SandboxDetails {
	sandboxDetailsList := []SandboxDetails{}

	s.RLock()
	defer s.RUnlock()

	for _, sandboxInstance := range s.instances {
		sandboxDetailsList = append(sandboxDetailsList, sandboxInstance.details)
	}

	return sandboxDetailsList
}

func (s *AzureSubscription) GetByName(name string) []SandboxDetails {

	sandboxDetailsList := []SandboxDetails{}

	s.RLock()
	defer s.RUnlock()

	for _, sandboxInstance := range s.instances {
		if sandboxInstance.details.Name == name {
			sandboxDetailsList = append(sandboxDetailsList, sandboxInstance.details)
		}
	}

	return sandboxDetailsList
}

func (s *AzureSubscription) GetByUUID(id string) (SandboxDetails, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	s.RLock()
	defer s.RUnlock()

	azureSandboxInstance := s.instances[uuid]
	if azureSandboxInstance == (AzureSandboxInstance{}) {
		return SandboxDetails{}, errors.New(error_not_found)
	}

	return azureSandboxInstance.details, nil
}

func (s *AzureSubscription) UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	s.Lock()

	azureSandboxInstance := s.instances[uuid]

	if azureSandboxInstance == (AzureSandboxInstance{}) {
		return SandboxDetails{}, errors.New(error_not_found)
	}

	azureSandboxInstance.details.ExpiresAt = expiresAt
	azureSandboxInstance.details.UpdatedAt = time.Now()
	s.instances[uuid] = azureSandboxInstance

	s.Unlock()

	return azureSandboxInstance.details, nil
}
