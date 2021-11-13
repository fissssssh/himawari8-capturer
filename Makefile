FILENAME=himawari8_capturer
SHELL=/bin/bash
OUTPUTDIR:=build
VERSION:=v0.0.0
.PHONY: clean build rebuild
rebuild: clean build

# clean execute files.
clean:
	@echo "Cleaning"
	rm -rf $(OUTPUTDIR)/$(VERSION) vendor

# buiil windows, linux and macos executed files.
build: cmd/cli/main.go
	@echo "Restore packages"
	go mod tidy && go mod vendor
	@echo "Build for linux"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(OUTPUTDIR)/$(VERSION)/$(FILENAME)_linux_amd64_$(VERSION) $<
	@echo "Build for macos"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $(OUTPUTDIR)/$(VERSION)/$(FILENAME)_macos_amd64_$(VERSION) $<
	@echo "Build for windows"
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o $(OUTPUTDIR)/$(VERSION)/$(FILENAME)_windows_amd64_$(VERSION).exe $<