# UploadFunctionApp

<!-- TOC -->
* [UploadFunctionApp](#uploadfunctionapp)
    * [ThumbnailObserver](#thumbnailobserver)
    * [FileActivator](#fileactivator)
    * [PackageAggregator](#packageaggregator)
    * [File Common Package](#file-common-package)
    * [Project Structure](#project-structure)
    * [Dependencies](#dependencies)
        * [Local Dependencies setup](#local-dependencies-setup)
        * [Developing Python function using VS Code](#developing-python-function-using-vs-code)
        * [Publishing your function app to Azure](#publishing-your-function-app-to-azure)
<!-- TOC -->

There are 3 python functions that will be part of this function app. They all will use python 3.8.

This function app consist of the functions which are required as part of the File Upload workflow.

### ThumbnailObserver

This function is triggered on new blobs created inside of Azure Blob Storage under the temp thumbnail path.
The role of this function is to publish message to thumbnail processor topic.
Kafka output binding is used for publishing messaged to kafka topic.
Refer [here](https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-kafka-output) for more
information on this

### FileActivator

The role of this function is to update the content `status` as `UPLOADED` in `file_upload` container of the cosmos DB on
content upload

### PackageAggregator

Listens to cosmos DB change events in the `file_upload` container.
Publish message to the commit processor topic if document is ready (all contents status is `UPLOADED`) for commit phase.
Kafka output binding is used for publishing messaged to kafka topic.
Refer [here](https://learn.microsoft.com/en-us/azure/azure-functions/functions-bindings-kafka-output) for more
information on this

### File Common Package

- Add common code here.
- Currently, the common package supports the following features:
    - Logging
    - FileUpload Repository
    - Content Storage Util
    - Common util functions like api client util, constants

### Project Structure

* **requirements.txt** - Contains the list of Python packages the system installs when publishing to Azure.
* **host.json** - Contains global configuration options that affect all functions in a function app. This file does get
  published to Azure. To learn more, see [host.json](https://aka.ms/azure-functions/python/host.json).
* **.venv/** - (Optional) Contains a Python virtual environment used by local development.
* **tests/** - (Optional) Contains the test cases of your function app. For more information,
  see [Unit Testing](https://aka.ms/azure-functions/python/unit-testing).
* **.funcignore** - (Optional) Declares files that shouldn't get published to Azure. Usually, this file contains .venv/
  to ignore local Python virtual environment and tests/ to ignore test cases.
* A directory for each function app (like **ThumbnailObserver**) as well as directory for having common code named *
  *common**
* A directory hosting a function app must have **function.json**

Each function has its own code file and binding configuration file ([**function.json
**](https://aka.ms/azure-functions/python/function.json)).

### Dependencies

- **Resources**
    - `App Service Plan` = Hosts the Function App. Using Elastic Premium Plan.
    - `Storage account` = Functions relies on Azure Storage for operations such as managing triggers and logging
      function executions
    - `Application Insights` = Monitoring and logging
    - `Event Grid` = The Event Grid is used as a middle man between the Blob Storage and the trigger for
      *ThumbnailObserver* and *FileActivator*.
    - `Key Vault` = Environment variables can be safely stored in keys in the Key Vault and accessed via Application
      Settings. Key Vault is not implemented yet, but it may be implemented later. If that change has to be made it will
      not be difficult.
    - `CosmosDB` = Used for updating data in file_upload container as part of all functions in this function app.
- **Environment Variables (Application Settings)**
    - `APPINSIGHTS_INSTRUMENTATIONKEY` = The instrumentation key for the application insights that is attached to the
      function app
    - `AzureWebJobsStorage` = The connection string for the storage account in which the function app is stored. The
      Azure Functions runtime uses this storage account connection string for normal operation
    - `CosmosDbName` = The database name of the CosmosDB resource
    - `CosmosKey` = The key of the CosmosDB resource (Key value reference)
    - `CosmosAccountName` = The name of the CosmosDB resource
    - `CosmosConnectionString` = Cosmos connection string (Key value reference)
    - `FUNCTIONS_EXTENSION_VERSION` = should be set to ~4
    - `FUNCTIONS_WORKER_RUNTIME` = must be set to python
    - `BootstrapServer` = Confluent cloud cluster bootstrap server (Key value reference)
    - `ConfluentCloudUsername` = Confluent cloud API key (Key value reference)
    - `ConfluentCloudPassword` = Confluent cloud API secret (Key value reference)
    - `ThumbTopic` = ThumbProcessor Topic name
    - `CommitTopic` = CommitProcessor topic name
    - `StorageAccountUrl` = Storage account CDN url
    - `BlobContainer` = Storage account blob container name

#### Local Dependencies setup

To manually create a virtualenv on MacOS and Linux:

```
$ python -m venv .venv
```

After the init process completes and the virtualenv is created, you can use the following
step to activate your virtualenv.

```
$ source .venv/bin/activate
```

If you are a Windows platform, you would activate the virtualenv like this:

```
% .venv\Scripts\activate.bat
```

Once the virtualenv is activated, you can install the required dependencies.

```
$ pip install -r requirements.txt
```

#### Developing Python function using VS Code

If you have not already, please check out [quickstart](https://aka.ms/azure-functions/python/quickstart) to get you
started with Azure Functions developments in Python.

* To learn more about developing Azure Functions, please
  visit [Azure Functions Developer Guide](https://aka.ms/azure-functions/python/developer-guide).

* To learn more specific guidance on developing Azure Functions with Python, please
  visit [Azure Functions Developer Python Guide](https://aka.ms/azure-functions/python/python-developer-guide).

#### Publishing your function app to Azure

This folder (./src/functions/azure/UploadFunctionApp) will be deployed to a pre-configured Azure Function App with all
the necessary dependencies.

For more information on deployment options for Azure Functions, please visit
this [guide](https://docs.microsoft.com/en-us/azure/azure-functions/create-first-function-vs-code-python#publish-the-project-to-azure).
