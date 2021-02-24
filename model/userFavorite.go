package model

import "github.com/ijidan/jgo/jgo/jdatabase"

//AR
type UserFavorite struct {
	jdatabase.ActiveRecord

	//字段
	Id       int64 `json:"id"`
	Uid      int64 `json:"f_uid"`
	BType    int64 `json:"f_b_type"`
	BId      int64 `json:"f_b_id"`
	CTime    int64 `json:"f_ctime"`
	UTime    int64 `json:"f_utime"`
	IsCancel int64 `json:"f_is_cancel"`
}

//获取连接名
func (u *UserFavorite) GetConnectionName() string {
	return ""
}

//获取表前缀
func (u *UserFavorite) GetTablePrefix() string {
	return "t_"
}

//获取表名
func (u *UserFavorite) GetTableName() string {
	return "user_favorite"
}

//获取主键
func (u *UserFavorite) GetPrimaryKey() string {
	return "id"
}
