<div align="center">
    <h1>Gin + Postgres Template</h1>
    <h4>Get started fast and efficiently with a ready to go template</h4>
    currently in development...
</div>

## Quickstart Guide
This project runs containerized so there should be no issues getting it off the ground. Please make sure that either your docker daemon is running, or you have Docker Desktop open (which does essentially the same thing.)

### Starting for the first time
To start the project for the first time run:

```cli
docker compose up --build
```

This may take some time depending on your internet connection as it needs to grab the postgres image, as well as build the custom image as part of the project.

While that downloads feel free to open the codebase your editor of choice and start taking a look around the template. We are using Air for hot reloading and there's a volume connected to root of the project so you will be able to edit everything directly.

### Starting the project up in the future
If you've built the project before you can ignore the build flag when starting the project:

```cli
docker compose up
```
However, if something has gone horribly wrong you can try rebuilding with the `--build` flag again. Just keep in mind that this will create a dangling image that's used for the build step so please run the following as clean up:

```cli
docker image prune -f
```

This will remove all the unused dangling images lying around on your system.

### Stopping the project
Once you are done with development and want to shut down the project, first hit `Ctrl + C` in whatever session it is running. Once it has gracefully stopped (or just hard kill it :gigachad:) run the following:

```cli
docker compose down
```

### Data storage
Keep in mind that everything you store in the database is stored locally in a volume, so you should be able to pick up right where you left off with whatever had already been setup.

## Structure
### The Project
The project is structured as Controller Service Repository. The routes directory will maintain all the endpoints a user could call, grouped in their respective files. The same can be said for the controllers and services related to each. The current project directory structure is built as:
```
.
├── Dockerfile
├── README.md
├── controllers
│   └── user-controller.go
├── docker-compose.yml
├── dto
│   └── auth-body.dto.go
├── go.mod
├── go.sum
├── helpers
│   └── response-builder.go
├── initialize
│   ├── connect-database.go
│   └── sync-database.go
├── main.go
├── middleware
│   └── authentication.go
├── models
│   └── users.go
├── routes
│   ├── auth-routes.go
│   └── user-routes.go
└── services
    └── jwt_service.go
```

### Controllers
Will place the controller template structure here when I've built the template file.

### Services 
Will place the service template structure here when I've built the template file.

---

Current Dev tasks
straight coding tasks:
- [x] dockerfile written
- [x] docker compose together
- [x] setup connection between postgres and backend using gorm
- [x] basic user model to test jwt auth
- [x] connected Air to have hot reloading during development
- [x] signing jwt access and refresh tokens
- [x] writing the authentication middleware
- [ ] refresh access token

coding but still vague
- [ ] figure out how unit tests
- [ ] write the runner for the automated unit tests
- [ ] write the service & controller template we will use for each part of the project

research
- [ ] doing some additional research on bcrypt cost (i’m a little worried about it)

doco
- [x] writing a quickstart guide
- [ ] writing doco for all of this


