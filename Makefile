default:
	${MAKE} deps
	${MAKE} build
	${MAKE} test
deps:		
	${MAKE} clean
	glide install
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o ./build/snap-plugin-collector-kafka -ldflags "-w"
test:
	go test ./kafka
clean:
	glide cc && rm -rf ./vendor ./build

