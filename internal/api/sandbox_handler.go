package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/makirill/sandbox-azure/internal/log"
	"github.com/makirill/sandbox-azure/internal/models"
)

type SandboxHandler struct {
	instances models.SandboxController
}

// Helper to map the string status to the SandboxStatus enum
func toSandboxStatus(status string) SandboxStatus {
	var statusMap = map[string]SandboxStatus{
		"deleted": DELETED,
		"expired": EXPIRED,
		"failed":  FAILED,
		"pending": PENDING,
		"running": RUNNING,
		"stopped": STOPPED,
	}

	ret, ok := statusMap[strings.ToLower(status)]
	if !ok {
		ret = UNKNOWN
	}

	return ret
}

// Make sure we conform to the StrictServerInterface
var _ StrictServerInterface = (*SandboxHandler)(nil)

func String(s string) *string {
	return &s
}

func NewSandboxHandler(controller models.SandboxController) *SandboxHandler {

	return &SandboxHandler{
		instances: controller,
	}
}

func (sh *SandboxHandler) Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error) {
	status := OK

	// TODO: check DB connection

	return Health200JSONResponse{
		Status: &status,
	}, nil
}

func (sh *SandboxHandler) ListSandboxes(ctx context.Context, request ListSandboxesRequestObject) (ListSandboxesResponseObject, error) {

	detailsList, err := sh.instances.ListAll(request.Params.Limit, request.Params.Offset)
	if err != nil {
		return ListSandboxesdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	sandboxes := make([]Sandbox, 0, len(detailsList))
	for _, details := range detailsList {
		sandboxes = append(sandboxes, Sandbox{
			Name:      details.Name,
			Id:        details.UUID,
			Status:    toSandboxStatus(details.Status),
			CreatedAt: details.CreatedAt,
			ExpiresAt: details.ExpiresAt,
			UpdatedAt: details.UpdatedAt,
		})
	}

	return ListSandboxes200JSONResponse(sandboxes), nil
}

func (sh *SandboxHandler) CreateSandbox(ctx context.Context, request CreateSandboxRequestObject) (CreateSandboxResponseObject, error) {
	sandboxDetails, err := sh.instances.Create(request.Body.Name, request.Body.ExpiresAt)
	if err != nil {
		return CreateSandboxdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	log.Logger.Info("Sandbox created", "name", sandboxDetails.Name, "id", sandboxDetails.UUID)

	return CreateSandbox201JSONResponse{
		Body: Sandbox{
			Name:      sandboxDetails.Name,
			Id:        sandboxDetails.UUID,
			Status:    toSandboxStatus(sandboxDetails.Status),
			CreatedAt: sandboxDetails.CreatedAt,
			ExpiresAt: sandboxDetails.ExpiresAt,
			UpdatedAt: sandboxDetails.UpdatedAt,
		},
		Headers: CreateSandbox201ResponseHeaders{
			Location: "/sandboxes/" + sandboxDetails.UUID,
		},
	}, nil

}

func (sh *SandboxHandler) GetSandboxByName(ctx context.Context, request GetSandboxByNameRequestObject) (GetSandboxByNameResponseObject, error) {
	detailsList, err := sh.instances.GetByName(request.Name)
	if err != nil {
		return GetSandboxByNamedefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	sandboxes := make([]Sandbox, 0, len(detailsList))
	for _, details := range detailsList {
		sandboxes = append(sandboxes, Sandbox{
			Name:      details.Name,
			Id:        details.UUID,
			Status:    toSandboxStatus(details.Status),
			CreatedAt: details.CreatedAt,
			ExpiresAt: details.ExpiresAt,
			UpdatedAt: details.UpdatedAt,
		})
	}

	return GetSandboxByName200JSONResponse(sandboxes), nil
}

func (sh *SandboxHandler) DeleteSandbox(ctx context.Context, request DeleteSandboxRequestObject) (DeleteSandboxResponseObject, error) {
	sandboxDetails, err := sh.instances.Remove(request.Id)
	if err != nil {
		return DeleteSandboxdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	log.Logger.Info("Sandbox deleted", "name", sandboxDetails.Name, "id", sandboxDetails.UUID)

	return DeleteSandbox204Response{}, nil
}

func (sh *SandboxHandler) GetSandbox(ctx context.Context, request GetSandboxRequestObject) (GetSandboxResponseObject, error) {
	sandboxDetails, err := sh.instances.GetByUUID(request.Id)
	if err != nil {
		return GetSandboxdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	return GetSandbox200JSONResponse(Sandbox{
		Name:      sandboxDetails.Name,
		Id:        sandboxDetails.UUID,
		Status:    toSandboxStatus(sandboxDetails.Status),
		CreatedAt: sandboxDetails.CreatedAt,
		ExpiresAt: sandboxDetails.ExpiresAt,
		UpdatedAt: sandboxDetails.UpdatedAt,
	}), nil
}

func (sh *SandboxHandler) UpdateSandbox(ctx context.Context, request UpdateSandboxRequestObject) (UpdateSandboxResponseObject, error) {

	sandboxDetails, err := sh.instances.UpdateExpiration(request.Id, request.Body.ExpiresAt)
	if err != nil {
		return UpdateSandboxdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}}, nil
	}

	return UpdateSandbox200JSONResponse(Sandbox{
		Name:      sandboxDetails.Name,
		Id:        sandboxDetails.UUID,
		Status:    toSandboxStatus(sandboxDetails.Status),
		CreatedAt: sandboxDetails.CreatedAt,
		ExpiresAt: sandboxDetails.ExpiresAt,
		UpdatedAt: sandboxDetails.UpdatedAt,
	}), nil
}
