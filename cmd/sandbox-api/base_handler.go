package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	v1 "github.com/makirill/sandbox-azure/internal/api/v1"
	"github.com/makirill/sandbox-azure/internal/log"
	"github.com/makirill/sandbox-azure/internal/models"
)

type BaseHandler struct {
	sandboxes *models.SandboxList
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		sandboxes: models.NewSandboxList(),
	}
}

func (h *BaseHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {

	log.Logger.Info("Health check", "status", "OK")
	render.Render(w, r, &v1.HealthCheckResult{
		HTTPStatusCode: 200,
		Message:        "OK",
	})

}

func (h *BaseHandler) GetSandboxHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "sandboxName")
	log.Logger.Info("Get sandbox", "name", name)

	if sandbox := h.sandboxes.Get(name); sandbox != nil {
		log.Logger.Info("Sandbox found", "name", name)
		render.Status(r, 200)
		render.Render(w, r, &v1.Sandbox{
			Name:      sandbox.Name,
			CreatedAt: &sandbox.CreatedAt,
			UpdatedAt: &sandbox.UpdatedAt,
		})
	} else {
		log.Logger.Info("Sandbox not found", "name", name)
		render.Status(r, 404)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 404,
			Message:        "Sandbox not found",
		})
	}

}

func (h *BaseHandler) ListSandboxesHandler(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("List all sandboxes")

	if sandboxes := h.sandboxes.List(); len(sandboxes) > 0 {
		log.Logger.Info("Sandboxes found", "count", len(sandboxes))
		render.Status(r, 200)
		render.Render(w, r, &v1.SandboxList{
			Sandboxes: sandboxes,
			Len:       len(sandboxes),
		})
	} else {
		log.Logger.Info("Sandboxes not found")
		render.Status(r, 404)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 404,
			Message:        "Sandboxes not found",
		})
	}
}

func (h *BaseHandler) CreateSandboxHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "sandboxName")
	subscriptionID := r.URL.Query().Get("subscriptionId")
	if sandbox := h.sandboxes.Get(name); sandbox != nil {
		log.Logger.Info("Sandbox already exists", "name", name)
		render.Status(r, 409)
		render.Render(w, r, &v1.Sandbox{
			Name:      sandbox.Name,
			CreatedAt: &sandbox.CreatedAt,
			UpdatedAt: &sandbox.UpdatedAt,
			Message:   "Sandbox already exists",
		})
	} else {
		log.Logger.Info("Create sandbox", "name", name)
		sandbox, err := models.NewSandbox(name, subscriptionID)
		if err != nil {
			log.Logger.Error("Error creating sandbox", "name", name, "error", err.Error())
			render.Status(r, 500)
			render.Render(w, r, &v1.Error{
				HTTPStatusCode: 500,
				Message:        err.Error(),
			})
		} else {
			h.sandboxes.Add(sandbox)
			render.Status(r, 201)
			render.Render(w, r, &v1.Sandbox{
				Name: name,
			})
		}
	}
}

func (h *BaseHandler) DeleteSandboxHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "sandboxName")
	if sandbox := h.sandboxes.Remove(name); sandbox != nil {
		log.Logger.Info("Sandbox deleted", "name", name)
		render.Status(r, 204)
		render.Render(w, r, &v1.Sandbox{
			Name:      sandbox.Name,
			CreatedAt: &sandbox.CreatedAt,
			UpdatedAt: &sandbox.UpdatedAt,
			Message:   "Sandbox deleted",
		})
	} else {
		log.Logger.Info("Sandbox not found", "name", name)
		render.Status(r, 404)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 404,
			Message:        "Sandbox not found",
		})
	}
}
