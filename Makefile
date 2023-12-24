MODULE_NAME=merkle-file-uploader
BINARY_NAME=mfu
TEST_FOLDER=resources

.PHONY: start-server build-client

start-server:
	docker-compose up -d

stop-server:
	docker-compose down

build-client:
	go build -o $(BINARY_NAME) $(MODULE_NAME)/

test-upload: build-client
	./$(BINARY_NAME) client upload $(TEST_FOLDER)

test-download: build-client
	./$(BINARY_NAME) client download 1
