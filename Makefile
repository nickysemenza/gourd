VERSION          := $(shell git describe --tags --always --dirty="-dev")a
DATE             := $(shell date '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

DSN := postgres://gourd:gourd@localhost:5555/food?sslmode=disable

.PHONY: all
all: bin/gourd

# general dev
test: unit-test-go lint-go test-rs

dev-env:
	docker-compose -f dockerfiles/docker-compose-deps.yml -f dockerfiles/docker-compose-devdeps.yml up -d db meilisearch meilisearch-ui jaeger collector
dev-air: bin/air 
	HTTP_HOST=127.0.0.1 ./bin/air -c tooling/dev/air.conf

install-deps:
	go mod tidy
	go mod vendor
	go mod tidy
	cd ui && yarn

# frontend dev
dev-ui:
	cd ui && yarn run dev
cy:
	cd ui && yarn run cy:open

# rust dev
dev-rs:
	cd rust && RUST_BACKTRACE=1 RUST_LOG=html5ever=info,selectors=info,debug cargo watch -x 'run server'

test-rs:
	cd rust && cargo test

# db
dev-db:
	pgcli $(DSN)
dev-db-stats:
	docker logs gourd_db_1 2>&1 | pgbadger - --prefix '%t [%p]:[%l] user=%u, db=%d'
new-migrate/%: bin/migrate
	mkdir -p internal/db/migrations
	./bin/migrate create -dir internal/db/migrations -ext sql $(@F)
migrate: bin/migrate
	./bin/migrate -source file://internal/db/migrations -database $(DSN) up
migrate-down: bin/migrate
	./bin/migrate -source file://internal/db/migrations -database $(DSN) down


# golang
bin/%: $(shell find . -type f -name '*.go' | grep -v '_test.go')
	@mkdir -p $(dir $@)
	CGO_ENABLED=0 go build $(VERSION_FLAGS) -o $@ ./internal/cmd/$(@F)

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
bin/oapi-codegen:
	@mkdir -p $(dir $@)
	go build -o $@ ./vendor/github.com/deepmap/oapi-codegen/cmd/oapi-codegen

unit-test-go: 
	go test -v -race -cover ./...
integration-test-go: bin/go-acc 
	./bin/go-acc -o coverage-integration.txt ./... -- -race -tags integration -v
unit-cover-go: bin/go-acc
	./bin/go-acc -o coverage-full.txt ./... -- -race
	
lint-go: bin/golangci-lint
	bin/golangci-lint run || (echo "lint failed"; exit 1)	


generate: wasm-dev openapi # gen-db



.PHONY: openapi
openapi: internal/api/openapi.yaml
internal/api/openapi.yaml: tooling/schemas/gourd.yaml tooling/schemas/usda.yaml
	# generate merged bundle
	npx @redocly/openapi-cli bundle tooling/schemas/gourd.yaml --output internal/api/openapi.yaml 
	
	# ui hooks 1
	cd ui && npx openapi-codegen gen gourdApi
	
	# rust
	rm -rf rust/openapi/src/models
	npx @openapitools/openapi-generator-cli generate -i internal/api/openapi.yaml \
		-o rust/openapi -g rust --global-property models,supportingFiles,modelDocs=false
	cd rust/openapi && cargo clippy --fix --allow-staged --allow-dirty

	# go
	go generate ./internal/api

openapi-docs:
	npx @redocly/openapi-cli preview-docs internal/api/openapi.yaml -p 8081

# WASM generation
.PHONY: generate-wasm
generate-wasm: rust/wasm/pkg

rust/wasm/pkg: $(shell find rust/ -type f -name '*.rs')
	cd rust && wasm-pack build wasm

ui/src/wasm/package.json: rust/wasm/pkg/package.json
	cp -r rust/wasm/pkg/ ui/src/wasm/
wasm-dev: generate-wasm ui/src/wasm/package.json
	

# GO SQL GENERATION
gen-db:
	rm -rf db/models/
	sqlboiler psql --relation-tag rel --config ./tooling/sqlboiler.toml --output ./internal/db/models --add-soft-deletes
# misc dev

# https://jqplay.org/s/c1T3lLCJwH
get-detail/%: 
	curl -s http://localhost:4242/api/recipes/$(@F) | jq '.detail | del(.. | .id?)' > tmp1
	dyff yaml tmp1 | pbcopy
	rm tmp1

seed-testdata: bin/gourd
	./tooling/testdata/seed.mjs

# todo import usda
devdata: seed-testdata sync
sync: bin/gourd
	./bin/gourd sync
insert-album: 	
	PGPASSWORD=gourd psql -Atx "$(DSN)" -h localhost -U gourd -d food -p 5555 -c "INSERT INTO "public"."gphotos_albums" ("id", "usecase") VALUES ('AIbigFomDsn4esVUopzvXsZ5GDjY3EDb7L_A8sf1Wf7-IWHxykoMjVy-KeCTHW7nVIaTkJ8CAV8i', 'food'), ('AIbigFobaVQFOEyoZk2TKEkCNS-ffzesGu7n-OZy6-YXnKLgrYE4ALSW-LknhcbttNNifPPCm7sY','plants');"


# GCP_PROJECT := cloudrun1-278204
# GCR_IMAGE := gcr.io/$(GCP_PROJECT)/gourd-backend:$(VERSION)
# IMAGE := nicky/gourd-backend
# deploy: deploy-image deploy-run
# deploy-image:
# 	echo "building $(GCR_IMAGE)"
# 	gcloud builds submit --tag $(GCR_IMAGE)
# deploy-run:
# 	echo "deploying $(GCR_IMAGE)"
# 	terraform apply -var-file="prod.tfvars" -var="image_name=$(GCR_IMAGE)" -var="project_id=$(GCP_PROJECT)"
# docker-build:
# 	docker build -t $(IMAGE) .

# docker-push: docker-build
# 	docker push $(IMAGE):latest	