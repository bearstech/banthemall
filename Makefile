default: build

src/github.com/nranchev/go-libGeoIP:
	GOPATH=`pwd` go get github.com/nranchev/go-libGeoIP

lib: src/github.com/nranchev/go-libGeoIP

build: lib
	GOPATH=`pwd` go build banthemall

linux: lib
	GOPATH=`pwd` gxc build-linux banthemall
