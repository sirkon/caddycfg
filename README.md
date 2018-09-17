# Reflect based config unmarshaler for Caddy server

[![Build Status](https://travis-ci.org/sirkon/caddycfg.svg?branch=master)](https://travis-ci.org/sirkon/caddycfg) [comment]:[![Coverage Status](https://coveralls.io/repos/github/sirkon/caddycfg/badge.svg?branch=master)](https://coveralls.io/github/sirkon/caddycfg?branch=master)

## Installation

I hope it will be included into Caddy installation, but it is not for now. So use

```bash
go get github.com/sirkon/caddycfg
``` 

## Usage

The usage is simple, just like unmarshaling jsons with standard library tools:

```go
var c *caddy.Controller
var cfg ConfigStruct
if err := caddycfg.Unmarshal(c, &cfg); err != nil {
	return err
}
```
 

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

So, you see, this is practically the same as with JSONs.
