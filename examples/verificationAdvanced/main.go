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

	// if not using a Test Application key it might fail, if using a test key
	// you can send to 6100000000 and validate it
	phone, _ := utils.AskForInput("Enter phone number [6100000000]: ", "6100000000")
	verificationRequest, err := twizo.NewVerificationRequest(phone)
	if err != nil {
		panic(err)
	}
	verificationRequest.SetBodyTemplate("Your verification code is %token%")
	verificationRequest.SetTokenLength(10)
	verificationRequest.SetTokenType(twizo.VerificationTokenTypeAlpha)
	verificationResponse, err := verificationRequest.Submit()
	if err != nil {
		panic(err)
	}
	// now we assume that verificationResponse.MessageId was saved, and
	// retrieved in messageId
	messageID := verificationResponse.GetMessageID()

	// valid Token for a Test Application (alphanumeric) with number
	// 6100000000 is a123456789, ask user
	token, _ := utils.AskForInput("Enter token [a123456789]: ", "a123456789")
	verificationResponse, err = twizo.VerificationVerify(messageID, token)
	if err != nil {
		panic(err)
	}

	// check the status of the response
	if verificationResponse.IsTokenSuccess() {
		fmt.Println("- Submitted token was correct")
	} else {
		fmt.Println("- Submitted token was not correct")
	}

	// as a test we revalidated with the same token
	verificationResponse, err = twizo.VerificationVerify(messageID, token)
	if verificationResponse.IsTokenSuccess() {
		// this will never happen as the token was already validated
		fmt.Println("- Submiting token second time was correct")
	} else if verificationResponse.IsTokenAlreadyVerified() {
		fmt.Println("-- Token already verified, status is unknown")
	} else {
		fmt.Println("-- Submitted token was not correct")
	}

	// retieve the status, if some other application did not Validate
	verificationResponse, err = twizo.VerificationStatus(messageID)
	if err != nil {
		panic(err)
	}
	if verificationResponse.IsTokenSuccess() {
		fmt.Println("- Status of intial submitted token was correct")
	} else {
		fmt.Println("- Status of initial submitted token was incorrect")
	}
}
