PACKAGE      = github.com/twizoapi/lib-api-go
GOPATH       = $(CURDIR)/.gopath~
BASE         = $(GOPATH)/src/$(PACKAGE)
BIN          = $(GOPATH)/bin

GO_SRC          := $(wildcard $(CURDIR)/*.go)
GO_EXAMPLES_SRC := $(wildcard $(CURDIR)/examples/*/*.go)
GO_EXAMPLES_BIN := $(addprefix bin/,$(notdir $(GO_EXAMPLES_SRC:.go=)))

VERBOSE = 0
Q = $(if $(filter 1,${VERBOSE}),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

# Internal functions to find X in A by position and return Y from B (on same position)
# mapping GO_EXAMPLES_BIN -> GO_EXAMPLES_SRC
_pos = $(if $(findstring $1,$2),$(call _pos,$1,\
       $(wordlist 2,$(words $2),$2),x $3),$3)
pos = $(words $(call _pos,$1,$2))
lookup = $(word $(call pos,$1,$2),$3)

# If we don't set this some stacks may not be complete when encountering race
# conditions. Uses a bit more memory, but we usually have enough of that.
export GORACE="history_size=4"
export GOPATH

# these targets do not actually produce
# output files
PHONY := ''

default: test

# Create the "fake GOPATH settings"
$(BASE): ; $(info $(M) setting GOPATH…)
	$(Q) mkdir -p $(dir $@)
	$(Q) ln -sf $(CURDIR) $@

# Set up go lint
GO_LINT = $(BIN)/golint
$(GO_LINT): | $(BASE) ; $(info $(M) building golint…)
	$(Q) go get github.com/golang/lint/golint

PHONY += lint
lint: | $(GO_LINT) ; $(info $(M) running golint…) @ ## Run golint
	$(Q) cd $(BASE) && ret=0 && for pkg in $(GO_SRC); do \
		test -z "$$($(GO_LINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

# Set up go fmt
PHONY += fmt
fmt:
	$(Q) gofmt -s -w $(GO_SRC)
	$(Q) gofmt -s -w $(GO_EXAMPLES_SRC)

# Set up go metalinter
GO_METALINTER       = $(BIN)/gometalinter.v1
$(GO_METALINTER): | $(BASE) ; $(info $(M) building gometalinter…)
	$(Q) go get gopkg.in/alecthomas/gometalinter.v1; $(GO_METALINTER) --install

PHONY += metalinter
metalinter: | $(GO_METALINTER) ; $(info $(M) running metalinter…) @ ## Run golint
	$(Q) $(GO_METALINTER) \
		--deadline=120s \
		--vendor \
		--tests \
		--disable-all \
		--enable=vet \
		--enable=golint \
        --enable=ineffassign \
        --enable=errcheck \
        --enable=lll \
        --enable=deadcode \
        --line-length=120 \
        --vendored-linters

# Richgo
# GO_RICHGO           = $(GOPATH)/bin/richgo
#richtest: | $(GO_RICHGO) $(GO_HTTPMOCK)
#	$(GO_RICHGO) test -v -coverprofile=c.out

GO_HTTPMOCK = $(GOPATH)/src/gopkg.in/jarcoal/httpmock.v1
$(GO_HTTPMOCK): | $(BASE) ; $(info $(M) installing httpmock.v1…)
	$(Q) go get gopkg.in/jarcoal/httpmock.v1

PHONY += test_units
test_units: | $(GO_HTTPMOCK)
	$(Q) go test -v -coverprofile=c.out

PHONY += travis
travis: test_units examples

PHONY += test
test: test_units

PHONY += examples
examples: $(GO_EXAMPLES_BIN)

$(GO_EXAMPLES_BIN): $(GO_EXAMPLES_SRC) | $(BASE) ; $(info $(M) building $@…)
	$(Q) go build -o $@ $(call lookup,$@,$(GO_EXAMPLES_BIN),$(GO_EXAMPLES_SRC))

# Declare the contents of the .PHONY variable as phony.  We keep that
# information in a variable so we can use it in if_changed and friends.
.PHONY: $(PHONY)