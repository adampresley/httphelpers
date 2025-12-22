# Requests

## AuthorizationBearer

**AuthorizationBearer** returns the token portion of a Bearer authorization
header. If the header is missing or malformed, an error is returned.

```go
token, err := requests.AuthorizationBearer(r)
```

## Get

**Get** attempts to retrieve a value from an HTTP request from all
possible sources. Here is the order of precedence.

1. POST, PUT, or PATCH form data
2. URL query parameters
3. Path parameters

This method uses generics for the type that shold be returned. It supports
the following data types: `int, int32, int64, []int, []int32, []int64, uint, uint32, uint64, []uint, []uint32, []uint64, float32, float64, []float32, []float64, string, []string, bool`.
If a value is not found, the default zero value is returned for that type.

Here is a small sample:

```go
// r is *http.Request
names := requests.Get[[]string](r, "names")
age := requests.Get[int](r, "age")
```

## StringListFromRequest

**StringListFromRequest** takes a delimited string from FORM or URL and
returns a split string slice of values.

```go
// Example URL: /?input=1,5,10
inputs := requests.StringListFromRequest(r, "input", ",")

// result is []string{"1", "5", "10"}
```

## IsHtmx

**IsHtmx** returns true if the request came from the HTMX library.

```go
// r is a *http.Request struct
isHTMX := requests.IsHtmx(r)
```

## Body

**Body** reads the body content from an http.Request and unmarshals it
into a struct. It attempts to determine the correct MIME type. Currently, this
supports JSON and XML.

```go
type Person struct {
   Name string
   Age int
}

// r is an http.Request
person, err := requests.Body[Person](r)
```
