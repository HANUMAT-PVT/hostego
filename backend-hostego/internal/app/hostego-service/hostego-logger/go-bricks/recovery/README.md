# Recovery Module

Often there are panics causing the runtime to crash. This package is to protect the go-routines spin-up & gin server started with appropriate recovery mechanisms

### Sample 

    package main
    
    import (
	    "fmt"
	    "github.com/gin-gonic/gin"
	    "log"
	    "net/http"
    )
    
    func GetTest(c *gin.Context) {
	    fmt.Println("serving test")
	    c.JSON(http.StatusOK, gin.H{
		    "message": "hello world",
	    })
    }
    
    func CustomRecoveryFunction() {
	    if r := recover(); r != nil {
		    err, ok := r.(error)
		    if !ok {
			    err = fmt.Errorf("%v", r)
		    }
		    fmt.Println("error is", err)
    
		    newStack := stack(3)
		    log.Printf("%s\n", string(newStack))
	    }
    }
    
    func DoPanicMain(c *gin.Context) {
	    // this does not causes server crash
	    fmt.Println("serving panic")
	    b := 0
	    fmt.Println(2 / b)
    }
    
    func DoPanicAsync(c *gin.Context) {
	    Execute(DoPanicRoutine)
    }
    
    func DoPanicAsyncCustom(c *gin.Context) {
	    CustomRecoveryExecute(DoPanicRoutine, CustomRecoveryFunction)
    }
    
    func DoPanicRoutine() {
	    fmt.Println("serving panic go routine")
	    b := 0
	    fmt.Println(2 / b)
    }
    
    func AbortWithInternalServerError(c *gin.Context, err any) {
	    c.JSON(http.StatusInternalServerError, gin.H{
		    "message": "Internal server error",
	    })
    }
    
    func main() {
	    r := gin.New()
	    r.Use(gin.Logger())
	    r.Use(CustomRecovery(AbortWithInternalServerError))
	    r.Use(gin.Recovery())
	    r.GET("/test", GetTest)
	    r.GET("/panic", DoPanicMain)
	    r.GET("/panic/go", DoPanicAsync)
	    r.GET("/panic/go/custom", DoPanicAsyncCustom)
	    r.Run()
    }
