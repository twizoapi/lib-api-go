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
	//
	// Biovoice might not work with a test key
	//
	utils.CheckKey(true)

	registration, _ := utils.AskForInput("Delete biovoice registration for phone number [6100000000]: ", "6100000000")
	err := twizo.BioVoiceDeleteSubscription(registration)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted biovoice registration for [%s]\n", registration)
}
