DESTDIR ?= /
PREFIX ?= ${DESTDIR}usr/local

all: zendebug zenexplorer restroom-test breakroom

zenexplorer_sources := main.go zencodeStatements.go delegate.go zenroom.go
zendebug_sources := main.go
restroom_test_sources := main.go

zenexplorer: $(wildcard src/zenexplorer/*.go)
	go build -o zenexplorer $(foreach source, ${zenexplorer_sources}, src/zenexplorer/${source})

zendebug: $(wildcard src/zendebug/*.go)
	go build -o zendebug $(foreach source, ${zendebug_sources}, src/zendebug/${source})

restroom-test: $(wildcard src/restroom-test/*.go)
	go build -o restroom-test $(foreach source, ${restroom_test_sources}, src/restroom-test/${source})

.PHONY: breakroom
breakroom: $(wildcard src/breakroom/*.c)
	${CC} -Os -o src/breakroom/breakroom-read \
	 src/breakroom/main.c src/breakroom/bestline.c
	chmod +x src/breakroom/breakroom
	cp src/breakroom/breakroom-read src/breakroom/breakroom .

## in case we want to include the binary in the shell script
# gzip -1 -c src/breakroom/breakroom-read | \
# 	base64 -w 0 - > src/breakroom/breakroom-read.gz.b64
# ls -lh src/breakroom/breakroom-read.gz.b64
# cat src/breakroom/breakroom-read.gz.b64

install:
	-install zenexplorer zendebug restroom-test ${PREFIX}/bin
	-install breakroom breakroom-read ${PREFIX}/bin

clean:
	rm -f zendebug zenexplorer restroom-test
