package milo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type EntityModelMap map[reflect.Type]reflect.Type

type Storer interface {
	Transaction(ctx context.Context, fn func(txStore Storer) error) error

	FindAll(ctx context.Context, entities interface{}) error

	FindBy(ctx context.Context, entities interface{}, exprs ...Expression) error
	FindByForUpdate(ctx context.Context, entities interface{}, skipLocked bool, exprs ...Expression) error

	FindOneBy(ctx context.Context, entity interface{}, exprs ...Expression) error
	FindOneByForUpdate(ctx context.Context, entity interface{}, skipLocked bool, exprs ...Expression) error

	FindByID(ctx context.Context, entity interface{}, id interface{}) error
	FindByIDForUpdate(ctx context.Context, entity interface{}, id interface{}, skipLocked bool) error

	Save(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, entity interface{}) error
}

type Store struct {
	db             orm.DB
	entityModelMap EntityModelMap
}

var _ Storer = (*Store)(nil)

func NewStore(db orm.DB, entityModelMap EntityModelMap) (*Store, error) {
	for entityType, modelType := range entityModelMap {
		if entityType.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("entity type %s must be a pointer", entityType.String())
		}

		if modelType.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("model type %s must be a pointer", modelType.String())
		}

		modelInterfaceType := reflect.TypeOf((*Model)(nil)).Elem()

		if !modelType.Implements(modelInterfaceType) {
			return nil, fmt.Errorf("model type %s must implement %s", modelType.String(), modelInterfaceType.String())
		}
	}

	return &Store{
		db:             db,
		entityModelMap: entityModelMap,
	}, nil
}

func (s *Store) inTransaction() bool {
	_, ok := s.db.(*pg.Tx)
	return ok
}

func applyExpressionsToQuery(exprs []Expression, query *orm.Query) error {
	for _, e := range exprs {
		if len(e.exprs) > 0 {

			query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				err := applyExpressionsToQuery(e.exprs, q)
				return q, err
			})

		} else {

			var condition string
			var params []interface{}

			if e.op == OpIsNull || e.op == OpIsNotNull {
				condition = fmt.Sprintf("%s.%s %s", query.TableModel().Table().Alias, e.column, e.op)
			} else {
				condition = fmt.Sprintf("%s.%s %s ?", query.TableModel().Table().Alias, e.column, e.op)
				params = append(params, e.Value())
			}

			switch e.t {
			case expressionTypeAnd:
				query.Where(condition, params...)

			case expressionTypeOr:
				query.WhereOr(condition, params...)

			default:
				return fmt.Errorf("unknown expressionType: %s", reflect.TypeOf(e.t).String())
			}

		}
	}

	return nil
}

// Transaction runs function fn in a transaction. If fn returns an error, the transaction is rolled back. Otherwise, the transaction is committed.
func (s *Store) Transaction(ctx context.Context, fn func(txStore Storer) error) error {
	if s.inTransaction() {
		return errors.New("already in a transaction")
	}

	return s.db.(*pg.DB).RunInTransaction(ctx, func(tx *pg.Tx) error {
		txStore, err := NewStore(tx, s.entityModelMap)
		if err != nil {
			return errors.Wrap(err, "creating a new store for the transaction")
		}

		return fn(txStore)
	})
}

