package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	default_expiration_time = time.Hour * 24 * 7 // 7 days
)

type AzureSandbox struct {
	subscriptionID string
	instances      SandboxList
}

func InitAzureSandbox(subscriptionID string) *AzureSandbox {
	return &AzureSandbox{
		subscriptionID: subscriptionID,
		instances:      InitAzureSandboxList(),
	}
}

func (s *AzureSandbox) Create(name string) SandboxDetails {

	c := make(chan SandboxDetails)

	go func() {
		details := SandboxDetails{
			Name:      name,
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(default_expiration_time),
			Status:    Pending,
		}

		s.instances.Put(details)

		c <- details

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox creation

		details.Status = Running
		details.UpdatedAt = time.Now()

		s.instances.Put(details)
	}()

	return <-c
}

func (s *AzureSandbox) Remove(id string) (SandboxDetails, error) {
	c := make(chan SandboxDetails)
	e := make(chan error)

	go func() {
		details, err := s.instances.GetByUUID(id)
		if err != nil {
			e <- err
			return
		}

		if details.Status == Deleted || details.Status == Pending {
			e <- errors.New(error_wrong_status)
			return
		}

		details.Status = Deleted
		details.UpdatedAt = time.Now()

		s.instances.Put(details)

		c <- details

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox deletion

		s.instances.Remove(details.UUID)

	}()

	select {
	case details := <-c:
		return details, nil
	case err := <-e:
		return SandboxDetails{}, err
	}
}

func (s *AzureSandbox) ListAll() []SandboxDetails {
	return s.instances.ListAll()
}

func (s *AzureSandbox) GetByName(name string) []SandboxDetails {
	return s.instances.GetByName(name)
}

func (s *AzureSandbox) GetByUUID(id string) (SandboxDetails, error) {
	return s.instances.GetByUUID(id)
}

func (s *AzureSandbox) UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error) {
	details, err := s.instances.GetByUUID(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	if details.Status == Deleted {
		return SandboxDetails{}, errors.New(error_already_deleted)
	}

	details.ExpiresAt = expiresAt
	details.UpdatedAt = time.Now()

	s.instances.Put(details)

	return details, nil
}
