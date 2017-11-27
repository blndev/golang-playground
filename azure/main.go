package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/utils"
)

func main() {

	//fmt.Printf("German Cloud: %+v\n", azure.GermanCloud)

	authorizer, _ := utils.GetAuthorizer(azure.GermanCloud)
	vmClient := compute.NewVirtualMachinesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	vmClient.Authorizer = authorizer
	vmClient.BaseURI = azure.GermanCloud.ResourceManagerEndpoint

	list, err := vmClient.ListAll()
	onErrorFail(err, "ListAll failed")
	if list.Value != nil && len(*list.Value) > 0 {
		fmt.Println("VMs in subscription")
		for _, vm := range *list.Value {
			printVM(vm)
		}
	} else {
		fmt.Println("There are no VMs in this subscription")
	}
}

// printVM prints basic info about a Virtual Machine.
func printVM(vm compute.VirtualMachine) {
	tags := "\n"
	if vm.Tags == nil {
		tags += "\t\tNo tags yet\n"
	} else {
		for k, v := range *vm.Tags {
			tags += fmt.Sprintf("\t\t%s = %s\n", k, *v)
		}
	}
	fmt.Printf("Virtual machine '%s'\n", *vm.Name)
	elements := map[string]interface{}{
		"ID":       *vm.ID,
		"Type":     *vm.Type,
		"Location": *vm.Location,
		"Tags":     tags}
	for k, v := range elements {
		fmt.Printf("\t%s: %s\n", k, v)
	}
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}
