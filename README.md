# PostgreSQL module (for Helium) provides you connection to PostgreSQL server

![Codecov](https://img.shields.io/codecov/c/github/go-helium/postgres.svg?style=flat-square)
![CircleCI (all branches)](https://img.shields.io/circleci/project/github/go-helium/postgres.svg?style=flat-square)
[![Report](https://goreportcard.com/badge/github.com/go-helium/postgres)](https://goreportcard.com/report/github.com/go-helium/postgres)
[![GitHub release](https://img.shields.io/github/release/go-helium/postgres.svg)](https://github.com/go-helium/postgres)
![GitHub](https://img.shields.io/github/license/go-helium/postgres.svg?style=popout)

Module provides you connection to PostgreSQL server
- `*pg.DB` is a database handle representing a pool of zero or more underlying connections. It's safe for concurrent use by multiple goroutines

Configuration:
- yaml example
```yaml
posgres:
    address: string
    username: string
    password: string
    database: string
    debug: bool
    pool_size: int
```
- env example
```
POSTGRES_ADDRESS=string
POSTGRES_USERNAME=string
POSTGRES_PASSWORD=string
POSTGRES_DATABASE=string
POSTGRES_DEBUG=bool
POSTGRES_POOL_SIZE=int
```