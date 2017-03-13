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


	// can send to multiple recipients below
	numberlookupResponse, err := twizo.NumberLookupSubmit(phone)
	if err != nil {
		panic(err)
	}

	// print status of all messages
	for _, nlR := range numberlookupResponse.GetItems() {
		fmt.Printf(
			"- Numberlookup result for [%s] has status [%s] operator [%s]\n",
			nlR.GetNumber(),
			nlR.GetStatusMsg(),
			utils.AsString(nlR.GetOperator(), "-unknown-"),
		)
	}

	fmt.Println("Sleeping some time.")
	time.Sleep(10 * time.Second)
	// retrieve the status, of all messages sent above
	err = numberlookupResponse.Status()
	if err != nil {
		panic(err)
	}
	for _, nlR := range numberlookupResponse.GetItems() {
		fmt.Printf(
			"- Numberlookup result for [%s] has status [%s] operator [%s]\n",
			nlR.GetNumber(),
			nlR.GetStatusMsg(),
			utils.AsString(nlR.GetOperator(), "-unknown-"),
		)
	}
}
