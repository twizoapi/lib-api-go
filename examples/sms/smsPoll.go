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
	body, _ := utils.AskForInput("Enter body [Sample Sms]: ", "Sample Sms")
	sender, _ := utils.AskForInput("Enter sender [TwizoDemo]: ", "TwizoDemo")

	// can send to multiple recipients below
	smsRequest, err := twizo.NewSmsRequest([]twizo.Recipient{twizo.Recipient(phone)}, body, sender)
	if err != nil {
		panic(err)
	}

	smsRequest.SetResultType(twizo.ResultTypePolling)

	_, err = smsRequest.Submit()
	if err != nil {
		panic(err)
	}

	// now we can poll for the results, keep in mind all poll results will be retrieved not only
	// our sent message
	maxPolls := 10
	for i := 1; i <= maxPolls; i++ {
		fmt.Printf("Polling [%d/%d]\n", i, maxPolls)
		smsPollResults, err := twizo.SmsPollStatus()
		if err != nil {
			panic(err)
		}

		//fmt.Print(utils.GetJsonFor(smsPollResults))
		for _, smsResponse := range smsPollResults.GetItems() {
			fmt.Printf(
				"- Sms to [%s] has status [%s]\n",
				smsResponse.GetRecipient(),
				smsResponse.GetStatusMsg(),
			)
		}
		// delete the poll result we processed it above
		err = smsPollResults.Delete()
		if err != nil {
			// failed to delete
			panic(err)
		}

		time.Sleep(2 * time.Second)
	}
}
