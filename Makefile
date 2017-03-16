SOURCEDIR           = $(shell pwd)

ifndef GOPATH
    export GOPATH=$(SOURCEDIR)/../../../../
endif

GO_LINT             = $(GOPATH)/bin/golint
GO_METALINTER       = $(GOPATH)/bin/gometalinter.v1
GO_RICHGO           = $(GOPATH)/bin/richgo
GO_HTTPMOCK         = $(GOPATH)/src/gopkg.in/jarcoal/httpmock.v1

EXAMPLES            = $(wildcard $(SOURCEDIR)/examples/*/*.go)

# these targets do not actually produce
# output files
.PHONY: \
    all \
    travis \
    richtest \
    examples \
    test_units \
	lint

all: test

travis: test_units examples
test: test_units

$(GO_LINT): |
	go get github.com/golang/lint/golint

$(GO_METALINTER): |
	go get gopkg.in/alecthomas/gometalinter.v1

$(GO_RICHGO): |
	go get github.com/kyoh86/richgo

$(GO_HTTPMOCK): |
	@echo $(GO_HTTPMOCK)
	go get gopkg.in/jarcoal/httpmock.v1

lint: | $(GO_LINT) $(GO_METALINTER)
	$(GO_METALINTER) --disable-all --enable=golint

richtest: | $(GO_RICHGO) $(GO_HTTPMOCK)
	$(GO_RICHGO) test -v -coverprofile=c.out

test_units: | $(GO_HTTPMOCK)
	go test -v -coverprofile=c.out

examples:
	go build -o bin/numberlookupPoll      examples/numberlookup/numberlookupPoll.go
	go build -o bin/numberlookupSimple    examples/numberlookup/numberlookupSimple.go
	go build -o bin/smsSimple             examples/sms/smsSimple.go
	go build -o bin/smsPoll               examples/sms/smsPoll.go
	go build -o bin/verificationAdvanced  examples/verification/verificationAdvanced.go
	go build -o bin/verificationSimpleSms examples/verification/verificationSimpleSms.go