func (s *Store) FindAll(ctx context.Context, entities interface{}) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("entities must be a pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("entities must be a slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("entities must be a slice of pointers")
	}

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return errors.New(fmt.Sprintf("unable to find model type for entity type %s", entityType.String()))
	}

	modelsValue := reflect.New(reflect.SliceOf(modelType))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	query.Context(ctx)

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err := query.Select()
	if err != nil {
		return errors.Wrap(err, "selecting the model")
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return errors.Wrap(err, "converting model to entity")
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindBy(ctx context.Context, entities interface{}, exprs ...Expression) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("entities must be a pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("entities must be a slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("entities must be a slice of pointers")
	}

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return errors.New(fmt.Sprintf("unable to find model type for entity type %s", entityType.String()))
	}

	modelsValue := reflect.New(reflect.SliceOf(modelType))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	query.Context(ctx)
	err := applyExpressionsToQuery(exprs, query)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err = query.Select()
	if err != nil {
		return errors.Wrap(err, "selecting the model")
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return errors.Wrap(err, "converting model to entity")
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindByForUpdate(ctx context.Context, entities interface{}, skipLocked bool, exprs ...Expression) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("entities must be a pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("entities must be a slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("entities must be a slice of pointers")
	}

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return errors.New(fmt.Sprintf("unable to find model type for entity type %s", entityType.String()))
	}

	modelsValue := reflect.New(reflect.SliceOf(modelType))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	query.Context(ctx)
	err := applyExpressionsToQuery(exprs, query)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	var skipLockedSQL string
	if skipLocked {
		skipLockedSQL = " SKIP LOCKED"
	}

	query.For(fmt.Sprintf("UPDATE OF %s%s", query.TableModel().Table().Alias, skipLockedSQL))

	err = query.Select()
	if err != nil {
		return errors.Wrap(err, "selecting the model")
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return errors.Wrap(err, "converting model to entity")
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindOneBy(ctx context.Context, entity interface{}, exprs ...Expression) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	query.Context(ctx)
	err := applyExpressionsToQuery(exprs, query)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err = query.First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return ErrNotFound
		}

		return errors.Wrap(err, "selecting first row")
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return errors.Wrap(err, "converting model to entity")
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) FindOneByForUpdate(ctx context.Context, entity interface{}, skipLocked bool, exprs ...Expression) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	query.Context(ctx)
	err := applyExpressionsToQuery(exprs, query)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	var skipLockedSQL string
	if skipLocked {
		skipLockedSQL = " SKIP LOCKED"
	}

	query.For(fmt.Sprintf("UPDATE OF %s%s", query.TableModel().Table().Alias, skipLockedSQL))

	err = query.First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return ErrNotFound
		}

		return errors.Wrap(err, "selecting first row")
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return errors.Wrap(err, "converting model to entity")
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) FindByID(ctx context.Context, entity interface{}, id interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	query.Context(ctx)

	for _, pk := range query.TableModel().Table().PKs {
		query.Where(fmt.Sprintf("%s.%s = ?", query.TableModel().Table().Alias, pk.SQLName), id)
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err := query.First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return ErrNotFound
		}

		return errors.Wrap(err, "selecting first row")
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return errors.Wrap(err, "converting model to entity")
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) FindByIDForUpdate(ctx context.Context, entity interface{}, id interface{}, skipLocked bool) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	query.Context(ctx)

	for _, pk := range query.TableModel().Table().PKs {
		query.Where(fmt.Sprintf("%s.%s = ?", query.TableModel().Table().Alias, pk.SQLName), id)
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	var skipLockedSQL string
	if skipLocked {
		skipLockedSQL = " SKIP LOCKED"
	}

	query.For(fmt.Sprintf("UPDATE OF %s%s", query.TableModel().Table().Alias, skipLockedSQL))

	err := query.First()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return ErrNotFound
		}

		return errors.Wrap(err, "selecting first row")
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return errors.Wrap(err, "converting model to entity")
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) Save(ctx context.Context, entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "converting entity to model")
	}

	var tx *pg.Tx

	if s.inTransaction() {
		tx = s.db.(*pg.Tx)
	} else {
		tx, err = s.db.(*pg.DB).Begin()
		if err != nil {
			return errors.Wrap(err, "beginning transaction")
		}

		defer tx.Rollback()
	}

	if model, ok := model.(Hook); ok {
		store, err := NewStore(tx, s.entityModelMap)
		if err != nil {
			return errors.Wrap(err, "creating new store for before save hook")
		}

		err = model.BeforeSave(ctx, store, entity)
		if err != nil {
			return errors.Wrap(err, "calling before save hook")
		}
	}

	exists, err := tx.Model(model).WherePK().Exists()
	if err != nil {
		return errors.Wrap(err, "exists")
	}

	// Insert

	if !exists {
		_, err := tx.Model(model).Insert()
		if err != nil {
			return errors.Wrap(err, "inserting model")
		}

		err = insertRelated(tx, model)
		if err != nil {
			return errors.Wrap(err, "inserting related models (insert)")
		}

		if !s.inTransaction() {
			err = tx.Commit()
			if err != nil {
				return errors.Wrap(err, "committing transaction")
			}
		}

		return nil
	}

	// Update

	_, err = tx.Model(model).WherePK().Update()
	if err != nil {
		return errors.Wrap(err, "updating model")
	}

	err = deleteRelated(tx, model)
	if err != nil {
		return errors.Wrap(err, "deleting related models (update)")
	}

	err = insertRelated(tx, model)
	if err != nil {
		return errors.Wrap(err, "inserting related models (update)")
	}

	if !s.inTransaction() {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction")
		}
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "converting entity to model")
	}

	var tx *pg.Tx

	if s.inTransaction() {
		tx = s.db.(*pg.Tx)
	} else {
		tx, err = s.db.(*pg.DB).Begin()
		if err != nil {
			return errors.Wrap(err, "beginning transaction")
		}

		defer tx.Rollback()
	}

	if model, ok := model.(Hook); ok {
		store, err := NewStore(tx, s.entityModelMap)
		if err != nil {
			return errors.Wrap(err, "creating new store for before delete hook")
		}

		err = model.BeforeDelete(ctx, store, entity)
		if err != nil {
			return errors.Wrap(err, "calling before delete hook")
		}
	}

	_, err = tx.Model(model).WherePK().Delete()
	if err != nil {
		return err
	}

	err = deleteRelated(tx, model)
	if err != nil {
		return errors.Wrap(err, "deleting related models (delete)")
	}

	if !s.inTransaction() {
		err = tx.Commit()
		if err != nil {
			return errors.Wrap(err, "committing transaction")
		}
	}

	return nil
}

