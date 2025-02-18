// 生成 sql 子句
/* SELECT col1, col2, ...
   FROM table_name
   WHERE [ conditions ]
   GROUP BY col1
   HAVING [ conditions ]
*/
package clause

import (
	"fmt"
	"strings"
)

// generator 各个子句的生成规则
type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

// genBindVars 生成占位符 "?, ?, ?"
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// INSERT INTO $tableName ($fields)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// VALUES ($v1), ($v2), ...
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string // 用于存储占位符
	var sql strings.Builder
	var vars []interface{} // 存储所有的参数值
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{}) // eg.  []interface{}{1, "Tom", 20},
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// SELECT $fields FROM $tableName
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %v FROM %v", fields, tableName), []interface{}{}
}

// LIMIT $num
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// WHERE $desc
// eg. "age > ? AND name = ?", 18, "Tom"
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0].(string), values[1:]
	return fmt.Sprintf("WHERE %v", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
