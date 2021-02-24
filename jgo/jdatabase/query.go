package jdatabase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ijidan/jgo/jgo/jconfig"
	"github.com/ijidan/jgo/jgo/jlogger"
	"github.com/ijidan/jgo/jgo/jutils"
	"strconv"
	"strings"
)

//默认连接
const defaultConnectionName = "default"

//增删改查常量
const _select = "select"
const _update = "update"
const _delete = "delete"
const _insert = "insert"

//数据库查询
type Query struct {
	isDebug bool

	connect        *sql.DB //数据库连接
	operation      string  //操作类型
	field          string  //列
	connectionName string  //连接名称
	tablePrefix    string  //表前缀
	table          string  //表名
	primaryKey     string  //主键
	where          string
	limit          []int64
	order          []string
	group          []string
	data           map[string]interface{}
}

//是否调试
func (q *Query) SetIsDebug(isDebug bool) *Query {
	q.isDebug = isDebug
	return q
}

//查询
func (q *Query) Select(field string) *Query {
	q.operation = _select
	q.field = field
	return q
}

//插入
func (q *Query) Insert() *Query {
	q.operation = _insert
	return q
}

//更新
func (q *Query) Update() *Query {
	q.operation = _update
	return q
}

//删除
func (q *Query) Delete() *Query {
	q.operation = _delete
	return q
}

//设置数据库连接
func (q *Query) SetConnectionName(connectionName string) *Query {
	q.connectionName = connectionName
	return q
}

//设置表前缀
func (q *Query) SetPrefix(tablePrefix string) *Query {
	q.tablePrefix = tablePrefix
	return q
}

//设置主键
func (q *Query) SetPrimaryKey(primaryKey string) *Query {
	q.primaryKey = primaryKey
	return q
}

//表名
func (q *Query) From(table string) *Query {
	q.table = table
	return q
}

//条件
func (q *Query) Where(where string) *Query {
	q.where = where
	return q
}

//限制条件
func (q *Query) Limit(start int64, len int64) *Query {
	q.limit = []int64{start, len}
	return q
}

//排序
func (q *Query) Order(order []string) *Query {
	q.order = order
	return q
}

//分组
func (q *Query) Group(group []string) *Query {
	q.group = group
	return q
}

//设置数据
func (q *Query) SetData(value map[string]interface{}) *Query {
	q.data = value
	return q
}

//添加数据
func (q *Query) AddData(key string, value string) *Query {
	q.data[key] = value
	return q
}

//生成SQL
func (q *Query) GenSQL() string {
	s := ""
	switch q.operation {
	case _select:
		s = q.genSelectSQL()
		break
	case _insert:
		s = q.genInsertSQL()
		break
	case _update:
		s = q.genUpdateSQL()
		break
	case _delete:
		s = q.genDeleteSQL()
		break
	}
	return s
}

//生成查询SQL
func (q *Query) genSelectSQL() string {
	// select from table where a=1 and b=2 limit 1,2 order by xxx limit xxx
	s := fmt.Sprintf("select %s from ", q.field)
	if q.tablePrefix != "" {
		s += q.tablePrefix
	}
	s += q.table
	if len(q.where) > 0 {
		s += fmt.Sprintf(" where %s ", q.where)
	}
	if len(q.limit) > 0 {
		s += fmt.Sprintf(" limit %d,%d", q.limit[0], q.limit[1])
	}
	if len(q.order) > 0 {
		s += " order by "
		orderList := ""
		for k, v := range q.order {
			orderList += fmt.Sprintf(" %d %s,", k, v)
		}
		s += strings.TrimRight(orderList, ",")
	}
	if len(q.group) > 0 {
		s += " group by "
		groupList := ""
		for _, value := range q.group {
			groupList += fmt.Sprintf(" %s,", value)
		}
		s += strings.TrimRight(groupList, ",")
	}
	return s
}

//生成插入SQL
func (q *Query) genInsertSQL() string {
	s := "insert int64o "
	if len(q.tablePrefix) > 0 {
		s += q.tablePrefix
	}
	s += fmt.Sprintf(" %s ", q.table)
	var fileList []string
	var valueList []interface{}
	for k, v := range q.data {
		fileList = append(fileList, k)
		valueList = append(valueList, v)
	}
	val, _ := json.Marshal(valueList)
	value := string(val)
	value = strings.ReplaceAll(value, "[", "(")
	value = strings.ReplaceAll(value, "]", ")")

	insertStr := "( " + strings.Join(fileList, ",") + ") values " + value
	s += insertStr
	return s
}

//生成更新SQL
func (q *Query) genUpdateSQL() string {
	s := "update "
	if len(q.tablePrefix) > 0 {
		s += q.tablePrefix
	}
	s += fmt.Sprintf("%s set ", q.table)

	sendList := ""
	for k, v := range q.data {
		vStr := jutils.ConvertVal2String(v)
		sendList += fmt.Sprintf(" %s=%s,", k, vStr)
	}
	sendList = strings.TrimRight(sendList, ",")
	s += sendList
	if len(q.where) > 0 {
		s += fmt.Sprintf(" where %s", q.where)
	}
	return s
}

//生成删除SQL
func (q *Query) genDeleteSQL() string {
	s := "delete from  "
	if len(q.tablePrefix) > 0 {
		s += q.tablePrefix
	}
	if len(q.where) > 0 {
		s += fmt.Sprintf(" where %s", q.where)
	}
	return s
}

//获取SQL
func (q *Query) toString() string {
	return q.GenSQL()
}

//获取数据库连接
func (q *Query) GetConnection() *sql.DB {
	return q.connect
}

//获取数据库连接
func (q *Query) Connect(confName string) *sql.DB {
	cu := jconfig.ConfigUtil{}
	keyList := []string{
		"db_conf." + confName + ".host", "db_conf." + confName + ".port", "db_conf." + confName + ".dbName", "db_conf." + confName + ".username", "db_conf." + confName + ".password",
	}
	var host, portStr, dbName, username, password string
	cu.GetConfigList(keyList, &host, &portStr, &dbName, &username, &password)
	port, _ := strconv.ParseInt(portStr,10,64)
	//连接数据库
	connection, err := (&Connection{}).GetConnect(host, port, dbName, username, password)
	if err != nil {
		return nil
	}
	err = connection.Ping()
	if err != nil {
		return nil
	}
	q.connect = connection
	return connection

}

//事务
func (q *Query) Transaction(confName string, txFunc func() bool) error {
	if confName == "" {
		confName = defaultConnectionName
	}
	q.Connect(confName)
	connection := q.connect
	defer func() {
		CloseConnection(connection)
	}()
	tx, err := connection.Begin()
	if q.isDebug {
		jlogger.Notice("事务开启...")
	}
	if err != nil {
		return err
	}
	if ok := txFunc(); ok {
		if q.isDebug {
			jlogger.Notice("事务提交...")
		}
		return tx.Commit()
	} else {
		if q.isDebug {
			jlogger.Notice("事务回滚...")
		}
		return tx.Rollback()
	}
}
