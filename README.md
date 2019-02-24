# PostgreSQL module (for Helium) provides you connection to PostgreSQL server

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