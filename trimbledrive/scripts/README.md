# Scripts

This directory maintains the build, deployment and maintenance scripts like a script to re process messages in DLQ.

Note, the build scripts are expected to be executed from the root repository folder.
These scripts are reused in the github workflows and during local development. E.g. VSCode tasks are setup to run these scripts to execute local unit tests and build the binaries for local execution.

## Build Scripts

The `./scripts/build.sh` compile the Go code and update/sync all the dependencies between modules. The changes in the `*.dep` and `go.work` files done by this script are expected to be committed back to the repository. If this is not done the PR workflow will fail.

## Deployment Scripts

To deploy a personal stack that has all the changes made locally run `./scripts/deploy-stack.sh <stack_suffix>`
To destroy a personal stack (or any other development stack, e.g. created by PR workflow) run `./scripts/delete-stack.sh <stack_suffix>`

If `<stack_suffix>` argument is not given then a username from the local machine is used automatically as a stack suffix.

The `./scripts/deploy-stack.sh` script uses other sub-scripts that can be used also individually:

- `./scripts/build-push-images.sh <stack_suffix>` can be used to publish latest changes in the code to container images. Remember to execute the `./scripts/deploy-images.sh <stack_suffix>` after this so the latest images are deployed to containers as new revisions
- `./scripts/deploy-bicep.sh <stack_suffix>` can be used to deploy the changes in the infrastructure

## Maintenance Scripts

The `./scripts/cosmos-db-cleanup.sh` script can be used to cleanup the CosmosDB account from the databases that were not cleaned up by normal PR cleanup workflow by some reason.
