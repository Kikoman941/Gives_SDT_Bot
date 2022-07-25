.PHONY:

build:
	go build -o ./.bin/app cmd/app/main.go

run: build
	./.bin/app

lint:
	golangci-lint run -c .golangci.yaml

define ENV_EXAMPLE
IS_PROD=False
SUPERADMIN=0
BOT_TOKEN=token
BOT_POLLING_TIMEOUT=30s
PUBLISHER_TIMEOUT=5m
POSTGRESQL_DSN=dsn
endef

export ENV_EXAMPLE
env_example:
	@if [ ! -f ".env" ]; then \
  		echo "$$ENV_EXAMPLE" > ".env"; \
  	fi

docker_build:
	docker build -t SDT_Gives_Bot .

docker_run:
	docker run \
		--name SDT_Gives_Bot \
		-v "${PWD}"/.images:/.images \
		-d SDT_Gives_Bot