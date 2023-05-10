# Whitesource scan

Runs the whitesource scanner.

## Parameters

See the [action.yml](action.yml) inputs section for the parameter descriptions.

## How to use

Below you can find some examples on usage. The language specific examples only show the jobs, the basic definition can be applied to all workflows.

## Inputs for the workflow events

### Reusable Whitesource Scan

This is the recommended way to define the scan in repositories, and reuse it across other workflows, for example, to have the scan in Pull Requests, CI, CD or regular (nightly) scans.

```yaml
name: Reusable Whitesource Scan

on:
  workflow_call:
    inputs:
      branch-type:
        description: Default branch targets Trimble-Connect, feature branch targets Sandbox product.
        required: true
        default: default-branch
    secrets:
      WHITESOURCE_API_KEY:
        required: true
```

### Manual Whitesource Scan

The below example shows how to add a manual scan. Note, the inputs uses manual choice, which the user can select before running the scan, to target different types of scans.

```yaml
name: Manual Whitesource Scan

on:
  workflow_dispatch:
    inputs:
      product-name:
        type: choice
        description: The product name in WhiteSource to publish the report to.
        required: true
```

## Inputs for the scanner, and example usage

### Checkout code

If you do not need to set up and build the project with complex commands to have it ready for scanning, you can pass `checkout: true` for the inputs, to have the composite action checkout the repository. See, for example, [dotnet](#dotnet) & [nuget](#nuget) usage examples.

```yaml
with:
  checkout: true
```

#### Checkout with submodules

The composite action does not support submodules. To pull submodules you will need to do the checkout in the workflow. Use the below example, or see [checkout action's documentation](https://github.com/actions/checkout) for further information.

```yaml
  - uses: actions/checkout@v3
    with:
      submodules: "true"
      token: ${{ secrets.TC_BOT_PAT }}
```

### Java + Gradle

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Setup JDK 11
        uses: actions/setup-java@v2
        with:
          java-version: "11"
          distribution: "adopt"

      - name: Setup Gradle
        uses: gradle/gradle-build-action@v2

      - name: Whitesource Scan
        uses: ./.github/actions/whitesource-scanner
        with:
          api-key: ${{secrets.WHITESOURCE_API_KEY}}
          product-name: ${{inputs.product-name}}
```

### dotnet

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest

    steps:

      - name: Whitesource Scan
        uses: ./.github/actions/whitesource-scanner
        with:
          api-key: ${{secrets.WHITESOURCE_API_KEY}}
          product-name: ${{inputs.product-name}}
          checkout: true
          dotnet-version: 6.0.x
```

### nuget

```yaml
jobs:
  scan:
    runs-on: ubuntu-latest

    steps:

      - name: Whitesource Scan
        uses: ./.github/actions/whitesource-scanner
        with:
          api-key: ${{secrets.WHITESOURCE_API_KEY}}
          product-name: ${{inputs.product-name}}
          checkout: true
          nuget-version: 6.0.x
          dotnet-solution: Solution-name.sln
```
