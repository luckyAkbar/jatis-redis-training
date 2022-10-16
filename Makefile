SHELL:=/bin/bash

ifdef test_run
	TEST_ARGS := -run $(test_run)
endif

migrate_up=go run main.go migrate --direction=up --step=0
migrate_down=go run main.go migrate --direction=down --step=0
run_command=go run main.go server

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf

lint: check-cognitive-complexity
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=revive --enable=goimports  --enable=unconvert --enable=unparam --concurrency=2

check-modd-exists:
	@modd --version > /dev/null	

check-cognitive-complexity:
	find . -type f -name '*.go' -not -name "*.pb.go" -not -name "mock*.go" -not -name "generated.go" -not -name "federation.go" \
      -exec gocognit -over 15 {} +

migrate:
	@if [ "$(DIRECTION)" = "" ] || [ "$(STEP)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run main.go migrate --direction=$(DIRECTION) --step=$(STEP);\
    fi