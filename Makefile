cnf ?= Makefile.env
include $(cnf)

build:
	docker build -f Dockerfile -t $(IMG_NAME):$(IMG_TAG) .	

build-debug: 
	docker build -f Dockerfile.debug -t $(IMG_NAME).debug:$(IMG_TAG) .	

run:
	docker run -d -v"$(CURDIR)/$(APP_CFG_FILE):/tmp/config.cfg" -e IMSCFGFILE=/tmp/config.cfg -p $(APP_PORT):8080 --name $(APP_NAME) $(IMG_NAME):$(IMG_TAG)

run-debug:
	docker run -d -v"$(CURDIR)/$(APP_CFG_FILE):/tmp/config.cfg" -e IMSCFGFILE=/tmp/config.cfg -p $(APP_PORT):8080 -v"$(CURDIR)/.app:/app" --name $(APP_NAME) $(IMG_NAME).debug:$(IMG_TAG)

attach-net:
	docker network connect $(net) $(APP_NAME) 

run-db:
	docker run -d --env-file "$(CURDIR)/.config/example-config-srv-users-db.env" -v"ztm-srv-users-db-lib:/var/lib/mysql" -v"ztm-srv-users-db-log:/var/log/mysql" -p 13306:3306 -p 23060:33060 --name ztm-srv-users-db ztm-sql-database:1.0

all: build run-db run
