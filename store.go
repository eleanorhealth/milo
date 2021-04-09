package milo

import (
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type Field string

type Storer interface {
	Find(entities interface{}) error
	FindBy(entities interface{}, field Field, val interface{}) error
	FindOneBy(entity interface{}, field Field, val interface{}) error
	FindByID(entity interface{}, id interface{}) error
	Save(entity interface{}) error
	Delete(entity interface{}) error
}

type Store struct {
	db             orm.DB
	entityModelMap map[reflect.Type]reflect.Type
}

var _ Storer = (*Store)(nil)

func NewStore(db orm.DB, entityModelMap EntityModelMap) *Store {
	for entityType, modelType := range entityModelMap {
		if entityType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("entity type %s must be a pointer", entityType.String()))
		}

		if modelType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("model type %s must be a pointer", modelType.String()))
		}

		modelInterfaceType := reflect.TypeOf((*Model)(nil)).Elem()

		if !modelType.Implements(modelInterfaceType) {
			panic(fmt.Sprintf("model type %s must implement %s", modelType.String(), modelInterfaceType.String()))
		}
	}

	return &Store{
		db:             db,
		entityModelMap: entityModelMap,
	}
}

func (s *Store) Find(entities interface{}) error {
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

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelsValue := reflect.New(reflect.SliceOf(modelType))
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

func (s *Store) FindBy(entities interface{}, field Field, val interface{}) error {
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

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelsValue := reflect.New(reflect.SliceOf(modelType))
	models := modelsValue.Interface()

	query := s.db.Model(models)
	query.Where(fmt.Sprintf("%s.%s = ?", query.TableModel().Table().Alias, field), val)

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

func (s *Store) FindOneBy(entity interface{}, field Field, val interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)

	relations := s.db.Model(model).TableModel().Table().Relations
	for _, relation := range relations {
		query.Relation(relation.Field.GoName)
	}

	query.Where(fmt.Sprintf("%s.%s = ?", query.TableModel().Table().Alias, field), val)

	err := query.First()
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

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	query := s.db.Model(model)
	pk := query.TableModel().Table().PKs[0]

	return s.FindOneBy(entity, Field(pk.SQLName), id)
}

func (s *Store) Save(entity interface{}) error {
	entityType := reflect.TypeOf(entity)

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "calling FromEntity on %s", modelType.Elem().String())
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

	modelType, ok := s.entityModelMap[entityType]
	if !ok {
		return fmt.Errorf("unable to find model type for entity type %s", entityType.String())
	}

	modelValue := reflect.New(modelType.Elem())
	model := modelValue.Interface().(Model)

	err := model.FromEntity(entity)
	if err != nil {
		return errors.Wrapf(err, "calling FromEntity on %s", modelType.Elem().String())
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
