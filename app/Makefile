
.PHONY: all
all: outBin docker
	mkdir -p ./testDataTmp
	mkdir -p ./logsData

outBin:
	go build -o $@ .

.PHONY: test
test:
	@echo "Running test..."
	go test .
#	go test -race $(PKGS)

.PHONY: cover
cover:
	@echo "Running coverage testing..."
#	./scripts/cover.sh $(PKGS)

.PHONY: bench
BENCH ?= .
bench:
#	@$(foreach pkg,$(PKGS),go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS) $(pkg);)
	@echo "Running benchmark..."

.PHONY: updatereadme
updatereadme:
	@echo "README.md updating..."
#	rm -f README.md
#	cat .readme.tmpl | go run internal/readme/readme.go > README.md
#
#observer
.PHONY:docker
docker: outBin
	docker rmi -f myapp:falcon
	docker build --rm -t myapp:falcon \
				--build-arg "group=${USER}" \
				--build-arg "gid=$$(id -u)" \
				--build-arg "user=${USER}" \
				--build-arg	"uid=$$(id -g)" \
				-f docker/Dockerfile.img .

.PHONY: clean
clean:
	rm -rf outBin
	rm -rf ./testDataTmp
	docker rmi -f myapp:falcon


