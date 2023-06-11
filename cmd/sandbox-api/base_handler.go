package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v1 "github.com/makirill/sandbox-azure/internal/api/v1"
	"github.com/makirill/sandbox-azure/internal/log"
	"github.com/makirill/sandbox-azure/internal/models"
)

type BaseHandler struct {
	model models.Sandbox
}

func NewBaseHandler(subscriptionID string) *BaseHandler {
	return &BaseHandler{
		model: models.InitAzureSandbox(subscriptionID),
	}
}

func (h *BaseHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &v1.HealthCheckResult{
		Message: "OK",
	})

}

func (h *BaseHandler) GetSandboxHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	sandbox, err := h.model.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Get Sandbox by ID: %s, error: %s", uuid, err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &v1.Sandbox{
		UUID:      sandbox.UUID,
		Name:      sandbox.Name,
		CreatedAt: sandbox.CreatedAt,
		UpdatedAt: sandbox.UpdatedAt,
		ExpiresAt: sandbox.ExpiresAt,
		Status:    sandbox.Status.String(),
	})
}

func (h *BaseHandler) GetSandboxByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	sandboxList := h.model.GetByName(name)
	if len(sandboxList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Sandbox %s not found", name),
		})
		return

	}

	var responsePayload v1.SandboxList
	responsePayload.SandboxList = make([]v1.Sandbox, len(sandboxList))
	for i, sandbox := range sandboxList {
		responsePayload.SandboxList[i] = v1.Sandbox{
			UUID:      sandbox.UUID,
			Name:      sandbox.Name,
			CreatedAt: sandbox.CreatedAt,
			UpdatedAt: sandbox.UpdatedAt,
			ExpiresAt: sandbox.ExpiresAt,
			Status:    sandbox.Status.String(),
		}
	}

	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &responsePayload)
}

func (h *BaseHandler) ListSandboxesHandler(w http.ResponseWriter, r *http.Request) {
	sandboxList := h.model.ListAll()
	if len(sandboxList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, &v1.Error{
			Message: "No Sandboxes found",
		})
		return
	}

	var responsePayload v1.SandboxList
	responsePayload.SandboxList = make([]v1.Sandbox, len(sandboxList))
	for i, sandbox := range sandboxList {
		responsePayload.SandboxList[i] = v1.Sandbox{
			UUID:      sandbox.UUID,
			Name:      sandbox.Name,
			CreatedAt: sandbox.CreatedAt,
			UpdatedAt: sandbox.UpdatedAt,
			ExpiresAt: sandbox.ExpiresAt,
			Status:    sandbox.Status.String(),
		}
	}

	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &responsePayload)

}

func (h *BaseHandler) CreateSandboxHandler(w http.ResponseWriter, r *http.Request) {
	data := &v1.SandboxRequest{}
	if err := render.Bind(r, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Failed to parse request body: %s", err),
		})
		return
	}

	sandbox := h.model.Create(data.Name) // TODO: add expiration date

	w.Header().Set("Location", fmt.Sprintf("/api/v1/sandboxes/%s", sandbox.UUID))
	w.WriteHeader(http.StatusCreated)
	render.Render(w, r, &v1.Sandbox{
		UUID:      sandbox.UUID,
		Name:      sandbox.Name,
		CreatedAt: sandbox.CreatedAt,
		UpdatedAt: sandbox.UpdatedAt,
		ExpiresAt: sandbox.ExpiresAt,
		Status:    sandbox.Status.String(),
	})

	log.Logger.Info("Sandbox created", "uuid", sandbox.UUID, "name", sandbox.Name)
}

func (h *BaseHandler) DeleteSandboxHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	sandboxDetails, err := h.model.Remove(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Logger.Error("Failed to delete sandbox", "uuid", uuid, "error", err)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Error delete sandbox %s: %s", uuid, err),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Logger.Info("Sandbox deleted", "uuid", sandboxDetails.UUID, "name", sandboxDetails.Name)
}

func (h *BaseHandler) UpdateSandboxHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	data := &v1.SandboxUpdateRequest{}

	if err := render.Bind(r, data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Failed to parse request body: %s", err),
		})
		return
	}

	sandboxDetails, err := h.model.UpdateExpiration(uuid, data.ExpiresAt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, &v1.Error{
			Message: fmt.Sprintf("Update Sandbox %s : %s", uuid, err.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &v1.Sandbox{
		UUID:      sandboxDetails.UUID,
		Name:      sandboxDetails.Name,
		CreatedAt: sandboxDetails.CreatedAt,
		UpdatedAt: sandboxDetails.UpdatedAt,
		ExpiresAt: sandboxDetails.ExpiresAt,
		Status:    sandboxDetails.Status.String(),
	})

	log.Logger.Info("Sandbox updated", "uuid", sandboxDetails.UUID, "name", sandboxDetails.Name)
}
