## GO REST API

A REST API written in Go with a PostgreSQL database and JWT authentication.

[![Build Status](https://ebcp-dev.semaphoreci.com/badges/gorest-api/branches/master.svg)](https://ebcp-dev.semaphoreci.com/projects/gorest-api)

Made with:

- Go
- PostgreSQL
- Mux
  - https://github.com/gorilla/mux
- JWT - used for authenticating users.
  - https://github.com/dgrijalva/jwt-go
- pq - PosgreSQL driver for Go.
  - https://github.com/lib/pq
- UUID - used for parsing UUID.
  - https://github.com/google/uuid
- Semaphore - CI/CD

---

Routes:

- User routes:

  - [POST] /user - register user with email, password
    - {email, password}
  - [POST] /user/login - user login with email, password
    - {email, password}
  - [GET] /user/:id - retrieves a specific user
  - [GET] /users (Auth required) - retrieves list of users
  - [PUT] /user/:id (Auth required) - update user details
    - {email, password}
  - [DELETE] /user/:id (Auth required) - delete user by id

- Data routes:
  - [GET] /data/:id - retrieves a specific data
  - [POST] /data (Auth required) - register data with string, int attributes
    - {strattr, intattr}
  - [GET] /data (Auth required) - retrieves list of data
  - [PUT] /data/:id (Auth required) - update user details
    - {strattr, intattr}
  - [DELETE] /data/:id (Auth required) - delete data by id

---

CI/CD:

- Semaphore CI

---

Links:

- Guides:

  - Building a basic REST API with testing and CI/CD: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

  - Authentication with JWT:
    - https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/
    - https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0

- Docs:
  - Syntax and convention:
    - https://golang.org/doc/effective_go
    - https://medium.com/@kdnotes/golang-naming-rules-and-conventions-8efeecd23b68
    - https://devhints.io/go
