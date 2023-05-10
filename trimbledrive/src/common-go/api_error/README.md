# Exception
Go Module as a library for handling exceptions and error codes.

## Handling Errors
Common HTTP Response Status Codes are used. The response body for all error responses are in the same format and include the error code, message and http status code.

When you send requests to and get responses from the Service API, you might encounter two types of API errors

Client errors - Client errors are indicated by a 4xx HTTP response code. Client errors indicate that the service found a problem with the client request. This can include such things as an authentication failure or missing required parameters. Fix the issue in the client application before submitting the request again.
Server errors - Server errors are indicated by a 5xx HTTP response code and need to be resolved by Trimble. 

## Error Codes

Error codes are the values that represents exact reason for failure. Every error will be associated with a specific error code.
Consider the following examples
- If a file space not found, then `SpaceNotFound` will be sent as error code.
- If a invalid token found, then `UnAuthorized` will be sent as error code.

Note that error codes are different from http status codes.

## Steps to add a new error code
To add a new error code, add an entry by creating a `ApiError` instance in `error.code.go`. 

Example
```
InvalidParam = &ApiError{ErrorMessage: "Invalid parameter", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidParam"}
```

## How to Use
```
package main
import ("exception")

func validateParam(requestParam any) (any, error){
    if requestParam == nil {
        return nil, exception.InvalidParam
    }
}
```


