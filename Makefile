BINARY=tibberexporter

COMMIT_HASH=`git rev-parse --short=8 HEAD 2>/dev/null`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-s -w \
    -X github.com/metskem/tibberexporter/conf.CommitHash=${COMMIT_HASH} \
    -X github.com/metskem/tibberexporter/conf.BuildTime=${BUILD_TIME}"

clean:
	go clean
	if [ -f ./target/linux_amd64/${BINARY} ] ; then rm ./target/${BINARY}-linux_amd64 ; fi
	if [ -f ./target/darwin_amd64/${BINARY} ] ; then rm ./target/${BINARY}-darwin_amd64 ; fi
	if [ -f ./target/darwin_arm64/${BINARY} ] ; then rm ./target/${BINARY}-darwin_arm64 ; fi
	if [ -f ./target/linux_arm64/${BINARY} ] ; then rm ./target/${BINARY}-linux_arm64 ; fi

release: clean linux_amd64 darwin_amd64 linux_arm64 darwin_arm64

deps:
	go get -v ./...

linux_amd64: deps
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./target/linux_amd64/${BINARY} ${LDFLAGS} .

darwin_amd64: deps
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ./target/darwin_amd64/${BINARY} ${LDFLAGS} .

linux_arm64: deps
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./target/linux_arm64/${BINARY} ${LDFLAGS} .

darwin_arm64: deps
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./target/darwin_arm64/${BINARY} ${LDFLAGS} .
