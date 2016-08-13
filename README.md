[![Build Status](https://travis-ci.org/deepthawtz/kv.svg?branch=master)](https://travis-ci.org/deepthawtz/kv)

kv
==

`kv` is a tiny CLI tool for getting and setting key value pairs from Consul

## Install
```
go get -v github.com/deepthawtz/kv
```

## Configure

create a `$HOME/.kv.yaml` file
```
---
consul:
  host: myconsulhost.somewhere.com
  scheme: https
```

## Examples

### get

get ALL key/value pairs at a given key prefix
```bash
kv get --prefix env/myapp/stage
```

get specific key/value pair(s) at a given key prefix
```
kv get --prefix env/myapp/stage THING_TOKEN
```

get key/value pair(s) matching pattern at a given key prefix
```
kv get --prefix env/myapp/stage THING_*
```

### set

set provided key/value pair(s) at a given prefix (*NOTE:* in Consul 0.7 this
can be performed in a single transaction)
```
kv set --prefix env/myapp/stage YO=123 THING_TOKEN=xcnvbxcmhdf COOL_FACTOR=9000
```

### del

delete provided key(s) at a given prefix (*NOTE:* in Consul 0.7 this
can be performed in a single transaction)
```
kv del --prefix env/myapp/stage YO THING_TOKEN COOL_FACTOR
```
