## Sermo

A chat server written in Go with Gorilla WebSocket and JWT authentication.

[![Build Status](https://ebcp-dev.semaphoreci.com/badges/sermo/branches/master.svg?style=shields&key=eeebee0b-69c4-4904-9e70-dc9c7e8f6ffd)](https://ebcp-dev.semaphoreci.com/projects/sermo)

Made with:

- Go
- PostgreSQL
- Mux - web framework.
  - https://github.com/gorilla/mux
- Gorilla WebSocket - WebSocket implementation in Go.
  - https://github.com/gorilla/websocket
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

- Channel routes:
  - [GET] /channel/:id - retrieves a specific channel
  - [POST] /channel (Auth required) - register channel with string, int attributes
    - {channelname, maxpopulation}
  - [GET] /channel (Auth required) - retrieves list of channel
  - [PUT] /channel/:id (Auth required) - update channel details
    - {strattr, intattr}
  - [DELETE] /channel/:id (Auth required) - delete channel by id

---

Links:

- Guides:

  - Go WebRTC:

    - https://medium.com/@ramezemadaiesec/from-zero-to-fully-functional-video-conference-app-using-go-and-webrtc-7d073c9287da
    - https://github.com/pion/webrtc/blob/master/examples/README.md
    - https://github.com/pion/example-webrtc-applications
    - https://webrtcforthecurious.com/

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

Known Issues:

- Update query goes through even if id doesn't exist.
- Delete query goes through even if id doesn't exist.
