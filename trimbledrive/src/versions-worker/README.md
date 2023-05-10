# Trimble Drive Version Worker

Trimble Drive Version Worker service written in Go

- Dependency injection using [uber-go/fx](https://pkg.go.dev/go.uber.org/fx)
- Database actions using [mongo](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- Code Style based on [Google Go Style Guide](https://google.github.io/styleguide/go/)

## Responsibilities
- This service watches the `latest` collection with the help of kafka connectors. 
- Process the events and pushes the change document to `versions` collection.
- Publish events to dead letter topic upon failures after retries.

## Project Structure

| Folder/File | Description                                                                                                     |
|-------------|-----------------------------------------------------------------------------------------------------------------|
| main.go     | Entrypoint for the application. Start application by running this file. Contains Fx app initialization.         |
| config      | Stores configuration files and configuration related functions like logging, mongo db configuration.            |
| worker      | Contains the worker file to consume the event from any messaging systems.                                       |
| service     | Contains the functionality that compounds the core of the application. Functions in services must return error. |
| model       | Contains event model or other miscellaneous structs                                                             |
| Dockerfile  | To create image of the service using docker                                                                     |


