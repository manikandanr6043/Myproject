# Repository

Go Module as a library for handling request context.

- Maintains current active request's contextual information
- Uses Gin Context to store/retrieve context.

## How to Use

``` go
package main
import ("requestcontext")
func PrintRequestContext(c *gin.Context){
    ctx := requestContext.NewRequestContext(c)
    ctx.SetUserID("user")
    ctx := requestcontext.GetContextFromGin(c)
    fmt.Printf("User ID : %s", ctx.UserID())
}
```
