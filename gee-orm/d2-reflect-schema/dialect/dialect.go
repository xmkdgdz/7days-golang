// 将 Go 语言的类型映射为数据库中的类型

package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// 将 Go 语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// 返回某个表是否存在的 SQL 语句
	TableExistsSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册数据库方言
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
