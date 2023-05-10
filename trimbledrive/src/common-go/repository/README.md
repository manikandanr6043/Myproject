# Repository
Go Module as a library for performing database operations.

- Uses `go.mongodb.org/mongo-driver/mongo` to perform db operations.
- Exposes repository methods for consumption.

## How to Use
```
package main
import ("repository")
func GetItem(id string) any{
	item, err := repository.GetByID(id)
	if err !=nil{
		panic(err)
    }
	return item
}
```



