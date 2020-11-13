default: test

test:
	go test ./...

fuzz:
	go run testdata/fuzz/main.go -n 10000

README.md: README.md.tpl $(wildcard *.go)
	becca -package $(subst $(GOPATH)/src/,,$(PWD))
