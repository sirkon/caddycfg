test:
		PATH=${GOPATH}/bin:${PATH}
		go install github.com/sirkon/caddycfg
		go get -u github.com/stretchr/testify
		go test -test.v github.com/sirkon/ldetool/testing


