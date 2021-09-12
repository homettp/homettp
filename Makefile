VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean ui darwin linux windows post_clean

pre_clean:
	rm -rf dist

ui:
	yarn prod
	mkdir dist
	cp .env.example dist/.env

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/homettp
	cd dist && zip -r homettp_$(VERSION)_darwin_amd64.zip .env homettp
	rm -f dist/homettp

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/homettp
	cd dist && zip -r homettp_$(VERSION)_linux_amd64.zip .env homettp
	rm -f dist/homettp

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/homettp.exe
	cd dist && zip -r homettp_$(VERSION)_windows_amd64.zip .env homettp.exe
	rm -f dist/homettp.exe

post_clean:
	rm -rf dist/.env
