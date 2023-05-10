# Constants
Go Module as a library for defining constant variables. Maintains constant variables that can be shared across modules.

## Example
```
const AzureHeaderKey string = "X-Azure-Ref"
const RequestIdHeaderKey string = "X-Trimble-RefId"
```

## How to Use
```
package main
import ("constants")
func PrintConstants(){
	fmt.Printf("Constant Request ID Name : %s", constants.RequestID)
}
```



