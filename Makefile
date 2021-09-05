GOLANGCI_LINT    ?= $(CURDIR)/bin/golangci-lint
GOLANGCI_VERSION ?= v1.41.1
$(GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(PWD)/bin" $(GOLANGCI_VERSION)

lint: $(GOLANGCI_LINT)
	"$(GOLANGCI_LINT)" run -v