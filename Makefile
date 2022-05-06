WORKDIR ?= .work
BUILDS  ?= ${WORKDIR}/builds
PKGDIR  ?= ${WORKDIR}/pkg

.PHONY: all
all: bin pkg macapp dmg

.PHONY: bin
bin: ${BUILDS}/wpsermon
${BUILDS}/wpsermon: *.go pkg/*.go
	mkdir -p ${BUILDS}
	go build -o $@

.PHONY: pkg
pkg: ${BUILDS}/wpsermon
	rm -rf ${PKGDIR}
	mkdir -p ${PKGDIR}
	cp -Rv ${BUILDS}/wpsermon ${PKGDIR}
	cp -Rv assets ${PKGDIR}

.PHONY: macapp
macapp: ${WORKDIR}/WPSermon.app
${WORKDIR}/WPSermon.app: ${BUILDS}/wpsermon ${PKGDIR}/*
	rm -rf $@
	go run ~/code/macapp/main.go                 \
		--assets ${PKGDIR}                       \
		--bin wpsermon                           \
		--icon assets/WPC_logo_brown_stacked.png \
		--identifier "church.whisperingpines"    \
		--name "WPSermon" -o ${WORKDIR}

.PHONY: dmg
dmg: ${WORKDIR}/WPSermon.dmg
${WORKDIR}/WPSermon.dmg: ${WORKDIR}/WPSermon.app
	rm -f $@
	create-dmg $@ ${WORKDIR}/WPSermon.app
