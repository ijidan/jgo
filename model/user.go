package model

import "github.com/ijidan/jgo/jgo/jdatabase"

//AR
type User struct {
	jdatabase.ActiveRecord

	//字段
	Id                  int64    `json:"f_id"`
	Password            string `json:"f_password"`
	Nickname            string `json:"f_nickname"`
	Sex                 int64    `json:"f_sex"`
	Birthday            string `json:"f_birthday"`
	HeadUrl             string `json:"f_head_url"`
	MobileNation        string `json:"f_mobile_nation"`
	MobileNumMd5        string `json:"f_mobile_num_md5"`
	MobileNumCipherText string `json:"f_mobile_num_ciphertext"`
	WeiXinMd5           string `json:"f_weixin_md5"`
	WeiXinCipherText    string `json:"f_weixin_ciphertext"`
	EmailMd5            string `json:"f_email_md5"`
	EmailCipherText     string `json:"f_email_ciphertext"`
	CTime               int64    `json:"f_ctime"`
	UTime               int64   `json:"f_utime"`
	Level               int64    `json:"f_level"`
	Status              int64    `json:"f_status"`
}

//获取连接名
func (u *User) GetConnectionName() string {
	return ""
}

//获取表前缀
func (u *User) GetTablePrefix() string {
	return "t_"
}

//获取表名
func (u *User) GetTableName() string {
	return "user"
}

//获取主键
func (u *User) GetPrimaryKey() string {
	return "f_id"
}




