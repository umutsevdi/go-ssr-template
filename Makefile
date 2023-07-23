SRC := $(wildcard app/*.go)
BIN := webwatch
BUILDDIR := .

.PHONY: all clean test

all: $(BUILDDIR)/$(BIN)

test:
	cd app/; go test ./...

$(BUILDDIR)/$(BIN):$(SRC)
	cd app/; go build
	mv app/app $(BUILDDIR)/$(BIN)

clean:
	rm -f $(BUILDDIR)/$(BIN)
