package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/user_srv/global"
	"time"
)

// 优先被制动调用
func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	cnf := global.SeverConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Name)
	//dsn = "root:root@tcp(192.168.120.172:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//设置全局logger，打印每个sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
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