func insertRelated(db orm.DB, model Model) error {
	modelValue := reflect.ValueOf(model)

	relations := db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {

		// Many to many relationships are not supported.
		if relation.Type == orm.Many2ManyRelation {
			continue
		}

		relatedModelField := modelValue.Elem().FieldByName(relation.Field.GoName)

		// Don't insert nil values.
		if (relatedModelField.Kind() == reflect.Ptr || relatedModelField.Kind() == reflect.Slice) && relatedModelField.IsNil() {
			continue
		}

		// If relatedModelField isn't a pointer, get its address as that's what go-pg expects.
		if relatedModelField.Kind() != reflect.Ptr {
			relatedModelField = relatedModelField.Addr()
		}

		_, err := db.Model(relatedModelField.Interface()).Insert()
		if err != nil {
			return err
		}

	}

	return nil
}

func deleteRelated(db orm.DB, model Model) error {
	modelValue := reflect.ValueOf(model)

	relations := db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {

		relatedModelFieldValue := reflect.New(relation.Field.Type)
		if relatedModelFieldValue.Kind() != reflect.Ptr {
			relatedModelFieldValue = relatedModelFieldValue.Addr()
		}

		deleteQuery := db.Model(relatedModelFieldValue.Interface())

		for i, fk := range relation.JoinFKs {
			baseFK := relation.BaseFKs[i]
			modelFKValue := modelValue.Elem().FieldByName(baseFK.GoName)

			deleteQuery.Where(fmt.Sprintf("%s = ?", fk.SQLName), modelFKValue.Interface())
		}

		_, err := deleteQuery.Delete()
		if err != nil {
			return err
		}

	}

	return nil
}
