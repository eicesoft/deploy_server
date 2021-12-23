BINARY_NAME=deploy_server

build:
	go build -gcflags '-N -l -w -s' -o $(BINARY_NAME) -v

run:
	./$(BINARY_NAME)

build-image:
	docker build -t eicesoft/deploy-server:1.0.0 .
