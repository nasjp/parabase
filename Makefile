.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: cleanup
cleanup: ## cleanup test db
	docker compose exec db bash -c "mysql -u root --password=password -e 'show databases' | grep _test | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'"
	docker compose exec db bash -c "mysql -u root --password=password -e 'show databases' | grep _management | xargs -I DB mysql --password=password -e 'DROP DATABASE IF EXISTS DB'"

.PHONY: test
test: ## go test
ifeq ($(IS_DOCKER), true)
	go test ./... -race -v
else
	docker compose run --rm app bash -c "go test ./... -race -v"
endif
