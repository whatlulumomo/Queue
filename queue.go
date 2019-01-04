package main

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// TRANSACTIONCODE 消息类型
type TRANSACTIONCODE int

const (
	// DEAL 交易
	DEAL TRANSACTIONCODE = iota
	// TRANsFER 转账
	TRANsFER
)

// Transaction 交易信息
type Transaction struct {
	Type      TRANSACTIONCODE
	Coin      int
	Payer     string
	Receiver  string
	Receipt   string
	Timestamp string
	Signature [2]string // signatureR and signatureS, 对上述合并String进行ECC加密后的秘钥
	Hash      string
}

// CreateTransaction 创建一条交易记录
func CreateTransaction(Type TRANSACTIONCODE, Coin int, Payer string, Receiver string, Receipt string) *Transaction {
	trans := &Transaction{
		Type:      Type,
		Coin:      Coin,
		Payer:     Payer,
		Receiver:  Receiver,
		Receipt:   Receipt,
		Timestamp: time.Now().String(),
	}
	return trans
}

// Queue Transaction队列
type Queue struct {
	maxsize int // 0代表无最大容量限制 正数表示最大容量限制 负数报错
	data    list.List
	mutex   sync.Mutex
}

// Init 初始化
func (q *Queue) Init(maxsize int) error {
	if maxsize < 0 {
		return errors.New("Illegal Capacity, maxsize must be non-negative")
	}
	q.maxsize = maxsize
	return nil
}

// Push 推入元素
func (q *Queue) Push(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.data.PushBack(v)
	cursize := q.data.Len()

	// 队列超过最大长度则清空前1/4
	if q.maxsize != 0 && cursize >= q.maxsize { // 如果q.maxsize为0,代表无容量限制
		for i := 0; i < cursize*1/4; i++ {
			iter := q.data.Front()
			q.data.Remove(iter)
		}
	}
}

// Pop 移出元素
func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.data.Len() == 0 {
		return nil
	}

	iter := q.data.Front()
	v := iter.Value
	q.data.Remove(iter)
	return v
}

// Size 获取大小
func (q *Queue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.data.Len()
}

// Get 按索引获取元素
func (q *Queue) Get(index int) interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if index >= q.Size() {
		return nil
	}

	for iter := q.data.Front(); iter != nil; iter = iter.Next() {
		if index == 0 {
			return iter.Value
		}
		index--
	}
	return nil
}

// 测试函数
var s Queue

func main() {
	s.Init(40)
	var t sync.WaitGroup
	for i := 0; i < 100; i++ {
		t.Add(1)
		go func(i int) {
			s.Push(CreateTransaction(DEAL, i, "A", "B", "C"))
			t.Done()
		}(i)
	}

	t.Wait()
	println(s.Size())
	println(s.Pop().(*Transaction).Coin)
}
