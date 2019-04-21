build:
	CGO_ENABLED=0 go build -o ../target/client ./src/*.go

build/docker:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /client ./src/*.go