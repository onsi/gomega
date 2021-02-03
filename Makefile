
.PHONY: test
test:
	[ -z "`gofmt -s -w -l -e .`" ]
	go vet
	go run github.com/onsi/ginkgo/ginkgo -p -r --randomizeAllSpecs --failOnPending --randomizeSuites --race

docker_test:
	docker-compose build test && docker-compose run --rm test make test
