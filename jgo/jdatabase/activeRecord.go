package jdatabase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ijidan/jgo/jgo/jconfig"
	"github.com/ijidan/jgo/jgo/jlogger"
	"github.com/ijidan/jgo/jgo/jutils"
	"reflect"
	"strconv"
	"strings"
)

//活动记录
type ActiveRecord struct {
	isManualClose bool
	isDebug       bool
	connect       *sql.DB
	model         interface{}
	attributes    map[string]interface{}
	oldAttributes map[string]interface{}
	columns       []string
}

//是否手动关闭连接
func (ar *ActiveRecord) SetIsManualClose(isManualClose bool) *ActiveRecord {
	ar.isManualClose = isManualClose
	return ar
}

//是否调试
func (ar *ActiveRecord) SetIsDebug(isDebug bool) *ActiveRecord {
	ar.isDebug = isDebug
	return ar
}

//设置数据库连接
func (ar *ActiveRecord) SetConnection(connect *sql.DB) *ActiveRecord {
	ar.connect = connect
	return ar
}

//传递对象
func (ar *ActiveRecord) SetModel(model interface{}) *ActiveRecord {
	ar.model = model
	return ar
}

//获取连接名称
func (ar *ActiveRecord) GetConnectionName() string {
	return ""
}

//获取表前缀
func (ar *ActiveRecord) GetTablePrefix() string {
	return ""
}

//获取表名
func (ar *ActiveRecord) GetTableName() string {
	return ""
}

//获取主键
func (ar *ActiveRecord) GetPrimaryKey() string {
	return "id"
}

//是否主键
func (ar *ActiveRecord) IsPrimaryKey(key string) bool {
	return ar.GetPrimaryKey() == key
}

//获取表结构
func (ar *ActiveRecord) GetTableSchema() interface{} {
	connection := ar.GetDbConnection()
	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}

	fullName := ar.GetFullTableName()
	sqlString := fmt.Sprintf("SHOW FULL COLUMNS FROM `%s` ", fullName)
	if ar.isDebug {
		jlogger.Info(sqlString)
	}
	rows, err := connection.Query(sqlString)
	if err != nil {
		return nil
	}
	result := ar.parseRows(rows)
	return result[0]
}

//获取属性
func (ar *ActiveRecord) GetAttributes() map[string]interface{} {
	return nil
}

//获取属性值
func (ar *ActiveRecord) GetAttribute(string) interface{} {
	return nil
}

//设置属性
func (ar *ActiveRecord) SetAttribute(string, interface{}) bool {
	return true
}

//是否有属性
func (ar *ActiveRecord) HasAttribute(string) bool {
	return true
}

//插入
func (ar *ActiveRecord) Insert() interface{} {
	return nil
}

//更新
func (ar *ActiveRecord) Update() bool {
	return true
}

//保存
func (ar *ActiveRecord) Save() bool {
	return true
}

//删除
func (ar *ActiveRecord) Delete() bool {
	return true
}

//获取数据库连接
func (ar *ActiveRecord) GetDbConnection() *sql.DB {
	if ar.connect == nil {
		//获取变量
		cu := jconfig.ConfigUtil{}
		confName := defaultConnectionName
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
		if ar.isDebug {
			jlogger.Info("db connected")
		}
		ar.connect = connection
	}
	return ar.connect
}

//通过主键查询
func (ar *ActiveRecord) FindByPk(pk int64) interface{} {
	result := ar.FindByPks([]int64{pk})
	if result == nil {
		return nil
	}
	return result[0]
}

//通过批量主键查询
func (ar *ActiveRecord) FindByPks(pks []int64) []interface{} {
	_, _, _, primaryKey := ar.getChildConfig()
	result := ar.FindByAttributes(map[string]interface{}{primaryKey: pks})
	return result
}

//通过批量属性查询
func (ar *ActiveRecord) FindByAttributes(attributes map[string]interface{}) []interface{} {
	//调用子类函数
	connectionName, tablePrefix, tableName, primaryKey := ar.getChildConfig()
	connection := ar.GetDbConnection()

	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}

	where := ar.buildWhere(attributes)
	sqlSting := (&Query{}).SetConnectionName(connectionName).SetPrefix(tablePrefix).SetPrimaryKey(primaryKey).Select("*").From(tableName).Where(where).GenSQL()
	if ar.isDebug {
		jlogger.Info(sqlSting)
	}
	rows, err := connection.Query(sqlSting)
	if err != nil {
		return nil
	}
	resultList := ar.parseRows(rows)
	return resultList
}

