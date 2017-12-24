BINARY = muh-light-string
GOARCH = arm

all: build deploy

# Build the project
build: clean arm pack

clean:
	-rm -f ${BINARY}*

deploy:
	scp ${BINARY}.up pi:

arm:
	GOOS=linux GOARCH=${GOARCH} GOARM=6 go build -o ${BINARY} .

pack:
	upx -9 -o ${BINARY}.up ${BINARY}

.PHONY: arm build clean deploy pack
