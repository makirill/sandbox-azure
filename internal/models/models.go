package models

import (
	"time"

	"github.com/makirill/sandbox-azure/internal/azure"
)

type Sandbox struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	resources *azure.AzureResources
}

type SandboxList struct {
	Sandboxes []Sandbox
}

func NewSandbox(name string, subscriptionID string) (*Sandbox, error) {
	azureResources, err := azure.CreateSandbox(name, subscriptionID)
	if err != nil {
		return nil, err
	}

	sandbox := Sandbox{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		resources: azureResources,
	}
	return &sandbox, nil
}

func NewSandboxList() *SandboxList {
	return &SandboxList{
		Sandboxes: []Sandbox{},
	}
}

func (l *SandboxList) Add(s *Sandbox) {
	l.Sandboxes = append(l.Sandboxes, *s)
}

func (l *SandboxList) Remove(name string) *Sandbox {
	for i, s := range l.Sandboxes {
		if s.Name == name {
			l.Sandboxes = append(l.Sandboxes[:i], l.Sandboxes[i+1:]...)
			return &s
		}
	}
	return nil
}

func (l *SandboxList) List() []Sandbox {
	return l.Sandboxes
}

func (l *SandboxList) Get(name string) *Sandbox {
	for _, s := range l.Sandboxes {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (l *SandboxList) Update(name string) {
	for i, s := range l.Sandboxes {
		if s.Name == name {
			l.Sandboxes[i].UpdatedAt = time.Now()
		}
	}
}

func (l *SandboxList) Len() int {
	return len(l.Sandboxes)
}

func (l *SandboxList) isExist(name string) bool {
	for _, s := range l.Sandboxes {
		if s.Name == name {
			return true
		}
	}
	return false
}
