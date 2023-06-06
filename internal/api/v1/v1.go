package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/makirill/sandbox-azure/internal/models"
)

type Error struct {
	Err error `json:"-"` // low-level runtime error

	Message        string   `json:"message"`                   // user-facing
	AppCode        int64    `json:"code,omitempty"`            // application-specific error code
	ErrorText      string   `json:"error,omitempty"`           // application-level error message, for debugging
	ErrorMultiline []string `json:"error_multiline,omitempty"` // application-level error message, for debugging
}

type HealthCheckResult struct {
	Message string `json:"message"`
}

type Sandbox struct {
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type SandboxList struct {
	SandboxList []models.SandboxDetails `json:"sandbox_list"`
}

type SandboxRequest struct {
	Name string `json:"name"`
}

type SandboxUpdateRequest struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *HealthCheckResult) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Sandbox) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (sl *SandboxList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (sr *SandboxRequest) Bind(r *http.Request) error {
	if sr.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

func (sur *SandboxUpdateRequest) Bind(r *http.Request) error {
	if sur.ExpiresAt.IsZero() {
		return errors.New("missing expires_at")
	}

	return nil
}
