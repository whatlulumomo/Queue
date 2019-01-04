# Queue
Thread safe queue in golang.

Golang 连个线程安全的锁都没有，只能自己造轮子喽

```
// Queue Transaction队列
type Queue struct {
	maxsize int // 0代表无最大容量限制 正数表示最大容量限制 负数报错
	data    list.List
	mutex   sync.Mutex // 并发锁
}

// 实现的接口函数
func (q *Queue) Init(maxsize int) error
func (q *Queue) Push(v interface{}) 
func (q *Queue) Pop() interface{}
func (q *Queue) Size() int
func (q *Queue) Get(index int) interface{} 
```
