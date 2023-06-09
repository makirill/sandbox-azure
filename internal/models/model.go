package models

import "time"

type SandboxDetails struct {
	Name      string
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	Status    SandboxStatus
}

type Sandbox interface {
	Add(name string) SandboxDetails
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
