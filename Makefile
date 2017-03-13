SOURCEDIR           = $(shell pwd)

# set the GOPATH on execution only if not
# already set in environment
ifeq ($(strip $(GOPATH)),)
    SET_GOPATH = env GOPATH="$(SOURCEDIR)/../../../../"
endif

PACKAGE             = twizo-lib-go
BIN                 = $(GOPATH)/bin
BASE                = $(GOPATH)/src/$(PACKAGE)
GO_LINT             = $(BIN)/gometalinter.v1
GO_RICHGO           = $(BIN)/richgo
GO_HTTPMOCK         = $(BASE)/gopkg.in/jarcoal/httpmock.v1


# these targets do not actually produce
# output files
.PHONY: \
    all \
    richtest \
    test \
	lint

all: test

$(BASE):
	@mkdir -p $(dir $@)

$(GO_LINT): | $(BASE)
	go get gopkg.in/alecthomas/gometalinter.v1

$(GO_RICHGO): | $(BASE)
	go get github.com/kyoh86/richgo

$(GO_HTTPMOCK): | $(BASE)
	go get gopkg.in/jarcoal/httpmock.v1

lint: | $(GO_LINT)
	$(GO_LINT) --disable-all --enable=golint

richtest: | $(GO_RICHGO) $(GO_HTTPMOCK)
	$(SET_GOPATH) $(GO_RICHGO) test -v -coverprofile=c.out

test: | $(GO_HTTPMOCK)
	$(SET_GOPATH) go test -v -coverprofile=c.out

