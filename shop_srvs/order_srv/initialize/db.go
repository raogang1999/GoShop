package initialize

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/order_srv/global"
	"time"
)

// 优先被动调用
func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	cnf := global.SeverConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Name)
	//dsn = "root:root@tcp(192.168.120.172:3306)/shop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//设置全局logger，打印每个sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			//ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //取消USERS这种命名方式
		},
	})
	if err != nil {
		panic(err)
	}

}

func InitRedis() {
	cnf := global.SeverConfig.RedisInfo
	global.RedisDB = goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
	})
	_, err := global.RedisDB.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	pool := goredis.NewPool(global.RedisDB) // or, pool := redigo.NewPool(...)
	global.RedisLocker = redsync.New(pool)
}
