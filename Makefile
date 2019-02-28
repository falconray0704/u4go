
PKGS ?= github.com/falconray0704/u4go \
		github.com/falconray0704/u4go/app/cfg

#PKGS ?= $(shell glide novendor)


.PHONY: all
all: env
	@echo ""
	@echo "----------------------------------------"
	@echo ""
	@echo "Env built up, and supported entries are:"
	@echo "make test"
	@echo "make cover"
	@echo "make clean"
	@echo ""
	@echo "----------------------------------------"
	@echo ""

.PHONY: dependencies
dependencies:
	go get -u go.uber.org/zap
	go get -u github.com/stretchr/testify/assert

.PHONY: env
env:
	rm -rf ./logsData
	mkdir -p ./logsData
#	echo "" > ./logsData/log.root
#	$$(sudo chown root:root ./logsData/log.root)
#	$$(sudo chmod a-w ./logsData/log.root)

.PHONY: test
test: env
	@echo "Cleanning testing cache..."
	go clean -cache $(PKGS)
	@echo "Running test..."
#	go test -race $(PKGS)
	go test $(PKGS)

.PHONY: testTravisCI
testTravisCI:
	@echo "Running test..."
	go test $(PKGS)

.PHONY: cover
cover: env
	@echo "Cleanning testing cache..."
	go clean -cache $(PKGS)
	@echo "Running coverage testing..."
	./scripts/cover.sh $(PKGS)

.PHONY: bench
BENCH ?= .
bench:
#	@$(foreach pkg,$(PKGS),go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS) $(pkg);)
	@echo "Running benchmark..."

.PHONY: clean
clean:
	rm -rf ./logsData
	rm -rf ./cover.out
	rm -rf ./cover.html
	go clean -testcache $(PKGS)

