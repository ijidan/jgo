package jrabbitmq

import (
	"fmt"
	"github.com/ijidan/jgo/jgo/jlogger"
	"github.com/streadway/amqp"
)

//rabbit mq
type JRabbit struct {
	UserName string
	Password string
	Host     string
	Port     int64
}

//获取连接
func (jr *JRabbit) GetConnection() (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", jr.UserName, jr.Password, jr.Host, jr.Port)
	return amqp.Dial(uri)
}

//发布消息
func (jr *JRabbit) PublishMessage(exChangeName string, queueName string, messageContent string, durable bool, autoDelete bool) error {
	connection, err := jr.GetConnection()
	if err != nil {
		jlogger.Error("rabbitMQ dial error:" + err.Error())
		return err
	}
	channel, err1 := connection.Channel()
	if err1 != nil {
		jlogger.Error("rabbitMQ get channel error:" + err1.Error())
		return err1
	}
	q, err2 := channel.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	if err2 != nil {
		jlogger.Error("rabbitMQ queue declare error:" + err2.Error())
		return err2
	}
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageContent),
	}
	err3 := channel.Publish(exChangeName, q.Name, false, false, msg)
	if err3 != nil {
		jlogger.Error("rabbitMQ publish message error:" + err3.Error())
		return err3
	}
	jlogger.Info("rabbitMq published message:" + messageContent)
	return nil
}

//接受消息
func (jr *JRabbit) ReceiveMessage(queueName string, autoAck bool) error {
	connection, err := jr.GetConnection()
	if err != nil {
		jlogger.Error("rabbitMQ dial error:" + err.Error())
		return err
	}
	channel, err1 := connection.Channel()
	if err1 != nil {
		jlogger.Error("rabbitMQ get channel error:" + err1.Error())
		return err1
	}

	messageCh, err2 := channel.Consume(queueName, "", autoAck, false, false, false, nil)
	if err2 != nil {
		jlogger.Error("rabbitMQ consume error:" + err2.Error())
		return err2
	}
	for message := range messageCh {
		content := string(message.Body)
		jlogger.Info("rabbitMQ  received message:" + content)
	}
	return nil
}

const UserName = "admin"
const Password = "admin"
const Host = "172.16.1.84"
const Port = int64(5672)

//获取实例
func NewJRabbit() *JRabbit {
	jr := &JRabbit{
		UserName: UserName,
		Password: Password,
		Host:     Host,
		Port:     Port,
	}
	return jr
}
