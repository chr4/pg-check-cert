# pg-check-cert: Simple command line utility that can be used to monitor postgresql certificates

This script connects to a postgresql instance, checks the certificate and displays the amount of days left before it expires.
It's intended to be used for monitoring your postgresql certificates, using a monitoring tool like [Zabbix](http://www.zabbix.com/) or [Nagios](https://www.nagios.org/).

## Why openssl is not enough
I used to monitor my postgresql certificates using `openssl`. Unfortunately, the `openssl s_client` option does not support the postgresql handshake, and can therefore only look at the `.crt` file to monitor the expiration date.
As postgresql needs to be restarted after the `.crt` file was replaced, the actual file might be updated, but postgresql is still using the old certificate in-memory, until the server is restarted (as of postgresql-9.5, `reload` is not sufficient to read in the new certificate).

## Installation
Precompiled versions (linux-amd64, osx-amd64) are available on the [release page](https://github.com/chr4/pg-check-cert/releases).
Download the file, extract it and move it to e.g. `/usr/local/bin`.

## Build
```shell
go build -o pg-check-cert *.go
```

## Usage
```shell
pg-check-cert localhost:5432
```

## Thanks
- [thusoy/postgres-migm](https://github.com/thusoy/postgres-mitm/blob/master/postgres_get_server_cert.py) for the inspiration.
- `buf.go` and `conn.go` are taken from [lib/pq](https://github.com/lib/pq/), see copyright notice in the respective files.
