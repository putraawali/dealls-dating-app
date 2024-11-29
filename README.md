# Dating Apps

Simple dating apps with feature:

-   Register
-   Login
-   Verify Email
-   Get Available Partner
-   Swipe Left/Right Partner
-   Request Transaction for Premium User
-   Accept Transaction for Premium User

## How To Run Service

-   Create file .env and copy all key from file **.env.example**
-   run `go run main.go` on terminal from root project directory

## Project Structure

```
📦src
 ┣ 📂constants ------ Constants Variable
 ┣ 📂controllers ---- Controllers layer, handler function for each endpoint
 ┣ 📂dtos ----------- Data Transfer Object, mostly for Request and Response Body Params
 ┣ 📂mocks ---------- Mock dependency injection
 ┣ 📂models --------- Database model
 ┣ 📂pkg ------------ Package used by this project such as db connection, jwt, helpers, etc.
 ┃ ┣ 📂connections -- Databases connection
 ┃ ┃ ┣ 📂mocks ------ Mock of databases connection
 ┃ ┣ 📂helpers ------ Helper
 ┃ ┣ 📂jwt ---------- JWT Auth
 ┃ ┣ 📂middlewares -- Endpoint Middlewares
 ┃ ┗ 📂response ----- Custom response handler
 ┣ 📂repositories --- Repository layer, for query to databases
 ┃ ┣ 📂mocks -------- Mock of repository layer
 ┣ 📂usecases ------- Usecase/Service layer, for business logic
 ┃ ┣ 📂helpers ------ Helpers for usecase only
 ┃ ┣ 📂mocks -------- Mocks of usecases
 ┣ 📜dependency.go -- Dependency injection
 ┣ 📜module.go ------ Module to wrap Dependency Injection, Middlewares and Routes
 ┗ 📜routes.go ------ Router
```

## Additional Features

-   Log Request and Response
-   Log error with stack trace
-   Each request has unique request-id if no request-id from request headers
