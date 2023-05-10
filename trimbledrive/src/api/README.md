# Trimble Drive API Service

Trimble Drive API service written in Go

- Application backbone with [Gin Web Framework](https://github.com/gin-gonic/gin)
- Dependency injection using [uber-go/fx](https://pkg.go.dev/go.uber.org/fx)
- Database actions using [mongo](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- Code Style based on [Google Go Style Guide](https://google.github.io/styleguide/go/)

## Project Structure

| Folder/File     | Description                                                                                                   |
|-----------------|---------------------------------------------------------------------------------------------------------------|
| cmd/app/main.go | Entrypoint for the application. Start application by running this file. Contains Fx app initialization.       |
| config          | Stores configuration files and configuration related functions like logging, mongo db configuration           |
| docs            | Contains docs related to this application like API document using OpenAPI                                     |
| handler         | Contains the handler files for per API group.Input validation and response handling will be done here         |
| internal        | Contains the functionality that compounds the core of the application.Functions in services must return error |
| middleware      | Contains the request/response interceptors                                                                    |
| model           | Contains request/response structs                                                                             |
| router          | Contains router files that routes the API to respective handler                                               |
| server          | Contains gin based server initialization with graceful shutdown                                               |
| utils           | Contains utility functions like error handling helpers                                                        |
| Dockerfile      | To create image of the service using docker                                                                   |
