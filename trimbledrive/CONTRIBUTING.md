# Contributing

1. All changes MUST be made as pull requests and associated with JIRA issues. We follow the [GitHub flow](https://docs.github.com/en/get-started/quickstart/github-flow) and "branch deploy" strategy to enforce the code in the `main` branch is always potentially releasable.

    1. Create a Jira task for the feature for tracking.
    2. Create a *feature branch* with name `feature/TD-NNN-message` where TD-NNN is the Jira item ID. The branch can be created from the `Branches` tab or locally. The branch name must begin with `feature/`. When making a change remember to update as part of the same PR:
        1. Functional tests (postman)
        2. Performance tests (k6)
        3. API specification
    3. Create a pull request, select your branch as the source and `main` as the target. Add the relevant people as reviewers.
    4. We're much more likely to approve your changes if you:
        - Add tests for new functionality. The focus should be on API level functional tests with a target is to have close to 100% coverage on those. Unit tests are used when applicable.
        - Write a [good commit message](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
    5. Fix any issues found by the reviewers in the same branch.
    6. Once the pull request has been approved, merge the branch to `main`. Remember to mark the branch for deletion after merge.

2. Bugfixes should follow the same process, but using `bugfix/TD-NNN-message` as the branch prefix.

3. We maintain a single release branch: `main`. All commits to `main` are considered as potentially releasable/deployable to production.

## Development environment setup

1. Install [GO](https://go.dev/doc/install) (version >= 1.20)
2. Install [Python](https://www.python.org/downloads/) (version >= 3.8)
3. Install [Postman](https://www.postman.com/downloads/) App. Request access to Postman team.
4. Install [k6](https://k6.io/).
5. Install [Docker](https://docs.docker.com/engine/install/)
6. Install development environment: 
   - [VSCode](https://code.visualstudio.com/download) or [GoLang IDE](https://www.jetbrains.com/go/buy/#commercial) are recommended.
       - If you are working with VSCode on Windows install WSL extension: `ms-vscode-remote.remote-wsl`
   - [GoLang IDE](https://www.jetbrains.com/go/buy/#commercial) (Another option for development in go)
   - [PyCharm](https://www.jetbrains.com/pycharm/download) (Another option for development in python)
7. Request access to Last Pass team shared folder "Shared CDE Team Accounts"
8. Request access to Azure and make sure you has access to `RG-RTMZ-TDRV-DEVELOPMENT` resource group.
    1. Install [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli): `brew update && brew install azure-cli`
    2. Install [Azure Functions Core Tools](https://learn.microsoft.com/en-us/azure/azure-functions/functions-run-local?tabs=v4%2Cmacos%2Ccsharp%2Cportal%2Cbash): `brew tap azure/functions && brew install azure-functions-core-tools@4`
9. Request access to MongoDB Atlas
    1. [Install MongoDB shell](https://www.mongodb.com/try/download/shell).
10. Request access to Confluent Kafka
    1. [Install Confluent CLI](https://docs.confluent.io/confluent-cli/current/install.html#scripted-installation) and [configure the connection](https://docs.confluent.io/confluent-cli/current/connect.html).
11. Request access for the following code quality and security tools
    1. [SonarQube](https://sonar.trimble.tools/dashboard?id=TrimbleConnect.ConnectedDataEnvironment%3ATrimbleDrive)
    2. [Whitesource](https://saas.whitesourcesoftware.com/Wss/WSS.html#!project;id=8879180)
    3. [Snyk](https://app.snyk.io/org/construction-connected-data-environment)
12. Clone the repository.
13. Follow the instructions given in the respective submodule for local development

## Understanding the repository structure

This repository follows the monorepo approach and consists of:

- src: The application implementation is split across several sub services. Implementation of each of the sub service can be found under its respective directory.
  The Purpose of the sub service and its local development setup information will be part of each sub service's README.md
- tests: tests implementation and test data (including unit, integration, functional, performance, etc)
- scripts: Purpose of this directory to maintain scripts created as part of the development process like some data migration, Reprocessing messages in a DLQ so on.
- iac: Scripts related Infrastructure as code such as terraform.

## Development workflows

1. Build the application
    - To build manually run `./scripts/build.sh`, it will produce a binary in the `.bin/` folder
    - To build using VSCode use a pre-configured `build` task
2. To execute unit tests run `./scripts/run-unit-tests.sh`
3. Running service locally is not supported. For debugging the application deploy personal DEV stack to Azure (local debugging is not supported)
    - Run `scripts/deploy-stack.sh` to deploy the personal DEV stack
    - Personal isolated DEV stack will be created using the code in your working folder including all the local modifications. The stack name will follow the pattern `tdrive-<username>`.
4. Run postman tests against personal dev stack.
    Locally installed Postman app is typically used to run tests using live collections that are kept in sync with the JSON collection in the repo and that are run using Newman during the CI/CD process.
5. Running k6 performance tests
    To do initial performance testing and e.g. compare new implementation to old one or compare different implementation options it is possible to run perf tests locally first (e.g. against personal DEV stack).

   ```bash
      k6 run smoke.js
   ```

6. Use logs to debug the code.

7. Working with the API specification

   - The API specification is maintained manually (not auto-generated) in a file named `docs/apispec.yaml`. It is recommended to use [editor.swagger.io](https://editor.swagger.io/) to edit API specification document.
   - Locally the API specification can be linted using the [Spectral](https://marketplace.visualstudio.com/items?itemName=stoplight.spectral) VSCode extension configured with the [Trimble ruleset](https://raw.githubusercontent.com/jbend/tdp-spectral-rules/main/rulesets/trimble.yaml).
   - Another useful VSCode extension for viewing the API specification locally is [Swagger Viewer](https://marketplace.visualstudio.com/items?itemName=Arjun.swagger-viewer). The preview of the API specification is available by pressing `Shift + Alt + P`.

### Updating Dependencies

Dependencies are automatically updated when you build the code according to the instructions above.
The project is configured with dependabot that will create PRs when updates are discovered for dependencies.
Typically development team aims to keep all dependencies on the latest versions and update the packages regularly (on weekly basis).

## Code analysis

A set of code analyzers are used to help with enforcing a code style, detecting common coding problems benefiting from latest language features.

The `.editorconfig` file (see [editorconfig.org](https://editorconfig.org/)) is used to configure analyzers on the editor level and allow to show errors/warnings and suggest refactoring even without compiling the code.
