default:
	${MAKE} deps
	${MAKE} build
	${MAKE} test
deps:		
	${MAKE} clean
	glide install
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ./build/snap-plugin-collector-kafka
test:
	go test ./kafka
clean:
	glide cc && rm -rf ./vendor ./build

