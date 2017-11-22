package main

import (
	"fmt"
	twizo "github.com/twizoapi/lib-api-go"
	"github.com/twizoapi/lib-api-go/examples"
	"time"
)

func main() {
	utils.Main()

	twizo.APIKey = utils.SuppliedApiKey
	twizo.RegionCurrent = twizo.APIRegion(utils.SuppliedRegion)

	//
	// Note: error handeling was abreviated for example's sake
	//

	phone, _ := utils.AskForInput("Enter phone number [6100000000]: ", "6100000000")

	numberLookupRequest := twizo.NewNumberLookupRequest([]twizo.Recipient{twizo.Recipient(phone)})
	numberLookupRequest.SetResultType(twizo.ResultTypePolling)

	_, err := numberLookupRequest.Submit()
	if err != nil {
		panic(err)
	}

	// now we can poll for the results, keep in mind all poll results will be retrieved not only
	// our sent message
	maxPolls := 10
	for i := 1; i <= maxPolls; i++ {
		fmt.Printf("Polling [%d/%d]\n", i, maxPolls)
		numberLookupResult, err := twizo.NumberLookupPollStatus()
		if err != nil {
			panic(err)
		}
		for _, numberLookupMessage := range numberLookupResult.GetItems() {
			fmt.Printf(
				"- Numberlookup for [%s] has operator [%s]\n",
				numberLookupMessage.GetNumber(),
				utils.AsString(numberLookupMessage.GetOperator(), "-unknown-"),
			)
		}
		// delete the poll result we processed it above
		err = numberLookupResult.Delete()
		if err != nil {
			// failed to delete
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}