//通过主键更新
func (ar *ActiveRecord) UpdateByPk(attributes map[string]interface{}, pk int64) bool {
	result := ar.UpdateByPks(attributes, []int64{pk})
	return result
}

//通过批量主键更新
func (ar *ActiveRecord) UpdateByPks(attributes map[string]interface{}, pks []int64) bool {
	_, _, _, primaryKey := ar.getChildConfig()
	result := ar.UpdateByAttributes(attributes, map[string]interface{}{primaryKey: pks})
	return result
}

//通过批量属性更新
func (ar *ActiveRecord) UpdateByAttributes(attributes map[string]interface{}, condition map[string]interface{}) bool {
	//调用子类函数
	connectionName, tablePrefix, tableName, primaryKey := ar.getChildConfig()
	connection := ar.GetDbConnection()
	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}

	where := ar.buildWhere(condition)
	sqlSting := (&Query{}).SetConnectionName(connectionName).SetPrefix(tablePrefix).
		SetPrimaryKey(primaryKey).Update().From(tableName).SetData(attributes).Where(where).GenSQL()
	if ar.isDebug {
		jlogger.Info(sqlSting)
	}
	result, err := connection.Exec(sqlSting)
	if err != nil {
		return false
	}
	affectedRowsCnt, _ := result.RowsAffected()
	if ar.isDebug {
		s := fmt.Sprintf("affected rows Cnt:%d", affectedRowsCnt)
		jlogger.Notice(s)
	}
	if affectedRowsCnt > 0 {
		return true
	}
	return false
}

//通过主键删除
func (ar *ActiveRecord) DeleteByPk(pk int64) bool {
	result := ar.DeleteByPks([]int64{pk})
	return result
}

//通过批量主键删除
func (ar *ActiveRecord) DeleteByPks(pks []int64) bool {
	_, _, _, primaryKey := ar.getChildConfig()
	result := ar.DeleteByAttributes(map[string]interface{}{primaryKey: pks})
	return result
}

//通过批量属性删除
func (ar *ActiveRecord) DeleteByAttributes(attributes map[string]interface{}) bool {
	connectionName, tablePrefix, tableName, primaryKey := ar.getChildConfig()
	connection := ar.GetDbConnection()
	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}

	where := ar.buildWhere(attributes)
	sqlSting := (&Query{}).SetConnectionName(connectionName).SetPrefix(tablePrefix).SetPrimaryKey(primaryKey).Delete().From(tableName).Where(where).GenSQL()
	if ar.isDebug {
		jlogger.Info(sqlSting)
	}

	result, err := connection.Exec(sqlSting)
	if err != nil {
		return false
	}
	affectedRowsCnt, _ := result.RowsAffected()
	if ar.isDebug {
		s := fmt.Sprintf("affected rows Cnt:%d", affectedRowsCnt)
		jlogger.Notice(s)
	}
	if affectedRowsCnt > 0 {
		return true
	}
	return false
}

//原生查询
func (ar *ActiveRecord) QueryRaw(sqlString string) []interface{} {
	connection := ar.GetDbConnection()
	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}
	if ar.isDebug {
		jlogger.Info(sqlString)
	}
	rows, err := connection.Query(sqlString)
	if err != nil {
		return nil
	}
	result := ar.parseRows(rows)
	return result
}

//原生执行
func (ar *ActiveRecord) ExecRaw(sqlString string) bool {
	connection := ar.GetDbConnection()
	if ar.isManualClose != true {
		defer func() { CloseConnection(connection) }()
	}
	if ar.isDebug {
		jlogger.Info(sqlString)
	}
	lower := strings.ToLower(sqlString)
	_, err := connection.Exec(lower)
	if err != nil {
		return false
	}
	return true

}

//获取完整表名
func (ar *ActiveRecord) GetFullTableName() string {
	_, tablePrefix, tableName, _ := ar.getChildConfig()
	return tablePrefix + tableName
}

