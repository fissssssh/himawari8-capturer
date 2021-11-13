FILENAME=himawari8_capturer
SHELL=/bin/bash
OUTPUTDIR:=build
VERSION:=0.0.0
.PHONY: clean build rebuild
rebuild: clean build

# clean execute files.
clean:
	@echo "Cleaning……"
	rm -rf $(OUTPUTDIR)

# buiil windows, linux and macos executed files.
build: cmd/cli/main.go
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(OUTPUTDIR)/$(FILENAME)_linux_amd64_$(VERSION) $<
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $(OUTPUTDIR)/$(FILENAME)_darwin_amd64_$(VERSION) $<
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o $(OUTPUTDIR)/$(FILENAME)_windows_amd64_$(VERSION).exe $<