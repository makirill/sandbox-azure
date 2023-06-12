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
	Status    SandboxStatus
}

type SandboxList interface {
	Put(sandbox SandboxDetails) error
	Remove(id string) (SandboxDetails, error)
	ListAll() []SandboxDetails
	GetByName(name string) []SandboxDetails
	GetByUUID(id string) (SandboxDetails, error)
}

type Sandbox interface {
	Create(name string) SandboxDetails
	Remove(id string) (SandboxDetails, error)
	ListAll() []SandboxDetails
	GetByName(name string) []SandboxDetails
	GetByUUID(id string) (SandboxDetails, error)
	UpdateExpiration(id string, expiresAt time.Time) (SandboxDetails, error)
}

type SandboxStatus int64

const (
	Running SandboxStatus = iota
	Stopped
	Expired
	Pending
	Failed
	Deleted
)

func (s SandboxStatus) String() string {
	return [...]string{"Running", "Stopped", "Expired", "Pending", "Failed", "Deleted"}[s]
}