//获取改变了的属性
func (ar *ActiveRecord) GetChangedAttributes() map[string]interface{} {
	return nil
}

//获取子类相关配置
func (ar *ActiveRecord) getChildConfig() (string, string, string, string) {
	connectionName := ar.getChildFuncValue("GetConnectionName")
	tablePrefix := ar.getChildFuncValue("GetTablePrefix")
	tableName := ar.getChildFuncValue("GetTableName")
	primaryKey := ar.getChildFuncValue("GetPrimaryKey")
	return connectionName, tablePrefix, tableName, primaryKey
}

//调用子类函数
func (ar *ActiveRecord) getChildFuncValue(funcName string) string {
	v := reflect.ValueOf(ar.model)
	method := v.MethodByName(funcName)
	data := method.Call(make([]reflect.Value, 0))
	value := data[0]
	val := value.String()
	return val
}

//计算短名称
func (ar *ActiveRecord) computeShortName() string {
	t := reflect.TypeOf(ar)
	fullName := t.String()
	dotIdx := strings.Index(fullName, ".")
	shortName := fullName[dotIdx+1:]
	return shortName
}

//解析查询结果列表
func (ar *ActiveRecord) parseRows(rows *sql.Rows) []interface{} {
	//返回结果
	if rows == nil {
		return nil
	}
	columns, _ := rows.Columns()
	columnsLen := len(columns)
	values := make([][]byte, columnsLen)
	//这里表示一行填充数据
	scans := make([]interface{}, columnsLen)
	//这里scans引用values，把数据填充到[]byte里
	for k := range values {
		scans[k] = &values[k]
	}
	shortName := ar.computeShortName()
	//遍历
	var sList []interface{}
	for rows.Next() {
		//填充数据
		err := rows.Scan(scans...)
		if err != nil {
			break
		}
		sCopy := ar.model
		for colIdx, colVal := range values {
			column := columns[colIdx]
			t := reflect.TypeOf(sCopy)
			v := reflect.ValueOf(sCopy)
			tEle := t.Elem()
			vEle := v.Elem()
			//遍历字段
			for i := 0; i < tEle.NumField(); i++ {
				_v := vEle.Field(i)

				//字段
				field := tEle.Field(i)
				fieldType := field.Type
				fieldName := field.Name
				fieldTag := field.Tag.Get("json")
				if fieldName == shortName {
					continue
				}
				sColVal := string(colVal)
				if fieldTag == column {
					rColValue := ar.buildValue(sColVal, fieldType)
					_v.Set(rColValue)
				}

			}
		}
		sList = append(sList, sCopy)
	}
	return sList
}

//构造Value结构数据
func (ar *ActiveRecord) buildValue(value string, t reflect.Type) reflect.Value {
	convertedValue := jutils.ConvertValueFromString(value, t)
	return reflect.ValueOf(convertedValue)
}

//构造Where条件
func (ar *ActiveRecord) buildWhere(attributes map[string]interface{}) string {

	where := " 1= 1 "
	for k, v := range attributes {
		query := ""
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Bool:
			value := reflect.ValueOf(v).Bool()
			query = fmt.Sprintf(" and %s is %v", k, value)
			break
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			value := reflect.ValueOf(v).Int()
			query = fmt.Sprintf(" and %s=%d", k, value)
			break
		case reflect.Float32,
			reflect.Float64:
			value := reflect.ValueOf(v).Float()
			query = fmt.Sprintf(" and %s=%f", k, value)
			break
		case reflect.String:
			value := reflect.ValueOf(v).String()
			if find := strings.Contains(value, "%"); find == true {
				query = fmt.Sprintf(" and %s like '%s'", k, v)
			} else {
				query = fmt.Sprintf(" and %s='%s'", k, v)
			}
			break
		case reflect.Slice:
			val, _ := json.Marshal(v)
			value := string(val)
			value = strings.ReplaceAll(value, "[", "(")
			value = strings.ReplaceAll(value, "]", ")")
			query = fmt.Sprintf(" and %s in %s", k, value)
			break
		default:
			query = fmt.Sprintf(" and %s='%s'", k, v)
			break
		}
		where += query
	}
	return where
}
