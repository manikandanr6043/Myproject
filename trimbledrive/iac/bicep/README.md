# Bicep Infrastructure As Code

The Bicep allow to deploy full infrastructure needed for functional development service stack.

## Supported Stacks

Deployment of following stacks are supported:

- personal development stack (expected to be deployed by developer from the local machine using all the local changes done in the code and infrastructure)
- TMP stack deployment for feature branches (expected to be deployed as part of the GitHub workflow for Pull Requests) and automatically destroyed when PR is closed
- INTegration stack deployment (expected to be deployed automatically as part of GitHub workflow when Pull Request is merged to `main` branch)

Whether the STAGE and PROD stacks will be deployed using this bicep definition is open at the moment and will be decided later based on how infrastructure for consumer facing stacks will be managed.

## Stacks and Deployment Parameters

The deployment is controlled by number of input parameters.

The default values of parameters are set for the stacks in the DEVELOPMENT resource group (this covers personal stacks and TMP stacks). The only required parameter to be passed is the stack suffix that identifies the specific stack.
All development stacks use a shared infrastructure and have only logical isolation on the execution time by the name of the stack. E.g. following infrastructure is shared between all dev stacks:

- image registry
- Container App Environment
- Storage account
- CosmosDB account
- FrontDoor profile
- Log Analytics Workspace and Application Insights

The situation is different for the INT stack. It is deployed to a separate resource group and uses a dedicated infrastructure. So many deployment parameters need a values different from default. INT stack deployment parameters are captured in the `int.parameters.json` parameters file.

## Prerequisites

Before deploying this bicep template it is expected that following infrastructure already exists in Azure

- The Azure Resource Group is created (`RG-RTMZ-TDRV-DEVELOPMENT` for dev stacks, `RG-RTMZ-TDRV-INTEGRATION` for INT stack)
- Azure Key Vault (`KV-RTMZ-TDRV-DEVELOPMENT` for dev stacks, `KV-RTMZ-TDRV-INTEGRATION` for INT stack) is created and secrets are initialized with valid values:
  - `MONGO-URI`
  - `CONFLUENT-CLOUD-BOOTSTRAP-SERVER`
  - `CONFLUENT-CLOUD-API-KEY`
  - `CONFLUENT-CLOUD-API-SECRET`
  - `CLIENT-ID`
