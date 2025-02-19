// Go 对象转换为关系型数据库中的表结构

package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represents a column of database
type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

// Schema represents a table of database
type Schema struct {
	Model      interface{} // 被映射的对象
	Name       string      // 表名
	Fields     []*Field    // 字段列表
	FieldNames []string    // 字段名(列名)
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}

// Parse 将任意的对象解析为 Schema 实例
// 例如，&User{}
func Parse(obj interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(obj)).Type()
	schema := &Schema{
		Model:    obj,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.FieldMap[field.Name] = field
		}
	}
	return schema
}

// RecordValues 将对象转换为 insert 值
func (s *Schema) RecordValues(obj interface{}) []interface{} {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, objValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
