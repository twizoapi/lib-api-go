package main

import (
	"fmt"
	"os"

	twizo "github.com/twizoapi/lib-api-go"
	"github.com/twizoapi/lib-api-go/examples"
)

func main() {
	utils.Main()

	twizo.APIKey = utils.SuppliedApiKey
	twizo.RegionCurrent = twizo.APIRegion(utils.SuppliedRegion)

	//
	// Note: error handling was abbreviated for example's sake
	//
	identifier, _ := utils.AskForInput("Create backup codes for  []: ", "")
	response, err := twizo.BackupCodeStatus(identifier)
	if aErr, ok := err.(*twizo.APIError); ok && aErr.NotFound() {
		fmt.Printf("Identifier [%s] was not found.\n", identifier)
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"[%s] as [%d] backup tokens left\n",
		response.GetIdentifier(),
		response.GetAmountOfCodesLeft(),
	)
}
