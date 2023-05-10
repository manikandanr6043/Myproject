#!/bin/bash

retval=0

for module in ./src/api ./src/versions-worker ./src/commit-processor ./src/common-go/api_error ./src/common-go/constants ./src/common-go/requestcontext ./src/common-go/repository ./src/common-go/configuration
do
  go test -v -race -vet=off $module/... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
  if [ "${PIPESTATUS[0]}" != 0 ]
  then
    retval=1;
  fi
done

exit $retval
