package main

import (
	"fmt"

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
	identifier, _ := utils.AskForInput("Delete backup codes for  []: ", "")
	err := twizo.BackupCodeDelete(identifier)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Backup codes for [%s] have been deleted\n", identifier)
}
