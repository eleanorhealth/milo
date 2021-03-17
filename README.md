# Milo

A utility for https://github.com/go-pg/pg that makes persisting DDD aggregates easier.

## Running Tests and Example

```bash
$ docker-compose up -d
```

Tests:
```bash
$ go test -v ./...
```

Example:
```bash
$ go run example/cmd/example/main.go
```

## Known Limitations
* IDs must be set on related models as Milo does not do this automatically.
* Using foreign keys (defined in SQL) with `has one` relationships does not work (inserts are not ordered correctly).
* Finding by ID with composite primary keys does not work.
* Saving `many to many` relationships are not supported (this type of relationship can be handled using hooks).
