# HttPoet

[![Build Status](https://travis-ci.org/khicago/httpoet.svg?branch=main)](https://travis-ci.org/khicago/httpoet)

A client tool for making http requests

## Usage 

### Request

```go
import (
	"github.com/khicago/httpoet"
)

func ExampleGet() {
	poet := httpoet.New("https://www.your_awsome_site.com")
	result := poet.Get("relative/path")
	body, err := result.Body()
	if err == nil {
		panic(err)
	}
	...
}
```

all http methods supported.

#### Create poet with default headers

### Header

```go
h := httpoet.H{ "key": "val" }
// or
h := httpoet.Hs{ "key": { "val", "val2"} } 
```

#### Header in action 

set or update default header to the poet

```go
h := httpoet.H{"key": "val"}
poet := httpoet.New("https://www.your_awsome_site.com").SetBaseH(h)
hs := httpoet.Hs{"key": {"v0", "v1"}}
poet.SetBaseH(hs)
```

temporary usage in request

```go
hs := httpoet.Hs{"key": {"v0", "v1"}}
result := poet.Get("relative/path", httpoet.OSetHeaders(hs))
```

#### Header combination methods

```go
var h httpoet.IHeader = httpoet.H {"k":"v"} // or httpoet.Hs{"k":{"v"}}

h.WithKV("k2", "v2")
h.WithKVAppend("k2", "v3")
h.WithH(httpoet.Hs{"k":{"v"}})
h.WithHAppend(httpoet.Hs{"k":{"v"}})
```

### Result 

```go
func ExampleResult() {
	poet := httpoet.New("https://www.your_awsome_site.com")
	result := poet.Get("relative/path")
	v := httpoet.D{}
	if err := result.ParseJson(v); err == nil {
		panic(err)
	}
	fmt.Println(v)
}
```

### Options

options can be used in the request, to modify the request contents before the msg has been sent.

options by defaults:

- `httpoet.OBackground()`
- `httpoet.OTimeout(d time.Duration)`
- `httpoet.OSetHeaders(headers IHeader)`
- `httpoet.OAddHeaders(headers IHeader)`
- `httpoet.OSetHeader(key string, value ...string)`
- `httpoet.OAppendHeader(key string, value ...string)`
- `httpoet.OAppendQuery(key string, val string)`
- `httpoet.OAppendQueryH(queries IQuery)`

to create a custom option, the httpoet.Option interface should be implemented

```go
type Option func(req *RequestBuilder) (fnDefer func())
``` 

### Query Tools

query could be used to create queries easily

```go
q := httpoet.Q{"s": "some query contents"} // or httpoet.Qs{"chars": {"a","b","c"}}
pathName := "awesome"
url := q.WriteToPth("the/%s/path", pathName)
poet.Post(url, httpoet.D{"input_val":"some example"})
```