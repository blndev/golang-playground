package main

import (
	"fmt"
	"os"
	"strconv" //for parsebool

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/utils"
)

// Function to analyze a bool Flag value containes in *sring and return true if it is set
// Default Value for nil or not convertable is false
func isFlagged(tagValue *string) bool {
	if tagValue == nil {
		return false
	}
	b, err := strconv.ParseBool(*tagValue)
	if err == nil {
		return b
	}
	return false
}

func main() {

	//fmt.Printf("German Cloud: %+v\n", azure.GermanCloud)

	//TODO add ENV Flag to identify GErman or Global Azure CLoud
	authorizer, _ := utils.GetAuthorizer(azure.GermanCloud)
	vmClient := compute.NewVirtualMachinesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	vmClient.Authorizer = authorizer
	vmClient.BaseURI = azure.GermanCloud.ResourceManagerEndpoint

	//TODO find a way to receive VMs by using a query
	list, err := vmClient.ListAll()
	onErrorFail(err, "ListAll failed")
	if list.Value != nil && len(*list.Value) > 0 {
		fmt.Println("VMs in subscription")
		for _, vm := range *list.Value {
			//display VMs only if they contains specific tags
			if vm.Tags != nil {
				if val, ok := (*vm.Tags)["scaleGroup"]; ok && isFlagged(val) {
					printVM(vm)
				} else {
					if vm.Name != nil {
						fmt.Printf("%s is not  part of target VMs", *vm.Name)
					}
				}
			}
		}
	} else {
		fmt.Println("There are no VMs in this subscription")
	}
}

// printVM prints basic info about a Virtual Machine.
func printVM(vm compute.VirtualMachine) {
	//see https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachine for details
	tags := "\n"
	if vm.Tags == nil {
		tags += "\t\tNo tags yet\n"
	} else {
		for key, value := range *vm.Tags {
			tags += fmt.Sprintf("\t\t%s = %s\n", key, *value)
		}
	}
	fmt.Printf("Virtual machine '%s'\n", *vm.Name)

	//see https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachineProperties
	var vmProps = *vm.VirtualMachineProperties
	var test = (*vmProps.HardwareProfile).VMSize
	elements := map[string]interface{}{
		"ID":                *vm.ID,
		"Type":              *vm.Type,
		"Location":          *vm.Location,
		"Tags":              tags,
		"ProvisioningState": *vmProps.ProvisioningState,
		"Computername":      (*vmProps.OsProfile).ComputerName,
		"OS":                (*(*vmProps.StorageProfile).OsDisk).OsType,

		"VMSize": (*vmProps.HardwareProfile).VMSize, //compare with const values from https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachineSizeTypes
		"test":   test}
	for k, v := range elements {
		fmt.Printf("\t%s: %s\n\n\n", k, v)
	}
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}
