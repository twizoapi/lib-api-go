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
	identifier, _ := utils.AskForInput("Verify backup codes for  []: ", "")
	token, _ := utils.AskForInput("Verification token  []: ", "")
	backupVerifyRequest, err := twizo.BackupCodeVerify(identifier, token)
	if err != nil {
		panic(err)
	}

	if backupVerifyRequest.GetVerificationResponse().IsTokenSuccess() {
		fmt.Println("Token was correct")
	} else {
		fmt.Println("Token was not correct")
	}
}
