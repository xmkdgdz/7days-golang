package session

import (
	"geeorm/clause"
	"reflect"
)

// 实现记录增删查改相关的代码
/*
	1. 调用 clause.Set() 构造每一个子句
	2. 调用 clause.Build() 按照传入的顺序构造出最终的 SQL 语句
	3. 构造完成后，调用 Raw().Exec() 方法执行。
*/

// Insert 将（多个）已存在对象的值平铺插入数据库
func (s *Session) Insert(values ...interface{}) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	recordValues := make([]interface{}, 0)
	// 设置表字段
	table := s.Model(values[0]).RefTable()
	s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
	// 插入具体数据
	for _, value := range values {
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Find 将查询到的结果保存在对象切片中
func (s *Session) Find(values interface{}) error {
	objSlice := reflect.Indirect(reflect.ValueOf(values))
	objType := objSlice.Type().Elem() // 获取切片的单个元素的类型
	table := s.Model(reflect.New(objType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		obj := reflect.New(objType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, obj.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		objSlice.Set(reflect.Append(objSlice, obj))
	}
	return rows.Close()
}
