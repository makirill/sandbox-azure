package models

import (
	"time"

	"github.com/google/uuid"
)

type AzureSandboxInstance struct {
	details SandboxDetails
}

type AzureSubscription struct {
	subscriptionID string
	instances      []AzureSandboxInstance
}

func InitAzureSubscription(subscriptionID string) *AzureSubscription {
	return &AzureSubscription{
		subscriptionID: subscriptionID,
		instances:      []AzureSandboxInstance{},
	}
}

func (s *AzureSubscription) Add(name string) (*SandboxDetails, error) {
	azureDetails := AzureSandboxInstance{
		details: SandboxDetails{
			Name:      name,
			UUID:      uuid.New().String(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		},
	}

	// TODO: Create the sandbox in Azure here
	// and add the Azure Resources details to the AzureSandboxIndex structure
	s.instances = append(s.instances, azureDetails)

	return &azureDetails.details, nil
}

func (s *AzureSubscription) Remove(uuid string) (*SandboxDetails, error) {
	var sandboxDetails *SandboxDetails

	for i, sandboxInstance := range s.instances {
		if sandboxInstance.details.UUID == uuid {
			sandboxDetails = &sandboxInstance.details

			// TODO: Delete the sandbox in Azure here
			s.instances = append(s.instances[:i], s.instances[i+1:]...)
			break
		}
	}

	return sandboxDetails, nil
}

func (s *AzureSubscription) ListAll() []SandboxDetails {
	sandboxDetailsList := []SandboxDetails{}

	for _, sandboxInstance := range s.instances {
		sandboxDetailsList = append(sandboxDetailsList, sandboxInstance.details)
	}

	return sandboxDetailsList
}

func (s *AzureSubscription) GetByName(name string) []SandboxDetails {

	sandboxDetailsList := []SandboxDetails{}

	for _, sandboxInstance := range s.instances {
		if sandboxInstance.details.Name == name {
			sandboxDetailsList = append(sandboxDetailsList, sandboxInstance.details)
		}
	}

	return sandboxDetailsList
}

func (s *AzureSubscription) GetByUUID(uuid string) *SandboxDetails {
	var sandboxDetails *SandboxDetails

	for _, sandboxInstance := range s.instances {
		if sandboxInstance.details.UUID == uuid {
			sandboxDetails = &sandboxInstance.details
			break
		}
	}

	return sandboxDetails
}

func (s *AzureSubscription) UpdateExpiration(uuid string, expiresAt time.Time) *SandboxDetails {
	var sandboxDetails *SandboxDetails

	for i := range s.instances {
		if s.instances[i].details.UUID == uuid {
			s.instances[i].details.ExpiresAt = expiresAt
			sandboxDetails = &s.instances[i].details
			break
		}
	}

	return sandboxDetails
}
