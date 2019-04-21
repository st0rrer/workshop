build:
	CGO_ENABLED=0 go build -o ../target/web-service ./src/*.go

build/docker:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /web-service ./src/*.go