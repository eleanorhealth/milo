# Milo

A utility package for https://github.com/go-pg/pg that makes persisting DDD aggregates easier.

## Quick Start

The best place to start exploring Milo is by taking a look at the [examples](/examples). Both the [simple](/examples/simple) and [store](/examples/store) examples borrow concepts from [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) and [Domain-Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design), but the [store](/examples/store) example shows how to call Milo from your own stores.

Run the example:
```bash
$ docker-compose up postgres
$ go run examples/simple/cmd/example/main.go
```

In the [examples/simple/cmd/example/main.go](/examples/simple/cmd/example/main.go), we see that Milo allows us to persist a `Customer` entity to the database:

```go
store, err := milo.NewStore(db, storage.MiloEntityModelMap)
if err != nil {
	log.Fatal(err)
}

customer := &domain.Customer{
	ID:        entityid.DefaultGenerator.Generate(),
	NameFirst: "Jane",
	NameLast:  "Doe",
}

err = store.Save(context.Background(), customer)
if err != nil {
	log.Fatal(err)
}
```

To make this work, Milo needs to be configured to understand how to map entities to storage models. This is done by passing a `milo.EntityModelMap` to `NewStore`:

```go
store, err := milo.NewStore(db, storage.MiloEntityModelMap)
```

Code from [examples/simple/storage/milo.go](/examples/simple/storage/milo.go):

```go
var MiloEntityModelMap = milo.EntityModelMap{
	reflect.TypeOf(&domain.Customer{}): reflect.TypeOf(&customer{}),
}
```

The last step is to implement `FromEntity` and `ToEntity` on the storage model. These two methods are what Milo calls to transform entities and models to and from eachother. Code from [examples/simple/storage/customer.go](/examples/simple/storage/customer.go):

```go
func (c *customer) FromEntity(e interface{}) error {
	entity := e.(*domain.Customer)

	c.ID = entity.ID.String()

	c.NameFirst = entity.NameFirst
	c.NameLast = entity.NameLast

	return nil
}

func (c *customer) ToEntity() (interface{}, error) {
	entity := &domain.Customer{}

	entity.ID = entityid.ID(c.ID)

	entity.NameFirst = c.NameFirst
	entity.NameLast = c.NameLast

	return entity, nil
}
```

### FindBy and FindOneBy

FindBy and FindOneBy (and variants) take in an additional argument of one or more `Expression`. Here is an example of `Equal` and `NotEqual`:

```go
// Find all customers with a first name of John.
customers := []*domain.Customer{}
store.FindBy(context.Background(), &customers, milo.Equal("name_first", "John"))

// Find the first customer that does not have the first name of John.
customer := &domain.Customer{}
store.FindOneBy(context.Background(), customer, milo.NotEqual("name_first", "John"))
```

Above, the first arguments to `Equal` and `NotEqual` are the column names you wish to apply the expression to.

You may also use the `And` and `Or` functions to create slightly more advanced expressions:

```go
// Find the first customer named John Smith.
customer := &domain.Customer{}
store.FindOneBy(context.Background(), customer, milo.And(milo.Equal("name_first", "John"), milo.Equal("name_last", "Smith"))

// Find all customers with the first name of John or Sally.
customers := []*domain.Customer{}
store.FindBy(context.Background(), &customers, milo.Or(milo.Equal("name_first", "John"), milo.Equal("name_first", "Sally"))
```

See [expression.go](/expression.go) for a full list of expression functions.

### Transactions

Milo supports database transactions through the `Transaction` method. In the example below, the last names of the customers John and Sally are updated in a single transaction:

```go
err := store.Transaction(context.Background(), func(txStore *milo.Store) error {
	var error err

	customer := &domain.Customer{}
	err = txStore.FindOneBy(context.Background(), customer, milo.Equal("name_first", "John"))
	if err != nil {
		return err
	}
	customer.NameLast = "Doe"

	customer2 := &domain.Customer{}
	err = txStore.FindOneBy(context.Background(), customer2, milo.Equal("name_first", "Sally"))
	if err != nil {
		return err
	}
	customer2.NameLast = "Doe"

	err = txStore.Save(context.Background(), customer)
	if err != nil {
		return err
	}

	err = txStore.Save(context.Background(), customer2)
	if err != nil {
		return err
	}

	return nil
})

if err != nil {
	// The transaction was rolled back.
}
```

## Running Tests

```bash
$ docker-compose up -d
$ go test -v ./...
```

## Known Limitations
* IDs must always be set on models as Milo does not do this automatically.
* Using foreign keys (defined in SQL) with `has one` relationships do not work (inserts are not ordered correctly).
* Saving `many to many` relationships is not supported (this type of relationship can be handled using hooks).
