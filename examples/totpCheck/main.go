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

	totpIdentifier, _ := utils.AskForInput("Identifier [TwizoTest]: ", "TwizoTest")

	response, err := twizo.TotpCheck(totpIdentifier)
	if err != nil {
		fmt.Printf("%#v", err)
		panic(err)
	}

	fmt.Printf("Totp Issuer for identifier [%s] is [%s]\n", response.GetIssuer())
}
