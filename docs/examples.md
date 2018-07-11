# Introduction
This document will show some examples, checking error responses are left up to the
user, and are removed here for readablity.

# Authentication
```go
import (
	twizo "github.com/twizoapi/lib-api-go"
)
func main() {
	twizo.ApiKey = "<API KEY>"
	// the default region is 'eu'
	twizo.RegionCurrent = "asia"
	// rest of code
}
```
# Credit balance
```go
response, err := twizo.BalanceGet()

fmt.Printf("Current balance for wallet: %s\n", response.GetWallet())
fmt.Printf(
  "  Current credit: %0.2f %s\n",
  response.GetCredit(),
  response.GetCurrencyCode(),
)
fmt.Printf("  Free verifications left: %d\n", response.GetFreeVerifications())
```
# Verification
## Submit
```go
verificationResponse, err := twizo.VerificationSubmit("1234567890")
```
## Verify
```go
err = verificationResponse.Verify("12345")

if !verificationResponse.IsTokenSuccess() {
  fmt.Println("Token was not correct")
}
```
## Status
```go
verificationResponse, err := twizo.VerificationStatus("1234567890")
fmt.Printf("Verification status is [%s]\n", verificationResponse.GetStatusMsg())

```
## Verification types
```go
verificationTypes := twizo.NewVerificationTypes()
verificationTypes.Fetch()

fmt.Printf("%v", verificationTypes)
```
## Widget
### Create
```go
sessionRequest := twizo.NewWidgetSessionRequest()
sessionRequest.SetAllowedTypes([]string{"sms", "call"})
sessionRequest.SetRecipient(twizo.Recipient("1234567890"))
response, err := sessionRequest.Submit()
```

### Status
```go
response, err := sessionRequest.Submit()
```
### Verify
```go
err = response.Verify()
if response.IsTokenSuccess() {
  fmt.Println("Token was correct")
} else {
  fmt.Println("Token was not correct")
}
```

# Registration widget
```go
sessionRequest := twizo.NewRegistrationWidgetSessionRequest()
sessionRequest.SetAllowedTypes([]string{"line", "backupcode"})
sessionRequest.SetRecipient(twizo.Recipient("1234567890"))	sessionRequest.SetBackupCodeIdentifier("<IDENTIFIER>")
response, err := sessionRequest.Submit()
```

# Backup Codes
## Create
```go
backupCreateRequest, err := twizo.BackupCodeCreate("<IDENTIFIER>")
fmt.Printf(
  "Created [%d] backup tokens for [%s]:\n\t%s\n",
  len(backupCreateRequest.GetCodes()),
  backupCreateRequest.GetIdentifier(),
  strings.Join(backupCreateRequest.GetCodes(), "\n\t"),
)
```

## Verify
```go
backupVerifyRequest, err := twizo.BackupCodeVerify("<IDENTIFIER>", "<TOKEN>")
if backupVerifyRequest.GetVerificationResponse().IsTokenSuccess() {
  fmt.Println("Token was correct")
} else {
  fmt.Println("Token was not correct")
}
```

## Check remaining
```go
response, err := twizo.BackupCodeStatus("<IDENTIFIER>")
fmt.Printf(
  "[%s] as [%d] backup tokens left\n",
  response.GetIdentifier(),
  response.GetAmountOfCodesLeft(),
)

```

## Update
```go
backupCreateRequest, err := twizo.BackupCodeUpdate("<IDENTIFIER>")
fmt.Printf(
  "Updated [%d] backup tokens for [%s]:\n\t%s\n",
  len(backupCreateRequest.GetCodes()),
  backupCreateRequest.GetIdentifier(),
  strings.Join(backupCreateRequest.GetCodes(), "\n\t"),
)
```
## Delete
```go
err := twizo.BackupCodeDelete("<IDENTIFIER>")
```
# Time-based One-Time Password
## Create
```go
response, err := twizo.TotpCreate("<IDENTIFIER>", "<ISSUER>")
if aErr, ok := err.(*twizo.APIError); ok && aErr.Conflict() {
  fmt.Printf("Token was already created for this identifier.\n")
  return
}

// for Create the url/secret is supplied
fmt.Printf("Token was created:\n  URL: %s\n  Secret: %s\n", response.GetURL(), *response.GetURLSecret())
```
## Verify
```go
response, err := twizo.TotpVerify("<IDENTIFIER>", "<TOKEN>")
if aErr, ok := err.(*twizo.APIError); ok && aErr.NotFound() {
  fmt.Printf("Identifier was not found.\n")
  return
}

if !response.GetVerificationResponse().IsTokenSuccess() {
  fmt.Println("Token was not correct.")
}
```
## Status
```go
response, err := twizo.TotpCheck("<IDENTIFIER>")
// for Check the url/secret is NOT supplied (only for create)
fmt.Printf("Totp Issuer for identifier [%s] is [%s]\n", response.GetIssuer())
```

## Delete
```go
err := twizo.TotpDelete("<IDENTIFIER>")
```

# Biovoice
## Create
```go
bioVoiceResponse, err := twizo.BioVoiceCreateRegistration("1234567890")
```
## Status Registration
```go
bioVoiceResponse, err := twizo.BioVoiceCheckRegistration("1234567890")
```
## Status Subscription
```go
bioVoiceResponse, err := twizo.BioVoiceCheckSubscription("1234567890")
```
## Delete subscription
```go
err :=  twizo.BioVoiceDeleteSubscription("1234567890")
```
# SMS
## Submit simple
```go
smsResponses, err := twizo.SmsSubmit("1234567890", "Body of the sms", "Sender")
```
## Statuses (based on SmsSubmit)
```go
for _, smsResponse := range smsResponses.GetItems() {
  fmt.Printf(
    "Sms [%s] to recipient [%s] has status [%s]\n",
    smsResponse.GetMessageID(),
    smsResponse.GetRecipient(),
    smsResponse.GetStatusMsg(),
  )
}
```
## Status based on message-id
```go
smsResponse, err := twizo.SmsStatus("<MessageId>")
```

## Delivery reports
```go
smsPollResults, err := twizo.SmsPollStatus()
for _, smsResponse := range smsPollResults.GetItems() {
  fmt.Printf(
    "- Sms to [%s] has status [%s]\n",
    smsResponse.GetRecipient(),
    smsResponse.GetStatusMsg(),
  )
}
// delete the poll result we processed it above
err = smsPollResults.Delete()
```
# Numberlookup
## Simple
```go
numberlookupResponse, err := twizo.NumberLookupSubmit("1234567890")
```
## Status
```go
// print status of all messages
for _, nlR := range numberlookupResponse.GetItems() {
  fmt.Printf(
    "- Numberlookup result for [%s] has status [%s] operator [%v]\n",
    nlR.GetNumber(),
    nlR.GetStatusMsg(),
    nlR.GetOperator(),
  )
}
```
