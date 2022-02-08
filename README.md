# Endpoints

The project uses `swaggo` to generate an `OpenAPI` specification for each endpoint. To view this documentation using swagger run `go run .` and visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html#/).
 (Note: you may have to recompile the docs with `swag init`)

# Tests

This project contains both unit and integration tests. Unit tests are located within each package with the suffix `_test.go` and they test particular methods or functions.
Integrations tests are defined in files in the `app` directory, divided by a per resource basis. These tests execute different requests to different endpoints using also other resources and validating the entire lifecycle of the application.
The tests files are named with the format `<resource_name>_test.go`.

Tests can be run on a per package basis with the command `go test <package>`. For example if we want to run the controller tests we can do `go test ./app/controllers` or `go test ./app` if we want to run
only the integration tests. All tests from all packages can be running by calling `go test ./...`
# Packages
The project is divided in different packages, inside each package `_test.go` files represent unit tests.

## app
The main package, here we define our `App` structure which holds
 all the relevant server data like the db connection and the router.
 
 
`app_test.go` contains useful function for our integration tests
 
## app/controller

The code here is separated on different files depending on the resource. 
Each file contains the appropiate methods which are then mapped to endpoint in `app.go`
 
## app/repos

This package has our data abstractions using the repository pattern.

## app/middleware

This package holds all of our middlewares, these are used in conjunction with
the router in order to provide auth validation or admin validation 
on specific endpoints
  
## app/models

This package holds all of our database models and response models.

# Authentication

When a user registers, the password a hash is stored using `bcrypt`, later when a user
tries to login the password is validated against that hash and a `JWT` token is emitted
in order to validate the session.

Registering happens on the endpoint `POST api/user` while login occurs in `POST api/session`

# Migrations

Migrations are executed automatically when the app starts. See the `runner.go` file in the `app/migrations` package.

# Environment variables
A `.env` file should be created in the root directory with the following environmental variables defined.
```
DB_NAME=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
TEST_DB_NAME=
TEST_DB_HOST=
TEST_DB_PORT=
TEST_DB_USER=
TEST_DB_PASSWORD=
JWT_SECRET=
 ```
