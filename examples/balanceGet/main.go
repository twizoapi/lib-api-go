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
	response, err := twizo.BalanceGet()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current balance for wallet: %s\n", response.GetWallet())
	fmt.Printf(
		"  Current credit: %0.2f %s\n",
		response.GetCredit(),
		response.GetCurrencyCode(),
	)
	fmt.Printf("  Free verifications left: %d\n", response.GetFreeVerifications())

}
