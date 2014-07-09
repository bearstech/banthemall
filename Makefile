default: build

src/github.com/nranchev/go-libGeoIP:
	GOPATH=`pwd` go get github.com/nranchev/go-libGeoIP

src/github.com/stretchr/testify/assert:
	GOPATH=`pwd` go get github.com/stretchr/testify/assert

lib: src/github.com/nranchev/go-libGeoIP

lib-test: src/github.com/stretchr/testify/assert

test: lib lib-test
	GOPATH=`pwd` go test -cover banthemall banthemall/combined banthemall/metrics

build: lib
	GOPATH=`pwd` go build banthemall
	GOPATH=`pwd` go build banthemall/static

linux: lib
	GOPATH=`pwd` gxc build-linux banthemall
