# Trimble Drive

![latest deployment to INT](https://github.com/trimble-oss/trimble-drive/actions/workflows/CD.yml/badge.svg)

This repository contains the Trimble Drive Service.
The service provides capabilities for managing files (content blobs with metadata) in the context of Connected Data Environment serving the Industry Clouds needs

## References

- [Documentation on Google Drive](https://drive.google.com/drive/folders/1Qy9K8fck74nqA0KZoYvvwnZ03BDEh33G)
- [JIRA board](https://jira.trimble.tools/secure/RapidBoard.jspa?rapidView=14593)
- [Trimble Drive Dev Google Group](https://groups.google.com/a/trimble.com/g/trimble-drive-dev-ug) - used for alarms and notifications on dev stacks

Please monitor the quality fo the code by checking the following reports:

- [SonarQube](https://sonar.trimble.tools/dashboard?id=TrimbleConnect.ConnectedDataEnvironment%3ATrimbleDrive)
- [Whitesource](https://saas.whitesourcesoftware.com/Wss/WSS.html#!project;id=8879180)
- [Snyk](https://app.snyk.io/org/construction-connected-data-environment)

## Service Instances

The service is deployed to following customer facing environments

- `STAGE`: TODO url regions
- `PROD`: TODO url regions

Permanent development environment for continuous integration and deployment (latest and greatest potentially releasable version):

- `INTegration`: https://tdrive-api-int-e5bfh0cfeygqebgh.z01.azurefd.net

In addition number of temporary environments can be active at any moment of time:

- personal development and/or experimental stacks
- `TMP` stacks that are very short-lived used during CI workflow
