package main;

import (
	twizo "github.com/twizoapi/lib-api-go"
	"github.com/twizoapi/lib-api-go/examples"
	"fmt"
	"time"
)

func main() {
	utils.Main();

	twizo.APIKey = utils.SuppliedApiKey
	twizo.RegionCurrent = twizo.APIRegion(utils.SuppliedRegion)

	//
	// Note: error handeling was abreviated for example's sake
	//

	phone, _ := utils.AskForInput("Enter phone number [6100000000]: ", "6100000000")
	body, _  := utils.AskForInput("Enter body [Sample Sms]: ", "Sample Sms")

	// can send to multiple recipients below
	smsResponses, err := twizo.SmsSubmit(phone, body, "TwizoDemo")
	if err != nil {
		panic(err)
	}

	// print status of all messages
	for _, smsResponse := range smsResponses.GetItems() {
		fmt.Printf(
			"Sms [%s] to recipient [%s] has status [%s]\n",
			smsResponse.GetMessageID(),
			smsResponse.GetRecipient(),
			smsResponse.GetStatusMsg(),
		)
	}

	fmt.Println("Sleeping some time.")
	time.Sleep(5 * time.Second)
	// retrieve the status, of all messages sent above
	err = smsResponses.Status()
	if err != nil {
		panic(err)
	}
	for _, smsResponse := range smsResponses.GetItems() {
		fmt.Printf(
			"Sms [%s] to recipient [%s] has status [%s]\n",
			smsResponse.GetMessageID(),
			smsResponse.GetRecipient(),
			smsResponse.GetStatusMsg(),
		)
	}
}
