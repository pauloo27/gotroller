build: build-common build-cli build-gui

build-common:
	go build -v

build-cli:
	cd ./cli; go build -v -o ../gotroller

build-gui: 
	cd ./gui; go build -v -o ../gotroller-gui

test:
	go test -v
