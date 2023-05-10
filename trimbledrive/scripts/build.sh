#!/bin/bash

for module in ./src/api ./src/versions-worker ./src/commit-processor ./src/common-go/api_error ./src/common-go/constants ./src/common-go/requestcontext ./src/common-go/repository ./src/common-go/configuration
do
  (cd $module || exit; go mod tidy)
done

go work sync

(
  cd ./src/api || exit
  # test both debug and release mode
  go build -tags=debug -v -o ../../.bin/
  go build -v -o ../../.bin/
  cp config*.yaml ../../.bin/
)

(
  cd ./src/versions-worker || exit
  # test both debug and release mode
  go build -tags=debug -v -o ../../.bin/
  go build -v -o ../../.bin/
)

  (
    cd ./src/commit-processor || exit
    # test both debug and release mode
    go build -tags=debug -v -o ../../.bin/
    go build -v -o ../../.bin/
  )

