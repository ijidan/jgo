package jutils

import (
	"container/list"
	"sync"
)

//数组堆栈
type ArrayStack struct {
	name string //名称
	list *list.List
	lock sync.Mutex //并发锁
}

//入栈
func (stack *ArrayStack) Push(ele interface{}) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	//插入数据
	stack.list.PushFront(ele)
}

//出栈
func (stack *ArrayStack) Pop() interface{} {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	ele := stack.list.Front()
	if ele==nil{
		return nil
	}
	stack.list.Remove(ele)
	value:=ele.Value
	return value
}

//获取名称
func (stack *ArrayStack) GetName() string {
	return stack.name
}

//获取元素
func (stack *ArrayStack) GetElementList() *list.List {
	return stack.list
}

//获取长度
func (stack *ArrayStack) Size() int64 {
	return int64(stack.list.Len())
}

//实例
var arrayStack map[string]*ArrayStack
var lock sync.Mutex

//获取实例
func NewArrayStack(name string) *ArrayStack {
	lock.Lock()
	defer lock.Unlock()
	if arrayStack == nil {
		arrayStack = make(map[string]*ArrayStack)
	}
	if _, ok := arrayStack[name]; !ok {
		stack := ArrayStack{
			name: name,
			list: list.New(),
			lock: sync.Mutex{},
		}
		arrayStack[name] = &stack
	}
	return arrayStack[name]
}
