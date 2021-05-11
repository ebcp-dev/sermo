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

---

Routes:

- Auth routes:
  - GET /users - retrieves list of users
  - GET /user/:id - retrieves a specific user
  - POST /user ({email, password}) - register user with email, password
  - POST /user/login ({email, password}) - user login with email, password
  - PUT /user/:id ({email, password}) - update user details
  - DELETE /user/:id - delete user by id

---

CI/CD:

- Semaphore CI

---

Links:

https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/

https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0
