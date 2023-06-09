package models

import "time"

const (
	error_not_found       = "not found"
	error_wrong_status    = "wrong status"
	error_already_deleted = "already deleted"
)

type SandboxDetails struct {
	Name      string
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	Status    string
}

type SandboxData interface {
	Insert(name string, expireTime time.Time) (string, error)
	Delete(id string) (bool, error)
	GetAll(limit int, offset int) ([]SandboxDetails, error)
	GetByName(name string) ([]SandboxDetails, error)
	GetByID(id string) (SandboxDetails, error)
	UpdateExpiration(id string, expiresAt time.Time) (bool, error)
	UpdateStatus(id string, status string) (bool, error)
}

type SandboxController interface { //TODO: find a better name
	Create(name string, expireTime time.Time) (SandboxDetails, error)
	Remove(id string) (SandboxDetails, error)
	ListAll(limit int, offset int) ([]SandboxDetails, error)
	GetByName(name string) ([]SandboxDetails, error)
	GetByUUID(id string) (SandboxDetails, error)
	UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error)
}
