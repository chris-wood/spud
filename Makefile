INCLUDES+= -I=.
INCLUDES+= -I=$(GOPATH)/src/
INCLUDES+= -I=/usr/local/include

check:
	go test ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -type f)

configure:
	./PREREQ.sh
