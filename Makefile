all: zendebug zenexplorer

zenexplorer_sources := main.go zencodeStatements.go delegate.go zenroom.go server.go
zendebug_sources := main.go

zenexplorer: $(wildcard src/zenexplorer/*.go)
	go build -o zenexplorer $(foreach source, ${zenexplorer_sources}, src/zenexplorer/${source})

zendebug: $(wildcard src/zendebug/*.go)
	go build -o zendebug $(foreach source, ${zendebug_sources}, src/zendebug/${source})

clean:
	rm -f zendebug zenexplorer
