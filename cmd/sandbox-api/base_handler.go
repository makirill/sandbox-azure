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
	model models.Sandbox
}

// TDOD: add subscriptionID here
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		model: models.InitAzureSubscription("test"),
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

	if sandboxDetailsList := h.model.GetByName(name); len(sandboxDetailsList) == 0 {
		log.Logger.Info("Sandbox not found", "name", name)
		render.Status(r, 404)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 404,
			Message:        "Sandbox not found",
		})
	} else {
		log.Logger.Info("Sandbox found", "name", name)
		render.Status(r, 200)
		render.Render(w, r, &v1.SandboxList{
			SandboxDetailsList: sandboxDetailsList,
			Len:                len(sandboxDetailsList),
		})
	}

}

func (h *BaseHandler) ListSandboxesHandler(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("List all sandboxes")

	if sandboxDetailsList := h.model.ListAll(); len(sandboxDetailsList) == 0 {
		log.Logger.Info("Sandboxes not found")
		render.Status(r, 404)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 404,
			Message:        "Sandbox not found",
		})
	} else {
		log.Logger.Info("Sandboxes found", "count", len(sandboxDetailsList))
		render.Status(r, 200)
		render.Render(w, r, &v1.SandboxList{
			SandboxDetailsList: sandboxDetailsList,
			Len:                len(sandboxDetailsList),
		})
	}

}

func (h *BaseHandler) CreateSandboxHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "sandboxName")

	if sandboxDetailsList := h.model.GetByName(name); len(sandboxDetailsList) > 0 {
		log.Logger.Info("Sandbox already exists", "name", name)
		render.Status(r, 409)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 409,
			Message:        "Sandbox already exists",
		})
	} else {
		log.Logger.Info("Create sandbox", "name", name)
		sandboxDetails, err := h.model.Add(name)
		if err != nil {
			log.Logger.Error("Error creating sandbox", "name", name, "error", err.Error())
			render.Status(r, 500)
			render.Render(w, r, &v1.Error{
				HTTPStatusCode: 500,
				Message:        "Error creating sandbox",
			})
		} else {
			log.Logger.Info("Sandbox created", "name", name)
			render.Status(r, 200)
			render.Render(w, r, &v1.Sandbox{
				Name:      sandboxDetails.Name,
				CreatedAt: &sandboxDetails.CreatedAt,
				UpdatedAt: &sandboxDetails.UpdatedAt,
				ExpiresAt: &sandboxDetails.ExpiresAt,
			})
		}
	}

}

func (h *BaseHandler) DeleteSandboxHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "sandboxUUID")
	sandboxDetails, err := h.model.Remove(uuid)
	if err != nil {
		log.Logger.Error("Error deleting sandbox", "uuid", uuid, "error", err.Error())
		render.Status(r, 500)
		render.Render(w, r, &v1.Error{
			HTTPStatusCode: 500,
			Message:        "Error deleting sandbox",
		})
	} else {
		if sandboxDetails == nil {
			log.Logger.Info("Sandbox not found", "uuid", uuid)
			render.Status(r, 404)
			render.Render(w, r, &v1.Error{
				HTTPStatusCode: 404,
				Message:        "Sandbox not found",
			})
		} else {
			log.Logger.Info("Sandbox deleted", "uuid", uuid)
			render.Status(r, 200)
			render.Render(w, r, &v1.Sandbox{
				Name:      sandboxDetails.Name,
				CreatedAt: &sandboxDetails.CreatedAt,
				UpdatedAt: &sandboxDetails.UpdatedAt,
				ExpiresAt: &sandboxDetails.ExpiresAt,
			})
		}
	}
}
