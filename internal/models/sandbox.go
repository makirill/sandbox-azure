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
	GetAll() ([]SandboxDetails, error)
	GetByName(name string) ([]SandboxDetails, error)
	GetByID(id string) (SandboxDetails, error)
	UpdateExpiration(id string, expiresAt time.Time) (bool, error)
	UpdateStatus(id string, status string) (bool, error)
}

type Sandbox interface {
	Create(name string, expireTime time.Time) (SandboxDetails, error)
	Remove(id string) (SandboxDetails, error)
	ListAll() ([]SandboxDetails, error)
	GetByName(name string) ([]SandboxDetails, error)
	GetByUUID(id string) (SandboxDetails, error)
	UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error)
}
