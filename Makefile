build: build-cli build-gui

build-cli:
	cd ./cli; go build -v -o ../gotroller

build-gui: 
	cd ./gui; go build -v -o ../gotroller-gui
