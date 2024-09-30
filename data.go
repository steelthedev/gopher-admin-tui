package main

import (
	"fmt"

	"gorm.io/gorm"
)

// Model fot the generic database model

type Model interface {
	GetID() interface{}
	GetTableName() string
}

// Reader interface for all read operations
type Reader interface {
	ReadAll(tableName string, filter map[string]interface{}) ([]Model, error)
	ReadByID(tableName string, id interface{}) (Model, error)
}

// Writer interface for all write operations
type Writer interface {
	Create(tableName string, model Model) error
	Update(tableName string, model Model) error
	Delete(tableName string, id interface{}) error
}

type Table interface {
	GetTableNames() ([]string, error)
	GetTableSchema(tableName string) (map[string]string, error)
}

// The general CRUD interface comprising of read and write interfaces
type CRUD interface {
	Reader
	Writer
	Table
}

// A struct that implement Gorm's CRUD interface (normal db operations with gorm)
type GormCRUD struct {
	DB *gorm.DB
}

// Implementing reading all records using Reader Interface
func (g *GormCRUD) ReadAll(tableName string, filter map[string]interface{}) ([]Model, error) {
	var models []Model
	if err := g.DB.Table(tableName).Where(filter).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (g *GormCRUD) ReadByID(tableName string, id interface{}) (Model, error) {
	var model Model
	if err := g.DB.Table(tableName).Where("id=?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (g *GormCRUD) Create(tableName string, model Model) error {

	if err := g.DB.Table(tableName).Create(model).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormCRUD) Update(tableName string, model Model) error {
	if err := g.DB.Table(tableName).Save(model).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormCRUD) Delete(tableName string, id interface{}) error {
	if err := g.DB.Table(tableName).Delete("id=?", id).Error; err != nil {
		return err
	}
	return nil
}
func (g *GormCRUD) GetTableNames() ([]string, error) {
	var tableNames []string
	err := g.DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableNames).Error
	if err != nil {
		fmt.Println("Error fetching table names:", err)
	}
	return tableNames, err
}
func (g *GormCRUD) GetTableSchema(tableName string) (map[string]string, error) {
	var columns []struct {
		ColumnName string
		DataType   string
	}
	err := g.DB.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = 'public' AND table_name = ?", tableName).Scan(&columns).Error
	if err != nil {
		return nil, err
	}

	schema := make(map[string]string)
	for _, column := range columns {
		schema[column.ColumnName] = column.DataType
	}
	return schema, nil
}

type DynamicModel struct {
	TableName string
	Schema    map[string]string
	Data      map[string]interface{}
}

func NewDynamicModel(tableName string, schema map[string]string) *DynamicModel {
	return &DynamicModel{
		TableName: tableName,
		Schema:    schema,
		Data:      make(map[string]interface{}),
	}
}
func (m *DynamicModel) GetID() interface{} {
	return m.Data["id"]
}

func (m *DynamicModel) GetTableName() string {
	return m.TableName
}

func (a *App) LoadTables() error {
	if a.Tables == nil {
		a.Tables = make(map[string]*DynamicModel)
	}

	tableNames, err := a.CRUD.GetTableNames()
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	for _, tableName := range tableNames {
		schema, err := a.CRUD.GetTableSchema(tableName)
		if err != nil {
			return fmt.Errorf("failed to get schema for table %s: %w", tableName, err)
		}

		model := NewDynamicModel(tableName, schema)
		a.Tables[tableName] = model
		// fmt.Printf("Loaded table: %s with schema: %v", tableName, schema)
	}

	return nil

}
