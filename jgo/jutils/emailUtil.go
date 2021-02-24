package jutils

import "gopkg.in/gomail.v2"

//邮件工具
type EmailUtil struct {

	SMTPHost   string
	SMTPPort   int64

	FromName     string //发件人名称
	FromAddress  string //发件人邮箱
	FromPassword string
}

func (u *EmailUtil) multiSend(addressList, subject string, content string, addressNameList []string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "xiaoxu@126.com")                     //发件人
	m.SetHeader("To", "xiaoxu@163.com")           //收件人
	m.SetAddressHeader("Cc", "test@126.com", "test")     //抄送人
	m.SetHeader("Subject", "Hello!")                     //邮件标题
	m.SetBody("text/html", "使用Go测试发送邮件!")     //邮件内容
	m.Attach("E:\\IMGP0814.JPG")       //邮件附件

	d := gomail.NewDialer("smtp.126.com", 465, "xiaoxu@126.com", "1234567")
	//邮件发送服务器信息,使用授权码而非密码
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
