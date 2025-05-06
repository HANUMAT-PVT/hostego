# Logger Wrapper

### Log with only common fields
``` go
commonFields: = map[string] interface {} {
        "userId": "12345",
        "ipAddress": "0.0.0.1",
}

// create new context with custom fields
newCtx: = context.WithValue(context.Background(), "commonFields", commonFields)

logger := GetNewLogger("zap", newCtx)
logger.Info("this is a info log")
```

> Key "commonFields" is a fixed common field

---

### Log with custom fields
``` go
commonFields: = map[string] interface {} {
        "userId": "12345",
        "ipAddress": "0.0.0.1",
}

// create new context with custom fields
newCtx: = context.WithValue(context.Background(), "commonFields", commonFields)

logger := GetNewLogger("zap", newCtx)
logger.Info("this is a info log with only common fields")

customFields: = map[string] interface {} {
    "key1": "val1",
}
// Add custom fields to loggerWrapper.customFields
logger.AddFields(customFields)
logger.InfofCf("this is a InfoF log with common fields & custom field: %v", "key1")
```
---
### Available loggers
  * zap: `GetNewLogger("zap", ctx)`
  * logrus: `GetNewLogger("logrus", ctx)`

### Log level
  * Set `LOG_LEVEL` env variable to `debug/info/warn`
  * default level is `info`