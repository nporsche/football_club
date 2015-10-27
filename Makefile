GOPATH:=$(CURDIR)/.godeps:$(CURDIR)
export GOPATH

target: dep
	go build -o ./bin/football_club ./src/...

dep:
	-mkdir .godeps
	go get github.com/go-sql-driver/mysql
	go get github.com/nporsche/goyaml
