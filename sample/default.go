package sample

import (
	"github.com/ijidan/jgo/jgo/jdatabase"
	"github.com/ijidan/jgo/model"
)

//事务样本
func SampleTransaction() bool{
	query := jdatabase.Query{}
	err := query.Transaction("", func() bool {
		//获取数据库连接
		connection := query.GetConnection()

		//AR
		ar := jdatabase.ActiveRecord{}
		ar.SetIsDebug(true)
		ar.SetIsManualClose(true)

		user := &model.User{}
		data1 := map[string]interface{}{
			"f_sex": 0,
		}
		result1 := ar.SetConnection(connection).SetModel(user).UpdateByPk(data1, 3681)

		userFavorite := &model.UserFavorite{}
		data2 := map[string]interface{}{
			"f_utime": 0,
		}
		result2 := ar.SetConnection(connection).SetModel(userFavorite).UpdateByPk(data2, 161)
		if result1 && result2 {
			return true
		} else {
			return false
		}
	})
	if err!=nil{
		return false
	}
	return true
}
