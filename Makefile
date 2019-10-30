GO ?= go

.PHONY: all
all: pinkie

%: cmd/%/main.go pkg/client/*.go
	$(GO) build -o $@ $<
