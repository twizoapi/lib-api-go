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
	response, err := twizo.ApplicationVerifyCredentials()
	if err != nil {
		panic(err)
	}
	if !response.IsKeyValid() {
		fmt.Printf("Api Key is not valid\n")
	} else {
		fmt.Printf("Api Key is valid:\n")
		fmt.Printf("\tTest Key: %v\n", response.IsTestKey())
		fmt.Printf("\tApplication Tag: %v\n", response.GetApplicationTag())
	}
}
