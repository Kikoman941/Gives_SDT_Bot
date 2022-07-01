.PHONY:

build:
	go build -o ./.bin/app cmd/app/main.go

run: build
	./.bin/app

lint:
	golangci-lint

define ENV_EXAMPLE
IS_PROD=False
SUPERADMIN=0
BOT_TOKEN=token
BOT_USERNAME=tg_username
BOT_POLLING_TIMEOUT=30s
POSTGRESQL_DSN=dsn
endef

export ENV_EXAMPLE
env_example:
	@if [ ! -f ".env" ]; then \
  		echo "$$ENV_EXAMPLE" > ".env"; \
  	fi

