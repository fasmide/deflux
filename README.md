# deflux
deflux connects to deCONZ rest api and listens for sensor updates and write these to InfluxDB

## Usage

Start off by `go get`'ting deflux:

```
go get github.com/fasmide/deflux
```

deflux tries to read `$(pwd)/deflux.yml` or `/etc/deflux.yml` in that order, if both fails it will try to discover deCONZ with their webservice and output a configuration sample to stdout. 