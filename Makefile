default:
	${MAKE} deps
	${MAKE} build
	${MAKE} test
deps:		
	${MAKE} clean
	glide install
build:
	go build -o ./build/snap-plugin-collector-kafka
test:
	go test ./kafka
clean:
	glide cc && rm -rf ./vendor ./build

