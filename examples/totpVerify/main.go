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

	totpIdentifier, _ := utils.AskForInput("Identifier [TwizoTest]: ", "TwizoTest")
	totpToken, _ := utils.AskForInput("Token []: ", "")

	response, err := twizo.TotpVerify(totpIdentifier, totpToken)
	if aErr, ok := err.(*twizo.APIValidationError); ok && aErr.UnprocessableEntity() {
		fmt.Printf("%s\n", aErr.Detail())
		os.Exit(1)
	}
	if aErr, ok := err.(*twizo.APIError); ok && aErr.NotFound() {
		fmt.Printf("Identifier was not found.\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%#v", err)
		panic(err)
	}

	if response.GetVerificationResponse().IsTokenSuccess() {
		fmt.Println("Token was correct.")
	} else {
		fmt.Println("Token was not correct.")
	}
}
