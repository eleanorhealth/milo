package milo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type Column string

type FieldColumnMap map[interface{}]Column

type ModelConfig struct {
	Model          reflect.Type
	FieldColumnMap FieldColumnMap
}

type EntityModelMap map[reflect.Type]ModelConfig

type Storer interface {
	Transaction(fn func(txStore *Store) error) error

	FindAll(entities interface{}) error

	FindBy(entities interface{}, exprs ...Expression) error
	FindByForUpdate(entities interface{}, exprs ...Expression) error

	FindOneBy(entity interface{}, exprs ...Expression) error
	FindOneByForUpdate(entity interface{}, exprs ...Expression) error

	FindByID(entity interface{}, id interface{}) error
	FindByIDForUpdate(entity interface{}, id interface{}) error

	Save(entity interface{}) error
	Delete(entity interface{}) error
}

type Store struct {
	db             orm.DB
	entityModelMap EntityModelMap
}

var _ Storer = (*Store)(nil)

func NewStore(db orm.DB, entityModelMap EntityModelMap) *Store {
	for entityType, modelConfig := range entityModelMap {
		if entityType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("entity type %s must be a pointer", entityType.String()))
		}

		if modelConfig.Model.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("model type %s must be a pointer", modelConfig.Model.String()))
		}

		modelInterfaceType := reflect.TypeOf((*Model)(nil)).Elem()

		if !modelConfig.Model.Implements(modelInterfaceType) {
			panic(fmt.Sprintf("model type %s must implement %s", modelConfig.Model.String(), modelInterfaceType.String()))
		}
	}

	return &Store{
		db:             db,
		entityModelMap: entityModelMap,
	}
}

func (s *Store) Transaction(fn func(txStore *Store) error) error {
	if _, ok := s.db.(*pg.Tx); ok {
		return errors.New("already in a transaction")
	}

	db := s.db.(*pg.DB)

	return db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		txStore := NewStore(tx, s.entityModelMap)
		return fn(txStore)
	})
}

