# Responses

## Html

**Html** writes content to the response writer with a `text/html` header.

```go
content := `<p>This is a test</p>`
responses.Html(w, http.StatusOK, content)
```

## HtmlOK

**HtmlOK** writes content to the response writer with a `text/html` header, an 
a 200 OK status.

```go
content := `<p>This is a test</p>`
responses.HtmlOK(w, content)
```

## IsSuccessRange

**IsSuccessRange** returns true if the status code falls within 200-299 range.

```go
status := someHttpCall()
isSuccess := responses.IsSuccessRange(status)
```

### Json

**Json** converts any arbitrary structure to JSON and writes it to an HTTP writer. If there 
is an error marshalling the value, it writes a 500 status code with a generic 
error message.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

responses.Json(w, http.StatusOK, output)
```

### JsonOK

**JsonOK** returns a _200 OK_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

responses.JsonOK(w, output)
```

### JsonBadRequest

**JsonBadRequest** returns a _400 Bad Request_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

responses.JsonBadRequest(w, output)
```

### JsonInternalServerError

**JsonInternalServerError** returns a _500 Internal Server Error_ status with an arbitrary structure converted to JSON.

```go
output := SomeType{
  Key1: "Adam",
  Key2: 10,
}

responses.JsonInternalServerError(w, output)
```

### JsonErrorMessage

**JsonErrorMessage** returns a specified status code along with a generic structure containing an error message
in JSON.

```go
responses.JsonErrorMessage(w, http.StatusInternalServerError, "something went wrong")
// The result written is {"message": "something went wrong"}
```

### JsonUnauthorized

**JsonUnauthorized** returns a status code _401 Unauthorized_ along with an arbitrary structure converted to JSON.
in JSON.

```go
value := map[string]any{"message": "not authorized"}
httphelpers.JsonErrorMessage(w, http.StatusInternalServerError, value)
// The result written is {"message": "not authorized"}
```

