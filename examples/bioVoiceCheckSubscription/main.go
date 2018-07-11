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
	//
	// Biovoice might not work with a test key
	//
	utils.CheckKey(true)

	registration, _ := utils.AskForInput("Check biovoice subscription for phone number [6100000000]: ", "6100000000")
	bioVoiceResponse, err := twizo.BioVoiceCheckSubscription(registration)
	if aErr, ok := err.(*twizo.APIError); ok && aErr.NotFound() {
		fmt.Printf("Biovoice subscription for [%s] not found.\n", registration)
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", bioVoiceResponse)
}
