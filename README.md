# PostgreSQL module (for Helium) provides you connection to PostgreSQL server

![Codecov](https://img.shields.io/codecov/c/github/go-helium/postgres.svg?style=flat-square)
[![Build Status](https://travis-ci.com/go-helium/redis.svg?branch=master)](https://travis-ci.com/go-helium/redis)
[![Report](https://goreportcard.com/badge/github.com/go-helium/postgres)](https://goreportcard.com/report/github.com/go-helium/postgres)
[![GitHub release](https://img.shields.io/github/release/go-helium/postgres.svg)](https://github.com/go-helium/postgres)
[![Sourcegraph](https://sourcegraph.com/github.com/go-helium/postgres/-/badge.svg)](https://sourcegraph.com/github.com/go-helium/postgres?badge)
![GitHub](https://img.shields.io/github/license/go-helium/postgres.svg?style=popout)

Module provides you connection to PostgreSQL server
- `*pg.DB` is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines

Configuration:
- yaml example
```yaml
posgres:
    hostname: string
    username: string
    password: string
    database: string
    debug: bool
    pool_size: int
    options: # optional
      host: string
      sslkey: string
      sslmode: string
      sslcert: string
      sslrootcert: string
```
- env example
```
POSTGRES_HOSTNAME=string
POSTGRES_USERNAME=string
POSTGRES_PASSWORD=string
POSTGRES_DATABASE=string
POSTGRES_DEBUG=bool
POSTGRES_POOL_SIZE=int
POSTGRES_OPTIONS_HOST=string
POSTGRES_OPTIONS_SSLKEY=string
POSTGRES_OPTIONS_SSLMODE=string
POSTGRES_OPTIONS_SSLCERT=string
POSTGRES_OPTIONS_SSLROOTCERT=string
```
