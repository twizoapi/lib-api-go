package main

import (
	"fmt"
	"strings"

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
	identifier, _ := utils.AskForInput("Update backup codes for  []: ", "")
	backupCreateRequest, err := twizo.BackupCodeUpdate(identifier)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Updated [%d] backup tokens for [%s]:\n\t%s\n",
		len(backupCreateRequest.GetCodes()),
		backupCreateRequest.GetIdentifier(),
		strings.Join(backupCreateRequest.GetCodes(), "\n\t"),
	)
}
