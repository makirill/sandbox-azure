package models

import (
	"errors"
	"time"

	"github.com/makirill/sandbox-azure/internal/log"
)

type AzureSandbox struct {
	subscriptionID string
	instances      SandboxData
}

func InitAzureSandbox(subscriptionID string) *AzureSandbox {
	// TODO: it should not be there
	sb, _ := InitAzureSandboxesPostgres("postgres://postgres:abc123@localhost/sandbox")

	return &AzureSandbox{
		subscriptionID: subscriptionID,
		instances:      sb,
	}
}

func (s *AzureSandbox) Create(name string) (SandboxDetails, error) {

	c := make(chan string)
	e := make(chan error)

	go func() {
		id, err := s.instances.Insert(name)
		if err != nil {
			e <- err
			return
		}

		c <- id

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox creation

		_, err = s.instances.UpdateStatus(id, "Running")
		if err != nil {
			log.Logger.Error("Failed to update status for sandbox: ", id, err)
		}

	}()

	select {
	case id := <-c:
		return s.instances.GetByID(id)
	case err := <-e:
		return SandboxDetails{}, err
	}
}

func (s *AzureSandbox) Remove(id string) (SandboxDetails, error) {
	c := make(chan bool)
	e := make(chan error)

	go func() {
		details, err := s.instances.GetByID(id)
		if err != nil {
			e <- err
			return
		}

		if details.Status == "Deleted" || details.Status == "Pending" {
			e <- errors.New(error_wrong_status)
			return
		}

		ok, err := s.instances.UpdateStatus(id, "Deleted")
		if err != nil {
			e <- err
			return
		}

		c <- ok

		time.Sleep(30 * time.Second) //TODO: replace it with the actual Azure sandbox deletion

		s.instances.Delete(id)

	}()

	select {
	case <-c:
		return s.instances.GetByID(id)
	case err := <-e:
		return SandboxDetails{}, err
	}
}

func (s *AzureSandbox) ListAll() ([]SandboxDetails, error) {
	return s.instances.GetAll()
}

func (s *AzureSandbox) GetByName(name string) ([]SandboxDetails, error) {
	return s.instances.GetByName(name)
}

func (s *AzureSandbox) GetByUUID(id string) (SandboxDetails, error) {
	return s.instances.GetByID(id)
}

func (s *AzureSandbox) UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error) {
	details, err := s.instances.GetByID(id)
	if err != nil {
		return SandboxDetails{}, err
	}

	if details.Status == "Deleted" {
		return SandboxDetails{}, errors.New(error_already_deleted)
	}

	_, err = s.instances.UpdateExpiration(id, expiresAt)
	if err != nil {
		return SandboxDetails{}, err
	}

	return s.instances.GetByID(id)
}
