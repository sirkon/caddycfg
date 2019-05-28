# Reflect based config unmarshaler for Caddy server

It is already a part of upcoming Caddy 2: https://github.com/mholt/caddy/issues/2295

[![Build Status](https://travis-ci.org/sirkon/caddycfg.svg?branch=master)](https://travis-ci.org/sirkon/caddycfg) [![codecov](https://codecov.io/gh/sirkon/caddycfg/branch/master/graph/badge.svg)](https://codecov.io/gh/sirkon/caddycfg) [![GoDoc](https://godoc.org/github.com/sirkon/caddycfg?status.svg)](https://godoc.org/github.com/sirkon/caddycfg) 

## Installation

```bash
go get github.com/sirkon/caddycfg
``` 

## Usage

The usage is simple, very much like unmarshaling jsons with standard library tools:

```go
var c *caddy.Controller
var cfg ConfigStruct
```

One use
```go
if err := caddycfg.Unmarshal(c, &cfg); err != nil {
    return err
}
```

Got a position of plugin name in config file
```go
head, err := caddycfg.UnmarshalHeadInfo(c, &cfg); err != nil {
    return err
}
```
where head is `Token`
```go
type Token struct {
    File  string
    Lin   int
    Col   int
    Value string
}
```

## Config types

> Please remember, this library grows from our need to reuse our existing pieces at my job, where we use JSON configs for our microservices. That's why it needs `json` tag for any field. Thats is not bad. It also supports `json.Unmarshaler` to the certain extent — value to be decoded must come in our piece, i.e. single `c.Next()` or `c.NextArg()` footprint which is to be returned by `c.Val()`


##### Example 1

Let we have plugin called _plugin_.
 
```
plugin {
    key1 value1
    key2 value2
}
```

This can be parsed with the following structure: 

```go
type pluginConfig struct {
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
} 
```


##### Example 2


A bit harder example

```
plugin {
    key1 value11 value12
    key2 value2
}
```

This can be parsed with

```go
type pluginConfig struct {
	Key1 []string `json:"key1"`
	Key2 string   `json:"key2"`
}
```


##### Example 3

Arguments appears before block

```
plugin arg1 arg2 {
    key1 value1
    key2 value2
}
```

This can be parsed with

```go
type pluginConfig struct {
	caddycfg.Args
	
	Key1 string `json:"key1"`
	Key2 string `json:"key2"`
}
```

`pluginConfig` has a method `Arguments` from embedded type `caddycfg.Args` to get these positional parameters

There is also alternative approach, which needs a bit more work yet provides better control over content: 

there's an interface 

```go
type ArgumentsConsumer interface {
    ConsumeArguments(head Token, args []Token) error
    Arguments() []string
}
```

if a structure implements it then its (of the structure) methods will be used to deal with positional parameters:
method `ConsumeArguments` will be called to consume positional arguments at once. Note though, block is not mandatory 
for types implementing `ArgumentsConsumer`


##### Example 4

Internal blocks

```
plugin arg {
    key1 subarg {
        key value
    }
    key2 value
}
```

Use

```go
type pluginConfig struct {
    caddycfg.Args
    
    Key1 subConfig `json:"key1"`
    Key2 int       `json:"key2"`
}

type subConfig struct {
    caddycfg.Args
    
    Key string `json:"key"`
}
```

##### Example 5

Parse both

```
plugin {
    a 1
    b 2
    c 3
}
```

and

```
plugin {
   someStrangeKeyName itsValue
}
```

with one type? It is easy! Just use

```go
var target map[string]string
```

for this.

##### Example 6

Don't bother with type for simple things like

```
plugin a b c
```
?

Use

```go
var target []string
```

then.

Remember, `[]string` can also be used to unmarshal this config:


```
plugin {
    a
    b
    c
}
```

## Structure validation

Data check may be needed at times. If destination type implements 

```go
type Validator interface {
	Err(head caddycfg.Token) error
}
```

function `Err(headToken)` will be called after unmarshaling 
