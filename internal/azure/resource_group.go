package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func (client *azureClient) createResourceGroup(resourceGroupName string, location string) (*armresources.ResourceGroup, error) {
	resourceGroupResp, err := client.resourceGroupClient.CreateOrUpdate(
		client.ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(location),
		},
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func (client *azureClient) getResourceGroup(resourceGroupName string) (*armresources.ResourceGroup, error) {

	resourceGroupResp, err := client.resourceGroupClient.Get(
		client.ctx,
		resourceGroupName,
		nil)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

func (client *azureClient) listResourceGroups() ([]*armresources.ResourceGroup, error) {

	resultPager := client.resourceGroupClient.NewListPager(nil)
	resourceGroups := make([]*armresources.ResourceGroup, 0)
	for resultPager.More() {
		pageResp, err := resultPager.NextPage(client.ctx)
		if err != nil {
			return nil, err
		}
		resourceGroups = append(resourceGroups, pageResp.ResourceGroupListResult.Value...)
	}

	return resourceGroups, nil
}

func (client *azureClient) checkExistenceResourceGroup(resourceGroupName string) (bool, error) {

	boolResp, err := client.resourceGroupClient.CheckExistence(
		client.ctx,
		resourceGroupName,
		nil)
	if err != nil {
		return false, err
	}
	return boolResp.Success, nil

}
