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
	totpIssuer, _ := utils.AskForInput("Issuer [TwizoTest]: ", "TwizoTest")

	response, err := twizo.TotpCreate(totpIdentifier, totpIssuer)
	if aErr, ok := err.(*twizo.APIError); ok && aErr.Conflict() {
		fmt.Printf("Token was already created for this identifier.\n")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%#v", err)
		panic(err)
	}

	fmt.Printf("Token was created:\n  URL: %s\n  Secret: %s\n", response.GetURL(), *response.GetURLSecret())
}
