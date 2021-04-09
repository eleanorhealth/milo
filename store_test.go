package milo

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type userEntityPtr struct {
	ID string

	NameFirst string
	NameLast  string

	Profile   *profileEntity
	Location  *locationEntity
	Addresses []*addressEntity
}

type profileEntity struct {
	ID string

	About         string
	FavoriteColor string
}

type locationEntity struct {
	ID string

	Latitude  string
	Longitude string
}

type addressEntity struct {
	ID string

	Street string
	City   string
	State  string
	Zip    string
}

type userModelPtr struct {
	tableName struct{} `pg:"users"`

	ID string `pg:"id"`

	NameFirst string `pg:"name_first"`
	NameLast  string `pg:"name_last"`

	Profile   *profileModel `pg:"rel:has-one"`
	ProfileID string        `pg:"profile_id"`

	Location *locationModel `pg:"rel:belongs-to,join_fk:user_id"`

	Addresses []*addressModel `pg:"rel:has-many,join_fk:user_id"`
}

var _ Model = (*userModelPtr)(nil)

func (u *userModelPtr) FromEntity(e interface{}) error {
	entity := e.(*userEntityPtr)

	u.ID = entity.ID
	u.NameFirst = entity.NameFirst
	u.NameLast = entity.NameLast

	if entity.Profile != nil {
		u.Profile = &profileModel{
			ID:            entity.Profile.ID,
			About:         entity.Profile.About,
			FavoriteColor: entity.Profile.FavoriteColor,
		}

		u.ProfileID = entity.Profile.ID
	}

	if entity.Location != nil {
		u.Location = &locationModel{
			ID:     entity.Location.ID,
			UserID: u.ID,

			Latitude:  entity.Location.Latitude,
			Longitude: entity.Location.Longitude,
		}
	}

	for _, a := range entity.Addresses {
		u.Addresses = append(u.Addresses, &addressModel{
			ID:     a.ID,
			UserID: u.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return nil
}

func (u *userModelPtr) ToEntity() (interface{}, error) {
	entity := &userEntityPtr{
		ID:        u.ID,
		NameFirst: u.NameFirst,
		NameLast:  u.NameLast,
	}

	if u.Profile != nil {
		entity.Profile = &profileEntity{
			ID:            u.Profile.ID,
			About:         u.Profile.About,
			FavoriteColor: u.Profile.FavoriteColor,
		}
	}

	if u.Location != nil {
		entity.Location = &locationEntity{
			ID:        u.Profile.ID,
			Latitude:  u.Location.Latitude,
			Longitude: u.Location.Longitude,
		}
	}

	for _, a := range entity.Addresses {
		entity.Addresses = append(entity.Addresses, &addressEntity{
			ID: a.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return entity, nil
}

type profileModel struct {
	tableName struct{} `pg:"profiles"`

	ID string `pg:"id"`

	About         string `pg:"about"`
	FavoriteColor string `pg:"favorite_color"`
}

type locationModel struct {
	tableName struct{} `pg:"locations"`

	ID     string `pg:"id"`
	UserID string `pg:"user_id"`

	Latitude  string `pg:"latitude"`
	Longitude string `pg:"longitude"`
}

type addressModel struct {
	tableName struct{} `pg:"addresses"`

	ID     string `pg:"id"`
	UserID string `pg:"user_id"`

	Street string `pg:"street"`
	City   string `pg:"city"`
	State  string `pg:"state"`
	Zip    string `pg:"zip"`
}

func TestStore_Pointer(t *testing.T) {
	assert := assert.New(t)

	// See docker-compose.yml
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:8200",
		User:     "postgres",
		Password: "password",
		Database: "milo",
	})
	defer db.Close()

	err := db.Ping(context.Background())
	assert.NoError(err)

	err = createSchema(db)
	assert.NoError(err)

	store := NewStore(db, map[reflect.Type]reflect.Type{
		reflect.TypeOf(&userEntityPtr{}): reflect.TypeOf(&userModelPtr{}),
	})

	user := &userEntityPtr{
		ID:        uuid.New().String(),
		NameFirst: "John",
		NameLast:  "Smith",

		Profile: &profileEntity{
			ID:            uuid.New().String(),
			About:         "Hi! I'm John.",
			FavoriteColor: "blue",
		},

		Location: &locationEntity{
			ID:        uuid.New().String(),
			Latitude:  "71.0589° W",
			Longitude: "42.3601° N",
		},

		Addresses: []*addressEntity{
			{
				ID:     uuid.New().String(),
				Street: "131 Tremont St",
				City:   "Boston",
				State:  "MA",
				Zip:    "02108",
			},
		},
	}

	err = store.Save(user)
	assert.NoError(err)

	user.NameFirst = "Jane"
	user.NameLast = "Doe"

	user.Profile.About = "Hey there! My name is Jane."
	user.Location = nil
	user.Addresses[0].Street = "101 Tremont St"

	err = store.Save(user)
	assert.NoError(err)

	user.Addresses = nil

	err = store.Save(user)
	assert.NoError(err)

	// Find
	foundUsers := []*userEntityPtr{}
	err = store.Find(&foundUsers)
	assert.NoError(err)
	assert.Len(foundUsers, 1)
	assert.Equal(user, foundUsers[0])

	// FindBy
	foundUsers = []*userEntityPtr{}
	err = store.FindBy(&foundUsers, Field("name_first"), "foo")
	assert.NoError(err)
	assert.Len(foundUsers, 0)

	// FindOneBy
	foundUser := &userEntityPtr{}
	err = store.FindOneBy(foundUser, Field("name_first"), user.NameFirst)
	assert.NoError(err)
	assert.Equal(user, foundUser)

	foundUser = &userEntityPtr{}
	err = store.FindOneBy(foundUser, Field("name_first"), "foo")
	assert.Error(err)
	assert.NotEqual(user, foundUser)

	// FindByID
	foundUser = &userEntityPtr{}
	err = store.FindByID(foundUser, user.ID)
	assert.NoError(err)
	assert.Equal(user, foundUser)

	// Delete
	err = store.Delete(user)
	assert.NoError(err)

	// Check if the user was deleted.
	foundUser = &userEntityPtr{}
	err = store.FindByID(foundUser, user.ID)
	assert.Error(err)
	assert.ErrorIs(err, pg.ErrNoRows)
}

type userEntity struct {
	ID string

	NameFirst string
	NameLast  string

	Profile   profileEntity
	Location  locationEntity
	Addresses []addressEntity
}

type userModel struct {
	tableName struct{} `pg:"users"`

	ID string `pg:"id"`

	NameFirst string `pg:"name_first"`
	NameLast  string `pg:"name_last"`

	Profile   profileModel `pg:"rel:has-one"`
	ProfileID string       `pg:"profile_id"`

	Location locationModel `pg:"rel:belongs-to,join_fk:user_id"`

	Addresses []addressModel `pg:"rel:has-many,join_fk:user_id"`
}

var _ Model = (*userModel)(nil)

func (u *userModel) FromEntity(e interface{}) error {
	entity := e.(*userEntity)

	u.ID = entity.ID
	u.NameFirst = entity.NameFirst
	u.NameLast = entity.NameLast

	u.Profile = profileModel{
		ID:            entity.Profile.ID,
		About:         entity.Profile.About,
		FavoriteColor: entity.Profile.FavoriteColor,
	}

	u.ProfileID = entity.Profile.ID

	u.Location = locationModel{
		ID:     entity.Location.ID,
		UserID: u.ID,

		Latitude:  entity.Location.Latitude,
		Longitude: entity.Location.Longitude,
	}

	for _, a := range entity.Addresses {
		u.Addresses = append(u.Addresses, addressModel{
			ID:     a.ID,
			UserID: u.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return nil
}

func (u *userModel) ToEntity() (interface{}, error) {
	entity := &userEntity{
		ID:        u.ID,
		NameFirst: u.NameFirst,
		NameLast:  u.NameLast,
	}

	entity.Profile = profileEntity{
		ID:            u.Profile.ID,
		About:         u.Profile.About,
		FavoriteColor: u.Profile.FavoriteColor,
	}

	entity.Location = locationEntity{
		ID:        u.Location.ID,
		Latitude:  u.Location.Latitude,
		Longitude: u.Location.Longitude,
	}

	for _, a := range u.Addresses {
		entity.Addresses = append(entity.Addresses, addressEntity{
			ID: a.ID,

			Street: a.Street,
			City:   a.City,
			State:  a.State,
			Zip:    a.Zip,
		})
	}

	return entity, nil
}

func TestStore_NonPointer(t *testing.T) {
	assert := assert.New(t)

	// See docker-compose.yml
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:8200",
		User:     "postgres",
		Password: "password",
		Database: "milo",
	})
	defer db.Close()

	err := db.Ping(context.Background())
	assert.NoError(err)

	err = createSchema(db)
	assert.NoError(err)

	store := NewStore(db, map[reflect.Type]reflect.Type{
		reflect.TypeOf(&userEntity{}): reflect.TypeOf(&userModel{}),
	})

	user := &userEntity{
		ID:        uuid.New().String(),
		NameFirst: "John",
		NameLast:  "Smith",

		Profile: profileEntity{
			ID:            uuid.New().String(),
			About:         "Hi! I'm John.",
			FavoriteColor: "blue",
		},

		Location: locationEntity{
			ID:        uuid.New().String(),
			Latitude:  "71.0589° W",
			Longitude: "42.3601° N",
		},

		Addresses: []addressEntity{
			{
				ID:     uuid.New().String(),
				Street: "131 Tremont St",
				City:   "Boston",
				State:  "MA",
				Zip:    "02108",
			},
		},
	}

	err = store.Save(user)
	assert.NoError(err)

	user.NameFirst = "Jane"
	user.NameLast = "Doe"

	user.Profile.About = "Hey there! My name is Jane."
	user.Location.Latitude = "71.7979° W"
	user.Location.Longitude = "21.6940° N"
	user.Addresses[0].Street = "101 Tremont St"

	err = store.Save(user)
	assert.NoError(err)

	foundUser := &userEntity{}
	err = store.FindByID(foundUser, user.ID)
	assert.NoError(err)
	assert.Equal(user, foundUser)
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*userModelPtr)(nil),
		(*profileModel)(nil),
		(*locationModel)(nil),
		(*addressModel)(nil),
	}

	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
		})
		if err != nil {
			return err
		}

		err = db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
