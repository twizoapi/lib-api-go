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
	registration, _ := utils.AskForInput("Create biovoice registration for phone number [6100000000]: ", "6100000000")
	bioVoiceResponse, err := twizo.BioVoiceCreateRegistration(registration)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", bioVoiceResponse)
}
