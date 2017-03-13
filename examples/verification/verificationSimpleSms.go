package main;

import (
	twizo "github.com/twizoapi/lib-api-go"
	"github.com/twizoapi/lib-api-go/examples"
	"fmt"
)

func main() {
	utils.Main();

	twizo.APIKey = utils.SuppliedApiKey
	twizo.RegionCurrent = twizo.APIRegion(utils.SuppliedRegion)

	//
	// Note: error handeling was abreviated for example's sake
	//

	// if not using a Test Application key it might fail, if using a test key
	// you can send to 6100000000 and validate it
	phone, _ := utils.AskForInput("Enter phone number [6100000000]: ", "6100000000")
	//verificationResponse, err := twizo.VerificationSubmit(twizo.Recipient(phone))
	verificationResponse, err := twizo.VerificationSubmit(phone)
	if err != nil {
		panic(err)
	}

	// valid Token for a Test Application with number 6100000000 is 012345, ask user
	token, _ := utils.AskForInput("Enter token [012345]: ", "012345")
	err = verificationResponse.Verify(token)
	if err != nil {
		panic(err)
	}

	if verificationResponse.IsTokenSuccess() {
		fmt.Println("Token was correct")
	} else {
		fmt.Println("Token was not correct")
	}
}
