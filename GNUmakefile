TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
SHELL:=/bin/bash

# For local testing, run `docker-compose up -d`
SERVER ?=http://127.0.0.1:8153/go/
export GOCD_URL=$(SERVER)
export GOCD_SKIP_SSL_CHECK=1


travis: before_install script after_success deploy_on_develop

before_install:
	go get -t -v ./...
	go get -u github.com/golang/lint/golint
	go get -u github.com/kardianos/govendor
	go get github.com/sergi/go-diff/diffmatchpatch


script: test
	chmod -R 777 ./godata/server
	make testacc

teardown_docker:
	docker-compose down

after_failure: teardown_docker upload_logs

after_success: teardown_docker upload_logs
	bash <(curl -s https://codecov.io/bash)
	go get github.com/goreleaser/goreleaser

deploy_on_tag:
	go get github.com/goreleaser/goreleaser
	gem install --no-ri --no-rdoc fpm
	go get
	goreleaser

deploy_on_develop:
	go get github.com/goreleaser/goreleaser
	gem install --no-ri --no-rdoc fpm
	go get
	goreleaser --snapshot

upload_logs:
	pip install awscli --upgrade --user
	AWS_DEFAULT_REGION=$(ARTIFACTS_REGION) \
		AWS_ACCESS_KEY_ID=$(ARTIFACTS_KEY) \
		AWS_SECRET_ACCESS_KEY=$(ARTIFACTS_SECRET) \
		aws s3 sync ./godata/server/ s3://$(ARTIFACTS_BUCKET)/drewsonne/terraform-provider-gocd/$(TRAVIS_BUILD_ID)/godata/

default: build

build: fmtcheck
	go install

test: fmtcheck before_install
	go test -v -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -v -timeout=30s -parallel=4

testacc: provision-test-gocd fmtcheck
	bash scripts/wait-for-test-server.sh
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./gocd"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

provision-test-gocd:
	docker-compose up -d

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
