VERSION = $(shell git describe --tags)

# GOFLAGS = -ldflags "-X main.Version=${VERSION} -s -w" -trimpath -mod=readonly
VERS = main.Version=${VERSION}
LDFLAGS = "-s -w -X '${VERS}'"
GOFLAGS = -ldflags ${LDFLAGS} -trimpath -buildvcs=true 
BIN = "hn"
GOBINARY = go

all: darwin-merge 
local: local-arch
darwin: darwin-merge

darwin-merge: darwin-arm darwin-amd
	@echo "Merging binaries"
	@lipo -create ${BIN}-arm ${BIN}-amd -output ${BIN}
	@rm -f ${BIN}-arm ${BIN}-amd 

darwin-arm: 
	@echo "Building darwin arm64"
	@env CGO_ENABLED=0 GOWORK=off GOARCH=arm64 GOOS=darwin ${GOBINARY} build -o ${BIN}-arm ${GOFLAGS}

darwin-amd: 
	@echo "Building darwin amd64"
	@env CGO_ENABLED=0 GOWORK=off GOARCH=amd64 GOOS=darwin ${GOBINARY} build -o ${BIN}-amd ${GOFLAGS}

local-arch: 
	@echo "Building for local arch (${LOCAL_OS} ${LOCAL_ARCH})"
	@env CGO_ENABLED=0 GOWORK=off ${GOBINARY} build -o ${BIN} ${GOFLAGS}

