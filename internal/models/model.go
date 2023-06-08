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
	Add(name string) (*SandboxDetails, error)
	Remove(uuid string) (*SandboxDetails, error)
	ListAll() []SandboxDetails
	GetByName(name string) []SandboxDetails
	GetByUUID(uuid string) *SandboxDetails
	UpdateExpiration(uuid string, expiresAt time.Time) *SandboxDetails
}

type SandboxStatus int64

const (
	Running SandboxStatus = iota
	Stopped
	Expired
	Pending
	Failed
)

func (s SandboxStatus) String() string {
	return [...]string{"Running", "Stopped", "Expired", "Pending", "Failed"}[s]
}
