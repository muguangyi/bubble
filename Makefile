export GO111MODULE=on

test:
	./codecov.sh

build:
	cd bubble-master && go build && cd -
	cd bubble-worker && go build && cd -
