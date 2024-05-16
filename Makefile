WORKDIR ?= .work
BUILDS  ?= ${WORKDIR}/builds

.PHONY: all
all: bin

.PHONY: bin
bin: ${BUILDS}/wpsermon
${BUILDS}/wpsermon: *.go pkg/*.go
	mkdir -p ${BUILDS}
	go build -o $@
