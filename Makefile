OUTPUT_DIR = build

clean:
	rm -rf $(OUTPUT_DIR)

test:
	go test

# Since this plugin uses sockets and Docker only supports sockets for linux/bsd
# we make sure that CGO is disabled and regardless of host system, the binary is
# crosscompiled for linux
# see https://golang.org/cmd/cgo/
build: test
	CGO_ENABLED=0 GOOS=linux go build -o $(OUTPUT_DIR)/sesame
