build:
	CGO_ENABLED=0 go build -o ../target/data-collector ./src/*.go

build/docker:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o /data-collector ./src/*.go