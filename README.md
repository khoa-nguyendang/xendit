# xendit
This project to test golang

# Prequiresites
### Golang
following tutorial at:
```
https://go.dev/doc/install
```
<br/>

### Protobuf
```
sudo apt install -y protobuf-compiler
```
<br/>

### GRPC

```
# readmore : https://grpc.io/docs/languages/go/quickstart/
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

```
<br/>

### CMake
following tutorial at:
```
https://cmake.org/install/
```
<br/>

### Docker and docker-compose
following tutorial at:
```
#docker
https://docs.docker.com/engine/install/ubuntu/

#compose
https://docs.docker.com/compose/install/
```

<br/>


<br/>


# Introduction

### Architecture and Patterns

 1. [CleanCode](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) Main architecture base on CleanCode.
 2. [Service-Endpoint-Transport](https://microservices.io/patterns/microservice-chassis.html) Pattern: Microservice chassis.

<br/>


### Opensource libraries

 1. [GRPC](https://grpc.io/) &emsp; -  gRPC
 2. [go-sql-driver](https://github.com/go-sql-driver/mysql) mysql driver.
 3. [zap](https://github.com/uber-go/zap) &emsp; -  performance logging library.
 4. [sqlx](https://github.com/jmoiron/sqlx) Extensions to database/sql.
 5. [viper](https://github.com/spf13/viper) Realtime configuration management.
 6. [Jaeger](https://www.jaegertracing.io/) &emsp; -  Jaeger tracing, tracing transactions between distributed services.
 7. [migrate](https://github.com/golang-migrate/migrate) &emsp; -  Database migration for either MySql, MongoDB.


### Troubleshouting


### First Deloyment tutorial

Step 1. Prepare new necessary folder, stop current running to load latest deployment from remote
```
make prepare
docker-compose down
```
### and for upward deployment

Step 3. Pull latest source, remember to stash your changes unless you may lost it.
```
git reset --hard HEAD
git pull
```

Step 4. Build docker-compose
```
docker-compose build
```

Step 5. Start docker-compose
```
docker-compose up -d
```

Step 6. View process
```
clear
docker ps
```

Step 7. Access swagger
```
http://localhost:3000/swagger/index.html
```

#### Stop current process
```
docker-compose down
```

#### display logs for specific service
```
docker logs name_of_service_here
#example : docker logs transaction-service
```

#### access to specific service
```
docker exec -it name_of_service_here bash
#example : docker exec transaction-service bash
```

#### generate swagger
```
cd cmd
swag init --parseDependency --parseInternal --parseDepth 1 
```