func (s *Store) FindAll(entities interface{}) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("must be pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("must be slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("must be slice of pointers")
	}

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelsValue := reflect.New(reflect.SliceOf(modelConfig.Model))
	models := modelsValue.Interface()

	query := s.db.Model(models)

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err := query.Select()
	if err != nil {
		return err
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return err
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindBy(entities interface{}, exprs ...Expression) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("must be pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("must be slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("must be slice of pointers")
	}

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelsValue := reflect.New(reflect.SliceOf(modelConfig.Model))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	err := s.applyExpressionsToQuery(exprs, query, modelConfig.FieldColumnMap)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err = query.Select()
	if err != nil {
		return err
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return err
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindByForUpdate(entities interface{}, exprs ...Expression) error {
	entitiesType := reflect.TypeOf(entities)

	if entitiesType.Kind() != reflect.Ptr {
		return errors.New("must be pointer")
	}

	if entitiesType.Elem().Kind() != reflect.Slice {
		return errors.New("must be slice")
	}

	entityType := entitiesType.Elem().Elem()

	if entityType.Kind() != reflect.Ptr {
		return errors.New("must be slice of pointers")
	}

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelsValue := reflect.New(reflect.SliceOf(modelConfig.Model))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	err := s.applyExpressionsToQuery(exprs, query, modelConfig.FieldColumnMap)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(models).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	query.For(fmt.Sprintf("UPDATE OF %s", query.TableModel().Table().Alias))

	err = query.Select()
	if err != nil {
		return err
	}

	entitiesValue := reflect.ValueOf(entities).Elem()

	for i := 0; i < modelsValue.Elem().Len(); i++ {
		modelValue := modelsValue.Elem().Index(i)
		model := modelValue.Interface().(Model)

		entity, err := model.ToEntity()
		if err != nil {
			return err
		}

		entitiesValue.Set(reflect.Append(entitiesValue, reflect.ValueOf(entity)))
	}

	return nil
}

func (s *Store) FindOneBy(entity interface{}, exprs ...Expression) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	err := s.applyExpressionsToQuery(exprs, query, modelConfig.FieldColumnMap)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	err = query.First()
	if err != nil {
		return err
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return err
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) FindOneByForUpdate(entity interface{}, exprs ...Expression) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	err := s.applyExpressionsToQuery(exprs, query, modelConfig.FieldColumnMap)
	if err != nil {
		return errors.Wrap(err, "applying expressions to query")
	}

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	query.For(fmt.Sprintf("UPDATE OF %s", query.TableModel().Table().Alias))

	err = query.First()
	if err != nil {
		return err
	}

	toEntity, err := model.ToEntity()
	if err != nil {
		return err
	}

	entityValue := reflect.ValueOf(entity)
	reflect.Indirect(entityValue).Set(reflect.Indirect(reflect.ValueOf(toEntity)))

	return nil
}

func (s *Store) FindByID(entity interface{}, id interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	pk := query.TableModel().Table().PKs[0]

	return s.FindOneBy(entity, Equal(Column(pk.SQLName), id))
}

func (s *Store) FindByIDForUpdate(entity interface{}, id interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	pk := query.TableModel().Table().PKs[0]

	return s.FindOneByForUpdate(entity, Equal(Column(pk.SQLName), id))
}

func (s *Store) Save(entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "calling FromEntity on %s", modelConfig.Model.Elem().String())
	}

	exists, err := s.db.Model(model).WherePK().Exists()
	if err != nil {
		return err
	}

	// Insert

	if !exists {
		_, err := s.db.Model(model).Insert()
		if err != nil {
			return err
		}

		err = s.insertRelated(model)
		if err != nil {
			return err
		}

		return nil
	}

	// Update

	_, err = s.db.Model(model).WherePK().Update()
	if err != nil {
		return err
	}

	err = s.deleteRelated(model)
	if err != nil {
		return err
	}

	err = s.insertRelated(model)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelConfig, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model config for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelConfig.Model.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "calling FromEntity on %s", modelConfig.Model.Elem().String())
	}

	_, err = s.db.Model(model).WherePK().Delete()
	if err != nil {
		return err
	}

	err = s.deleteRelated(model)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) insertRelated(model Model) error {
	modelValue := reflect.ValueOf(model)

	relations := s.db.Model(model).TableModel().Table().Relations
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

		_, err := s.db.Model(relatedModelField.Interface()).Insert()
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Store) deleteRelated(model Model) error {
	modelValue := reflect.ValueOf(model)

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {

		relatedModelFieldValue := reflect.New(relation.Field.Type)
		if relatedModelFieldValue.Kind() != reflect.Ptr {
			relatedModelFieldValue = relatedModelFieldValue.Addr()
		}

		deleteQuery := s.db.Model(relatedModelFieldValue.Interface())

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

func (s *Store) applyExpressionsToQuery(exprs []Expression, query *orm.Query, fieldColumnMap FieldColumnMap) error {
	for _, e := range exprs {
		if len(e.exprs) > 0 {

			query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				err := s.applyExpressionsToQuery(e.exprs, q, fieldColumnMap)
				return q, err
			})

		} else {

			var column Column
			var ok bool
			if column, ok = e.Field().(Column); !ok {
				column, ok = fieldColumnMap[e.Field()]
				if !ok {
					return fmt.Errorf("unable to find column for field %s", e.Field())
				}
			}

			switch e.t {
			case expressionTypeAnd:
				query.Where(fmt.Sprintf("%s.%s %s ?", query.TableModel().Table().Alias, column, e.op), e.Value())

			case expressionTypeOr:
				query.WhereOr(fmt.Sprintf("%s.%s %s ?", query.TableModel().Table().Alias, column, e.op), e.Value())

			default:
				return fmt.Errorf("unknown ExpressionType: %s", e)
			}

		}
	}

	return nil
}
