package main

import (
	"fmt"
	"strings"

	twizo "github.com/twizoapi/lib-api-go"
	"github.com/twizoapi/lib-api-go/examples"
)

func trimStrArr(s []string) []string {
	var r []string
	for _, str := range s {
		str = strings.Trim(str, " ")
		if len(str) > 0 {
			r = append(r, str)
		}
	}
	return r
}

func main() {
	utils.Main()

	twizo.APIKey = utils.SuppliedApiKey
	twizo.RegionCurrent = twizo.APIRegion(utils.SuppliedRegion)

	//
	// Note: error handling was abbreviated for example's sake
	//
	sessionRequest := twizo.NewWidgetSessionRequest()
	allowedTypesStr, _ := utils.AskForInput("AllowedType(s) (Delimiter: ', ') []: ", "")
	sessionRequest.SetAllowedTypes(trimStrArr(strings.Split(allowedTypesStr, ",")))

	if sessionRequest.GetAllowedTypes().Has(twizo.VerificationTypeSms, twizo.VerificationTypeCall, twizo.VerificationTypeBioVoice) {
		recipient, _ := utils.AskForInput("Number [6100000000]: ", "6100000000")
		sessionRequest.SetRecipient(twizo.Recipient(recipient))
	}

	if sessionRequest.GetAllowedTypes().Has(twizo.VerificationTypeBackupCode) {
		backupCodeIdentifier, _ := utils.AskForInput("BackupCodeIdentifier: ", "")
		sessionRequest.SetBackupCodeIdentifier(backupCodeIdentifier)
	}

	if sessionRequest.GetAllowedTypes().Has(twizo.VerificationTypeTotp) {
		totpIdentifier, _ := utils.AskForInput("TotpIdentifier: ", "")
		sessionRequest.SetTotpIdentifier(totpIdentifier)
	}

	if sessionRequest.GetAllowedTypes().Has(twizo.VerificationTypeTelegram, twizo.VerificationTypeLine, twizo.VerificationTypePush) {
		issuer, _ := utils.AskForInput("Issuer: ", "")
		sessionRequest.SetIssuer(issuer)
	}

	response, err := sessionRequest.Submit()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Please visit [https://widget.twizo.com/ver?sessionToken=%s&origin=] to complete the login procedure\n", response.GetSessionToken())
	_, err = utils.AskForInput("Press enter to continue: \n", "")

	err = response.Verify()
	if response.IsTokenSuccess() {
		fmt.Println("Token was correct")
	} else {
		fmt.Println("Token was not correct")
	}
}
