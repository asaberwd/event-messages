SERVICE = event-messages

clean:
	rm -rf ./bin

test:
	go test ./...

build: clean
	GOOS=linux GOARCH=amd64 go build -v -a -tags aws_lambda -o bin/$(SERVICE)-api -a --ldflags "-w \
	-X github.com/asaberwd/event-messages/build.GitCommit=$(COMMIT)" cmd/lambda/main.go

deploy-local: build
	sls deploy --stage local --verbose

deploy-test: build
	sls deploy --stage test --verbose
