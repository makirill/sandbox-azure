package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/google/uuid"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (client *azureClient) RegisterApplication(displayName string) (string, error) {
	requestBody := graphmodels.NewApplication()
	requestBody.SetDisplayName(&displayName)
	requestBody.SetIdentifierUris([]string{"api://" + displayName})

	result, err := client.graphServiceClient.Applications().Post(client.ctx, requestBody, nil)
	if err != nil {
		return "", err
	}

	return *result.GetAppId(), nil
}

func (client *azureClient) CreateServicePrincipal(appId string) (string, error) {
	requestBody := graphmodels.NewServicePrincipal()
	requestBody.SetAppId(&appId)

	result, err := client.graphServiceClient.ServicePrincipals().Post(client.ctx, requestBody, nil)
	if err != nil {
		return "", err
	}

	return *result.GetId(), nil
}

func (client *azureClient) setRoleAssignments(spID string, roleDefinitionID string) (string, error) {

	// TODO: Is it correct to assign the role to the service principal?
	scope := "subscriptions/" + client.subscriptionID

	clientFactory, err := armauthorization.NewClientFactory(client.subscriptionID, client.cred, nil)
	if err != nil {
		return "", err
	}

	roleProps := armauthorization.RoleAssignmentProperties{
		PrincipalID:      &spID,
		RoleDefinitionID: &roleDefinitionID,
		PrincipalType:    to.Ptr(armauthorization.PrincipalTypeServicePrincipal),
	}

	result, err := clientFactory.NewRoleAssignmentsClient().Create(
		client.ctx,
		scope,
		uuid.New().String(),
		armauthorization.RoleAssignmentCreateParameters{
			Properties: &roleProps,
		}, nil)

	if err != nil {
		return "", err
	}

	return *result.ID, nil

}

// The filter to apply on the operation. Use atScopeAndBelow filter to search below the given scope as well.
// https://learn.microsoft.com/en-us/rest/api/authorization/role-definitions/list?tabs=HTTP
func (client *azureClient) GetRoleDefinitions(filter string) ([]string, error) {
	clientFactory, err := armauthorization.NewClientFactory(client.subscriptionID, client.cred, nil)
	if err != nil {
		return nil, err
	}

	roleDefinitions := []string{}

	pager := clientFactory.NewRoleDefinitionsClient().NewListPager(
		"subscriptions/"+client.subscriptionID,
		&armauthorization.RoleDefinitionsClientListOptions{Filter: &filter})
	for pager.More() {
		page, err := pager.NextPage(client.ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			roleDefinitions = append(roleDefinitions, *v.ID)
		}
	}
	return roleDefinitions, nil
}
