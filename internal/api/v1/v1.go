package v1

import (
	"net/http"
	"time"

	"github.com/makirill/sandbox-azure/internal/models"
)

type Error struct {
	Err            error `json:"-"`                   // low-level runtime error
	HTTPStatusCode int   `json:"http_code,omitempty"` // http response status code

	Message        string   `json:"message"`                   // user-facing
	AppCode        int64    `json:"code,omitempty"`            // application-specific error code
	ErrorText      string   `json:"error,omitempty"`           // application-level error message, for debugging
	ErrorMultiline []string `json:"error_multiline,omitempty"` // application-level error message, for debugging
}

type HealthCheckResult struct {
	HTTPStatusCode int    `json:"http_code,omitempty"` // http response status code
	Message        string `json:"message"`
}

type Sandbox struct {
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type SandboxList struct {
	SandboxDetailsList []models.SandboxDetails `json:"sandbox_details_list"`
	//	Sandboxes []models.Sandbox `json:"sandboxes"`
	Len int `json:"len"`
}

func (p *Error) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *HealthCheckResult) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Sandbox) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *SandboxList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
