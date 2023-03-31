DESTDIR ?= /
PREFIX ?= ${DESTDIR}usr/local

all: zendebug zenexplorer restroom-test

zenexplorer_sources := main.go zencodeStatements.go delegate.go zenroom.go
zendebug_sources := main.go
restroom_test_sources := main.go

zenexplorer: $(wildcard src/zenexplorer/*.go)
	go build -o zenexplorer $(foreach source, ${zenexplorer_sources}, src/zenexplorer/${source})

zendebug: $(wildcard src/zendebug/*.go)
	go build -o zendebug $(foreach source, ${zendebug_sources}, src/zendebug/${source})

restroom-test: $(wildcard src/restroom-test/*.go)
	go build -o restroom-test $(foreach source, ${restroom_test_sources}, src/restroom-test/${source})

install:
	install zenexplorer zendebug restroom-test ${PREFIX}/bin
	
clean:
	rm -f zendebug zenexplorer restroom-test
