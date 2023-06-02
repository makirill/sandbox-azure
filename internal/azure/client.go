package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type azureClient struct {
	subscriptionID      string
	ctx                 context.Context
	cred                *azidentity.DefaultAzureCredential
	graphServiceClient  *msgraphsdk.GraphServiceClient
	resourceGroupClient *armresources.ResourceGroupsClient
}

type AzureResources struct {
	SubscriptionID string
	ResourceGroup  string
	ApplicationID  string
}

func newAzureClient(subscriptionID string) (*azureClient, error) {

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	graphClient, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, nil)
	if err != nil {
		return nil, err
	}

	resourceClientFactory, err := armresources.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		log.Fatal(err)
	}
	resourceGroupClient := resourceClientFactory.NewResourceGroupsClient()

	return &azureClient{
		subscriptionID:      subscriptionID,
		graphServiceClient:  graphClient,
		cred:                cred,
		ctx:                 context.Background(),
		resourceGroupClient: resourceGroupClient}, nil
}

var (
	location = "eastus"
)

func CreateSandbox(name string, subscriptionID string) (resources *AzureResources, err error) {

	azureClient, err := newAzureClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	exist, err := azureClient.checkExistenceResourceGroup(name)
	if err != nil {
		log.Fatal(err)
	}

	var resourceGroup *armresources.ResourceGroup
	if !exist {
		resourceGroup, err = azureClient.createResourceGroup(name, location)
		if err != nil {
			return nil, err
		}
	}

	appId, err := azureClient.RegisterApplication("test-sample-application")
	if err != nil {
		return nil, err
	}

	spID, err := azureClient.CreateServicePrincipal(appId)
	if err != nil {
		return nil, err
	}

	roleDefinitions, err := azureClient.GetRoleDefinitions("roleName eq 'Owner'")
	if err != nil {
		return
	}
	for _, role := range roleDefinitions {
		_, err := azureClient.setRoleAssignments(spID, role)
		if err != nil {
			return nil, err
		}
	}

	//	return nil, errors.New("not implemented")
	return &AzureResources{
		SubscriptionID: subscriptionID,
		ResourceGroup:  *resourceGroup.ID,
		ApplicationID:  appId,
	}, nil
}
