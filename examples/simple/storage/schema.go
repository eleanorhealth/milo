package storage

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*customer)(nil),
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
