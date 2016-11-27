all: local test alpine linux

test:
	docker run --rm -ti -v $(CURDIR):/usr/local/src/github.com/qnib/qcollect/ --workdir /usr/local/src/github.com/qnib/qcollect qnib/golang ./test.sh

local:
	./build.sh

alpine:
	docker run --rm -ti -v $(CURDIR):/usr/local/src/github.com/qnib/qcollect/ --workdir /usr/local/src/github.com/qnib/qcollect qnib/alpn-go-dev ./build.sh

linux:
	docker run --rm -ti -v $(CURDIR):/usr/local/src/github.com/qnib/qcollect --workdir /usr/local/src/github.com/qnib/qcollect qnib/golang ./build.sh
