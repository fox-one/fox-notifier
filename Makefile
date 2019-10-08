.PHONY: fox clean dist

CMD = $@
ENV = $*

VERSION = 0.0.1
BUILD 	= $(shell git rev-parse --short HEAD)
CONFIGDATA 	= $(shell cat config.${ENV}.yaml | base64)

GO = go

%-dev: GO = GOOS=linux GOARCH=amd64 go
%-prod: GO = GOOS=linux GOARCH=amd64 go

fox-%: out = notifier.${ENV}
fox-%:
	${GO} build -ldflags "-w -s 		\
	-X 'main.VERSION=${VERSION}' 		\
	-X 'main.BUILD=${BUILD}' 			\
	-X 'main.CONFIGDATA=${CONFIGDATA}' \
	" -o ${out}
	@echo ""
	@echo `file ${out}`, `du -h ${out} | cut -f1`

fox: fox-local
release: fox-prod

clean:
	${GO} clean
	rm -f ./notifier.*
