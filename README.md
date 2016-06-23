# Simple go cli that panics when your postgresql certificate is about to expire

This script connects to a postgresql instance on 127.0.0.1:5432, checks the certificate and panics when it's about to expire.
It's intended to be used for monitoring your postgresql certificates.

Something similar can be done using [this python script](https://github.com/thusoy/postgres-mitm/blob/master/postgres_get_server_cert.py) and `openssl` (if you like bloat :))

## Build

```shell
go build pg_check_cert.go buf.go conn.go
```

## Usage

```shell
./pg_check_cert
```

## Thanks
buf.go and conn.go are taken from [lib/pq](https://github.com/lib/pq/), see copyright notice in the files.
