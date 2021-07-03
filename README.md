# Endpoints

# Packages
The project is divided in different packages, inside each package `_test.go` files represent unit tests.

## app
The main package, here we define our `App` structure which holds
 all the relevant server data like the db connection and the router.
 
`app_test.go` contains our integration tests
 
## app/handlers

The code here is separated on different files depending on the endpoint. 
Each file contains a function to register all the routes for a specific REST
 resource, they also handle validation and response codes.
 
## app/repos

This package has our data abstractions using the repository pattern.

## app/middleware

This package holds all of our middlewares, these are used in conjunction with
the router in order to provide auth validation or admin validation 
on specific endpoints
  
## app/models

This package holds all of our database models.

# Authentication

When a user registers, the password a hash is stored using `bcrypt`, later when a user
tries to login the password is validated against that hash and a `JWT` token is emitted
in order to validate the session.

Registering happens on the endpoint `POST api/user` while login occurs in `POST api/session`

# Migrations

Migrations are executed automatically when the app starts. See the `runMigrations` function on `app.go`
