# Milo

A utility for https://github.com/go-pg/pg that makes persisting DDD aggregates easier.

## Quick Start

The best place to start exploring Milo is by taking a look at the [example](/example). It may help to have a high level understanding of [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) before looking at the example code.

Run the example:
```bash
$ docker-compose up postgres
$ go run example/cmd/main.go
```

In the [example/cmd/example/main.go](/example/cmd/example/main.go), we see that Milo allows us to persist a `Customer` entity (including the nested addresses) to the database:

```go
store, err := milo.NewStore(db, storage.MiloEntityModelMap)
if err != nil {
	log.Fatal(err)
}

customer := &domain.Customer{
    ID: entityid.DefaultGenerator.Generate(),

    NameFirst: "John",
    NameLast:  "Smith",

    Addresses: []*domain.Address{
        {
            ID: entityid.DefaultGenerator.Generate(),

            Street: "1 City Hall Square #500",
            City:   "Boston",
            State:  "MA",
            Zip:    "02201",
        },
    },
}

err = store.Save(customer)
if err != nil {
    log.Fatal(err)
}
```

To make this work, Milo needs to be configured to understand how to map entities to storage models. This is done by passing a `milo.EntityModelMap` to `NewStore`:

```go
store, err := milo.NewStore(db, storage.MiloEntityModelMap)
```

Code from [example/storage/milo.go](/example/storage/milo.go):

```go
var MiloEntityModelMap = milo.EntityModelMap{
	reflect.TypeOf(&domain.Customer{}): reflect.TypeOf(&customer{}),
}
```

The last step is to implement `FromEntity` and `ToEntity` on the storage model. These two methods are what Milo calls to transform entities and models to and from eachother. Code from [example/storage/customer.go](/example/storage/customer.go):

```go
func (c *customer) FromEntity(e interface{}) error {
	entity := e.(*domain.Customer)

	c.ID = entity.ID.String()

	c.NameFirst = entity.NameFirst
	c.NameFirst = entity.NameLast

	for _, a := range entity.Addresses {
		c.Addresses = append(c.Addresses, &address{
			ID:         a.ID.String(),
			CustomerID: c.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return nil
}

func (c *customer) ToEntity() (interface{}, error) {
	entity := &domain.Customer{}

	entity.ID = entityid.ID(c.ID)

	entity.NameFirst = c.NameFirst
	entity.NameLast = c.NameLast

	for _, a := range c.Addresses {
		entity.Addresses = append(entity.Addresses, &domain.Address{
			ID: entityid.ID(a.ID),

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return entity, nil
}
```

### FindBy and FindOneBy

FindBy and FindOnBy (and variants) take in an additional argument of one or more `Expression`. Here is an example of `Equal` and `NotEqual`:

```go
// Find all customers with a first name of John.
customers := []*domain.Customer{}
store.FindBy(&customers, milo.Equal("name_first", "John"))

// Find the first customer that does not have the first name of John.
customer := &domain.Customer{}
store.FindOneBy(customer, milo.NotEqual("name_first", "John"))
```

Above, the first arguments to `Equal` and `NotEqual` is the column name.

You may also use the `And` and `Or` functions to create slightly more advanced expressions:

```go
// Find the first customer named John Smith.
customer := &domain.Customer{}
store.FindOneBy(customer, milo.And(milo.Equal("name_first", "John"), milo.Equal("name_last", "Smith"))

// Find all customers with the first name of John or Sally.
customers := []*domain.Customer{}
store.FindBy(&customers, milo.Or(milo.Equal("name_first", "John"), milo.Equal("name_first", "Sally"))
```

See [expression.go](/expression.go) for a full list of expression functions.

### Transactions

Milo supports database transactions through the `Transaction` method. In the example below, the last names of the customers John and Sally are updated in a single transaction:

```go
err := store.Transaction(context.Background(), func(txStore *milo.Store) error {
	var error err

	customer := &domain.Customer{}
	err = store.FindOneBy(customer, milo.Equal("name_first", "John"))
	if err != nil {
		return err
	}
	customer.NameLast = "Doe"

	customer2 := &domain.Customer{}
	err = store.FindOneBy(customer2, milo.Equal("name_first", "Sally"))
	if err != nil {
		return err
	}
	customer2.NameLast = "Doe"

	err = store.Save(customer)
	if err != nil {
		return err
	}

	err = store.Save(customer2)
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
