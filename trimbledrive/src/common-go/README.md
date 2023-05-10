# Common-Go
This is a common directory that contains common functions as go modules that can be used by multiple sub services.

In order to make use of these modules, add an entry referring the relative path of the directory in `go.work` file in the root directory.

## Project Structure

| Folder        | Description                                                                        |
|---------------|------------------------------------------------------------------------------------|
| exception     | Go Module as a library for handling exceptions and error codes.                    |
| logger        | Go Module as a library for logging configurations.                                 |
| repository    | Go Module as a library for performing database operations.                         |
| configuration | Go Module as library for client configuration of monogo, kafka, logging and cosmos |
