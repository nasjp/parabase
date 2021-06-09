.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: login
login: ## login docker container
	@echo '<login>'
ifeq ($(IS_DOCKER), true)
	@:
else
	@docker compose run --rm app bash
endif

.PHONY: wait-db
wait-db: ## wait for the database to wake up
	@echo '<wait-db>'
ifeq ($(IS_DOCKER), true)
	@sh wait.sh
else
	@docker compose run --rm app sh -c "sh wait.sh"
endif


.PHONY: cleanup-test-db
cleanup-test-db: wait-db ## cleanup test db
	@echo '<cleanup-test-db>'
ifeq ($(IS_DOCKER), true)
	@mysql -u root --password=password -e 'show databases' | grep _test | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'
	@mysql -u root --password=password -e 'show databases' | grep _management | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'
else
	@docker compose exec mysql bash -c "mysql -u root --password=password -e 'show databases' | grep _test | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'"
	@docker compose exec mysql bash -c "mysql -u root --password=password -e 'show databases' | grep _management | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'"
endif

.PHONY: test
test: cleanup-test-db ## go test
	@echo '<test>'
ifeq ($(IS_DOCKER), true)
	@go test ./... -race -cover
else
	@docker-compose run --rm app bash -c "go test ./... -race -cover"
endif

.PHONY: lint
lint: ## go vet
	@echo '<lint>'
	@go vet ./...

.PHONY: build
build: ## go build
	@echo '<build>'
	@go build .
	@go build ./mysql

.PHONY: clean
clean: ## clean bin
	@echo '<clean>'
	@go clean -testcache

.PHONY: ci
ci: clean lint test build ## Continuous Integration

NUMBER=1 2 3 4 5 6 7 8 9 10 \
			 11 12 13 14 15 16 17 18 19 20 \
			 21 22 23 24 25 26 27 28 29 30 \
			 31 32 33 34 35 36 37 38 39 40 \
			 41 42 43 44 45 46 47 48 49 50 \
			 51 52 53 54 55 56 57 58 59 60 \
			 61 62 63 64 65 66 67 68 69 70 \
			 71 72 73 74 75 76 77 78 79 80 \
			 81 82 83 84 85 86 87 88 89 90 \
			 91 92 93 94 95 96 97 98 99 100

.PHONY: test-heavy
test-heavy: cleanup-test-db ## 100 times test
	echo '<test-heavy>'
	@for i in ${NUMBER}; do go test ./mysql/... -race -count=1 -v ; done
