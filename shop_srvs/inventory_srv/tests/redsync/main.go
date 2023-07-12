package main

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
	"time"
)

func main() {

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.120.172:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	rs := redsync.New(pool)

	mutexname := "my-global-mutex"
	gNums := 2
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < gNums; i++ {
		go func() {
			defer wg.Done()
			mutex := rs.NewMutex(mutexname)
			fmt.Println("开始获取锁")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}
			fmt.Println("获取锁成功")
			time.Sleep(2 * time.Second)
			fmt.Println("开始释放锁")
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("释放锁成功")
		}()
	}
	wg.Wait()

}
