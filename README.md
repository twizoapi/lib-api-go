<p align="center"><a href="https://www.twizo.com/" target="_blank">
    <img src="https://www.twizo.com/wp-content/themes/twizo/_/images/twizo-logo-0474ce6f.png" alt="Twizo">
</a></p>


# Twizo Go API #

Connect to the Twizo API using Go. This API includes functions to send verifications (2FA), SMS and Number Lookup.

## Requirements ##
* Go >= 1.7

## Get application secret and choose api region ##
To use the Twizo API client, the following things are required:

* Create a [Twizo account][twizo-account]
* Login on the Twizo portal
* Find your [application][twizo-application] secret
* Find your nearest api [region][twizo-doc-intro]

## Versioning
Each revision of the library is tagged and the version is updated accordingly.

Given Go's lack of built-in versioning, it is highly recommended you use a
[package management tool][package-management] in order to ensure a newer
version of the binding does not affect backwards compatibility.

## Installation
```sh
$ go get github.com/twizoapi/lib-api-go
```

## Getting started ##
Require the library using import, here we use twizo as an alias 

```go
import (
	twizo "github.com/twizoapi/lib-api-go"
)
```

Set the api key that was retrieved from the application section of the portal.

```go
func main() {
	twizo.ApiKey = "43reFDSrewrfet425rtefdGDSGds54twegdsgHaFST2refwd"
	// the default region is 'eu'
	twizo.RegionCurrent = "asia"
	// rest of code
}
```

### Verification ###
Create and send new verification

```go
verificationResponse, err := twizo.VerificationSubmit("610123456789")
if err != nil {
    // handle error
}
```

Verify token, based on the response above

```go
err = verificationResponse.Verify("12345")
if err != nil {
    // handle error
}

if verificationResponse.IsTokenSuccess() {
     // Token was correct
} else {
    // Token was not correct, expired, or other error
}
```

For more examples please see [Verification Examples][examples-verification]


### Sms ###
Submit the request using the simple method.  Using this method the Twizo Api will autodetect 
some of the settings for us.  The api will split up the message if the body is too long for
one SMS as well as some other settings.

```go
smsResponses, err := twizo.SmsSubmit("610123456789", "Greetings from Twizo", "TwizoDemo")
if err != nil {
    // handle error
}
```

Optional retrieve status of all sent sms using the response above, an smsResponses is a collection of smsResponse

```go
err = smsResponses.Status()
if err != nil {
    // handle error
}
for _, smsResponse := range smsResponses.GetItems() {
    fmt.Printf(
        "Sms [%s] to recipient [%s] has status [%s]\n",
        smsResponse.GetMessageID(),
        smsResponse.GetRecipient(),
        smsResponse.GetStatusMsg(),
    )
}
```

Optional retrieve status of sent sms using only the messageId, please note that this function
can only be used to retieve the status of one sms, using the response above will allow checking 
all messages in that request.

```go
smsResponse, err := twizo.SmsStatus("<MessageId>")
if err != nil {
        if (err.(*twizo.APIError).NotFound()) {
        	// Not found, it might have expired, or was never sent
        } else {
        	// handle other error
        }
} else {
        // smsResponse object was updated
}
```

For more examples please see [Sms Examples][examples-sms]

### Numberlookup ###
Create and submit a new numberlookup

```go
numberlookupResponse, err := twizo.NumberLookupSubmit("610123456789")
if err != nil {
    // handle error
}
```

Retrieve the result of the numberlookup

```go
err = numberlookupResponse.Status()
if err != nil {
    // handle error
}
for _, nlR := range numberlookupResponse.GetItems() {
    fmt.Printf(
        "- Numberlookup result for [%s] has status [%s] operator [%s]\n",
        nlR.GetNumber(),
        nlR.GetStatusMsg(),
        utils.AsString(nlR.GetOperator(), "-unknown-"),
    )
}
```

For more examples please see [Numberlookup Examples][examples-numberlookup]

## Examples ##
In the examples directory you can find a collection of examples of how to use the api. All examples can be
run using the following commands.

```sh
go run examples/verification/verificationSimple.go --key <key> --region <region>
go run examples/verification/verificationAdvanced.go --key <key> --region <region>
go run examples/sms/smsSimple.go --key <key> --region <region>
go run examples/sms/smsPoll.go --key <key> --region <region>
go run examples/numberlookupSimple.go --key <key> --region <region>
go run examples/numberlookupPoll.go --key <key> --region <region>
```


## Development
Pull requests from the community are welcome. If you submit one, please keep
the following guidelines in mind:

1. Code must be `go fmt` compliant.
2. All types, structs and funcs should be documented.
3. Ensure that `make test` succeeds.

Code can best be checked out like this (as an example, replace where needed)

```sh
export GOPATH="<DIR>"
git clone git@github.com:twizoapi/lib-api-go.git $(GOPATH)/src/github.com/twizoapi/lib-api-go
cd $(GOPATH)/src/github.com/twizoapi/lib-api-go
```

This way will allow tests to work

## ToDo ##
- [x] Add example for number lookup
- [x] Split up examples
- [ ] Add more tests for sms / verification / numberlookup
- [ ] Document all types / funcs / structs (point 2)

## License ##
[The MIT License][mit].

Copyright (c) 2016-2017 Twizo

## Support ##
Contact: [www.twizo.com][twizo] — support@twizo.com

[twizo]: http://www.twizo.com/
[mit]: https://opensource.org/licenses/mit-license.php
[twizo-account]: https://register.twizo.com/
[twizo-application]: https://portal.twizo.com/applications/
[twizo-doc-intro]: https://www.twizo.com/developers/documentation/#introduction_api-url
[package-management]: https://code.google.com/p/go-wiki/wiki/PackageManagementTools
[examples-verification]: examples/verification
[examples-sms]: examples/sms
[examples-numberlookup]: examples/numberlookup
