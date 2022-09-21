all: zendebug zenexplorer

zenexplorer_sources := src/zenexplorer/main.go src/zenexplorer/zencodeStatements.go src/zenexplorer/delegate.go
zendebug_sources := src/zendebug/main.go

zenexplorer: ${zenexplorer_sources}
	go build -o zenexplorer ${zenexplorer_sources}

zendebug: ${zendebug_sources}
	go build -o zendebug ${zendebug_sources}

clean:
	rm -f zendebug zenexplorer
