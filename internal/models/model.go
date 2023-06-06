package models

import "time"

type SandboxDetails struct {
	Name      string
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type Sandbox interface {
	Add(name string) (*SandboxDetails, error)
	Remove(uuid string) (*SandboxDetails, error)
	ListAll() []SandboxDetails
	GetByName(name string) []SandboxDetails
	GetByUUID(uuid string) *SandboxDetails
	UpdateExpiration(uuid string, expiresAt time.Time) *SandboxDetails
}
