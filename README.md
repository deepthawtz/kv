kv
==

`kv` is a tiny CLI tool for getting and setting key value pairs from Consul

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

set provided key/value pair(s) at a given prefix
```
kv set --prefix env/myapp/stage YO=123 THING_TOKEN=xcnvbxcmhdf COOL_FACTOR=9000
```
