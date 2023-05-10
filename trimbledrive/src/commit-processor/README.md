# Trimble Drive Commit Processor

Trimble Drive Commit Processor service written in Go

- Dependency injection using [uber-go/fx](https://pkg.go.dev/go.uber.org/fx)
- Database actions using [mongo](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- Events(kafka) related actions using [kafka-go](https://pkg.go.dev/github.com/segmentio/kafka-go) 
- Code Style based on [Google Go Style Guide](https://google.github.io/styleguide/go/)

## Responsibilities
- This service watches for the messages in `trimble.tdrive.commit_processor` topic
- Performs File Upload Commit Operation
- Publish events to dead letter topic upon failures after retries.

## Project Structure

| Folder/File | Description                                                                                                     |
|-------------|-----------------------------------------------------------------------------------------------------------------|
| main.go     | Entrypoint for the application. Start application by running this file. Contains Fx app initialization.         |
| config      | Stores configuration files and configuration related functions like logging, mongo db configuration.            |
| processor   | Contains the processor file to consume the event from any messaging systems.                                    |
| service     | Contains the functionality that compounds the core of the application. Functions in services must return error. |
| model       | Contains event model or other miscellaneous structs                                                             |
| Dockerfile  | To create image of the service using docker                                                                     |


