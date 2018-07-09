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

	err := twizo.TotpDelete(totpIdentifier)
	if err != nil {
		fmt.Printf("%#v", err)
		panic(err)
	}

	fmt.Printf("Totp for identifier [%s] deleted\n", totpIdentifier)
}
