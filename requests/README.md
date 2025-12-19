# Requests

## GetAuthorizationBearer

**GetAuthorizationBearer** returns the token portion of a Bearer authorization
header. If the header is missing or malformed, an error returned.

```go
token, err := httphelpers.GetAuthorizationBearer(r)
```

## GetFromRequest

**GetFromRequest** attempts to retrieve a value from an HTTP request from all
possible sources, similar to how PHP's `$_REQUEST` array works. Here is the
order of precedence.

1. FORM
2. Query
3. Multipart form data
4. Path

This method uses generics for the type that shold be returned. It supports
the following data types: `int, int32, int64, []int, []int32, []int64, uint, uint32, uint64, []uint, []uint32, []uint64, float32, float64, []float32, []float64, string, []string, bool`.
If a value is not found, the default zero value is returned for that type.

Here is a small sample:

```go
// r is *http.Request
names := httphelpers.GetFromRequest[[]string](r, "names")
age := httphelpers.GetFromRequest[int](r, "age")
```

## GetStringListFromRequest

**GetStringListFromRequest** takes a delimited string from FORM or URL and
returns a split string slice of values.

```go
// Example URL: /?input=1,5,10
inputs := httphelpers.GetStringListFromRequest(r, "input", ",")

// result is []string{"1", "5", "10"}
```

## IsHtmx

**IsHtmx** returns true if the request came from the HTMX library.

```go
// r is a *http.Request struct
isHTMX := httphelpers.IsHtmx(r)
```

## ReadBody

**ReadBody** reads the body content from an http.Request as data into
the provided destination variable. It attempts to determine the correct 
MIME type. Currently, this supports JSON and XML.

```go
type Person struct {
   Name string
   Age int
}

dest := &Person{}

// r is an http.Request
dest, err := httphelpers.ReadBody(r)
```

