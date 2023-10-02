# Pismo Test
This application runs an API to handle financial transactions.

### Requirements
- [Go](https://go.dev/)

Or you can just use Docker:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/compose-file/)

### Running with docker
First, you need to build the Docker container:
```bash
./script/build
```

This command is going to create an image called `pismo-api` which can be used to run the tests and the application itself.

Then, you can just run using the docker-compose:
```bash
docker-compose up
```

The application will be acessible through `http://localhost:3000` endpoint.

### Testing
You can run the tests with docker by running:
```bash
./script/test
```

Otherwise, you can just use the regular go command to run them:
```bash
go test -v ./...
```

### Migrations
This project uses the [golang-migration](https://github.com/golang-migrate/migrate) tool to track changes to the database schema.

To create a new migration, just run the following command and update the SQL up and down files accordingly:
```bash
migrate create -ext sql -dir db/migrations -seq <migration_name>
```

### Documentation

For Transaction API development, I have used Docker-compose, Postgres16, Golang. Postman for testing.
Code has following structure:
 ```
    app: create database insance using db folder sql script and starting the transaction endpoint to accept connection request.
    db: Contains the db table creation, insertion sql flies.
    handler: call to actual api endpoint reaches and validation done for account and transaction.
    model: account and transaction struct element.
    pkg/lib: helper function.
    repository: interface defined for db call and db function call defind.
    script: to start and test the code.
```
