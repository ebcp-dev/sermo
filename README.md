## Sermo

A REST API written in Go with a PostgreSQL database and JWT authentication.

Made with:

- Go
- PostgreSQL
- Mux - web framework.
  - https://github.com/gorilla/mux
- pq - PostgreSQL driver for Go.
  - https://github.com/lib/pq
- JWT - used for authenticating users.
  - https://github.com/dgrijalva/jwt-go
- UUID - used for parsing and storing UUID in database.
  - https://github.com/google/uuid
- Viper - used for configuring environment variables
  - https://github.com/spf13/viper
- Docker
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
  - [PUT] /data/:id (Auth required) - update data details
    - {strattr, intattr}
  - [DELETE] /data/:id (Auth required) - delete data by id

---

Links:

- Guides:

  - Building a basic REST API with testing and CI/CD:

    - https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

  - Authentication with JWT:

    - https://tutorialedge.net/golang/authenticating-golang-rest-api-with-jwts/
    - https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0

  - Deployment:
    - https://tutorialedge.net/golang/go-docker-tutorial/
    - https://semaphoreci.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker
    - https://towardsdatascience.com/use-environment-variable-in-your-next-golang-project-39e17c3aaa66

- Docs:
  - Syntax and convention:
    - https://golang.org/doc/effective_go
    - https://medium.com/@kdnotes/golang-naming-rules-and-conventions-8efeecd23b68
    - https://devhints.io/go
