package main

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/utils"
	"os"
	"regexp"  //to extract resourcegroup
	"strconv" //for parsebool
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

var vmClient compute.VirtualMachinesClient

func main() {

	//fmt.Printf("German Cloud: %+v\n", azure.GermanCloud)

	//TODO add ENV Flag to identify GErman or Global Azure CLoud
	authorizer, _ := utils.GetAuthorizer(azure.GermanCloud)
	vmClient = compute.NewVirtualMachinesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	vmClient.Authorizer = authorizer
	vmClient.BaseURI = azure.GermanCloud.ResourceManagerEndpoint

	//TODO find a way to receive VMs by using a query
	list, err := vmClient.ListAll()
	//TODO there is a "nextLink" in the result if the result is paged
	//see https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachineListResult
	//use ListAllNExtResults(previousResult)
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
	fmt.Println("\n---------------------------------------------------------------")
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

	//find resourcegroup
	r, _ := regexp.Compile("resourceGroups/(.*)/providers")
	//r.MatchString(*vm.ID)
	resoureGroup := r.FindStringSubmatch(*vm.ID)[1]
	fmt.Println(resoureGroup)

	//see https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachineProperties
	//statuses seems to be filled only by vmClient.Get(rg,vmName, compute.InstanceView)
	var vmProps = *vm.VirtualMachineProperties
	var statuses *[]compute.InstanceViewStatus
	if vmProps.InstanceView != nil {
		statuses = (*vmProps.InstanceView).Statuses
	}

	if statuses != nil {
		for _, status := range *statuses {
			fmt.Printf("Status %s, Message %s", *status.Code, *status.DisplayStatus)
		}
	} else {
		fmt.Println("no status found, get detailed vm data")
		//see https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachinesClient.Get
		detailedVM, err := vmClient.Get(resoureGroup, *vm.Name, compute.InstanceView)
		fmt.Println(*detailedVM.Name)
		//TODO refactor status extraction to separate function
		if (*detailedVM.VirtualMachineProperties).InstanceView != nil {
			statuses = (*(*detailedVM.VirtualMachineProperties).InstanceView).Statuses
		}

		if statuses != nil {
			for _, status := range *statuses {
				fmt.Printf("Status %s, Message %s", *status.Code, *status.DisplayStatus)
			}
		}
		if err != nil {
			//seems not to work, maybe recursion is not possible with go
			printVM(detailedVM)
		} else {
			fmt.Println(err)
		}
	}

	elements := map[string]interface{}{
		"ID":                *vm.ID,
		"Name":              *vm.Name,
		"Type":              *vm.Type,
		"Location":          *vm.Location,
		"Tags":              tags,
		"ProvisioningState": *vmProps.ProvisioningState,
		"Computername":      *(*vmProps.OsProfile).ComputerName,
		"OS-TYpe":           (*(*vmProps.StorageProfile).OsDisk).OsType,
		"VMSize":            (*vmProps.HardwareProfile).VMSize, //compare with const values from https://godoc.org/github.com/Azure/azure-sdk-for-go/arm/compute#VirtualMachineSizeTypes
		//"OS-DiskName":       *(*(*vmProps.StorageProfile).OsDisk).Name,
		//"Instance Status": (*(*vmProps.InstanceView).Statuses),
		"ResourceGroup": resoureGroup,
		"test":          statuses}
	for k, v := range elements {
		fmt.Printf("\t%s: %s\n", k, v)
	}
	fmt.Println("---------------------------------------------------------------")
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}
