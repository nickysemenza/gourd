VERSION          := $(shell git describe --tags --always --dirty="-dev")a
DATE             := $(shell date '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

GCP_PROJECT := cloudrun1-278204
GCR_IMAGE := gcr.io/$(GCP_PROJECT)/gourd-backend:$(VERSION)

deploy: deploy-image deploy-run
deploy-image:
	echo "building $(GCR_IMAGE)"
	gcloud builds submit --tag $(GCR_IMAGE)
deploy-run:
	echo "deploying $(GCR_IMAGE)"
	terraform apply -var-file="prod.tfvars" -var="image_name=$(GCR_IMAGE)" -var="project_id=$(GCP_PROJECT)"



.PHONY: all
all: bin/gourd

bin/%: $(shell find . -type f -name '*.go' | grep -v '_test.go')
	@mkdir -p $(dir $@)
	go build $(VERSION_FLAGS) -o $@ ./cmd/$(@F)

test: unit-test lint

dev: bin/gourd
	./bin/gourd server

dev-env:
	docker-compose up -d db jaeger
dev-air: bin/air 
	HTTP_HOST=127.0.0.1 ./bin/air -c air.conf
dev-ui:
	cd ui && yarn run start

bin/air:
	@mkdir -p $(dir $@)
	go build -o $@ ./vendor/github.com/cosmtrek/air

bin/golangci-lint:
	@mkdir -p $(dir $@)
	go build -o $@ ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint
bin/migrate:
	@mkdir -p $(dir $@)
	go build -tags 'postgres' -o $@ ./vendor/github.com/golang-migrate/migrate/v4/cmd/migrate
bin/go-acc:
	@mkdir -p $(dir $@)
	go build -o $@ ./vendor/github.com/ory/go-acc

unit-test: 
	go test -v -race -cover ./...
integration-test: 
	go test -v -race -cover ./... -tags integration	
test-full-cover: bin/go-acc
	./bin/go-acc -o coverage-full.txt ./... -- -race -tags integration
	
lint: bin/golangci-lint
	bin/golangci-lint run || (echo "lint failed"; exit 1)	

IMAGE := nicky/gourd-backend
docker-build:
	docker build -t $(IMAGE) .
docker-dev: docker-build

docker-push: docker-build
	docker push $(IMAGE):latest

dev-db:
	pgcli postgres://gourd:gourd@localhost:5555/food
dev-db-stats:
	docker logs gourd_db_1 2>&1 | pgbadger - --prefix '%t [%p]:[%l] user=%u, db=%d'
new-migrate/%: bin/migrate
	mkdir -p db/migrations
	./bin/migrate create -dir db/migrations -ext sql $(@F)
migrate: bin/migrate
	./bin/migrate -source file://db/migrations -database postgres://gourd:gourd@localhost:5555/food?sslmode=disable up
migrate-down: bin/migrate
	./bin/migrate -source file://db/migrations -database postgres://gourd:gourd@localhost:5555/food?sslmode=disable down

generate: generate-rust generate-openapi

generate-rust:
	cd rust && wasm-pack build

validate-openapi:
	./ui/node_modules/ibm-openapi-validator/src/cli-validator/index.js api/openapi.yaml -c api/.validaterc -v
generate-openapi: validate-openapi
	rm -rf ui/src/api/openapi-fetch
	rm -rf ui/src/api/openapi-hooks
	openapi-generator generate -i api/openapi.yaml -o ui/src/api/openapi-fetch -g typescript-fetch --config ui/openapi-typescript.yaml
	openapi-generator generate -i api/openapi.yaml -o rust/openapi -g rust --global-property models,supportingFiles,modelDocs=false


	mkdir -p ui/src/api/openapi-hooks/
	cd ui && yarn run generate-fetcher

	go generate ./api
cy:
	cd ui && yarn run cy:open

# https://jqplay.org/s/c1T3lLCJwH
get-detail/%: 
	curl -s http://localhost:4242/api/recipes/$(@F) | jq '.detail | del(.. | .id?)' > tmp1
	dyff yaml tmp1 | pbcopy
	rm tmp1

seed-testdata: bin/gourd
	./testdata/seed.sh

rs: 
	cd rust/s && cargo sqlx prepare -- --bin gourd