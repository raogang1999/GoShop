package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"shop/user_srv/model"
	"strings"
	"time"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(192.168.120.172:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //取消USERS这种命名方式
		},
	})
	if err != nil {
		panic(err)
	}

	//创建表
	//_ = db.AutoMigrate(&model.User{})

	//option := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("admin123", option)
	//my_password := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	////创建用户
	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("Tom%d", i),
	//		Mobile:   fmt.Sprintf("188888888%d", i),
	//		Password: my_password,
	//	}
	//	_ = db.Save(&user)
	//}

	//查询所有用户
	var users []model.User
	db.Find(&users)
	//打印
	for _, user := range users {
		fmt.Println(user.Mobile, user.NickName, user.Password)
	}

}

func password_fun() {
	fmt.Println(genMD5("123456"))

	option := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("123456", option)
	my_password := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(len(my_password))
	//解析
	passwordInfo := strings.Split(my_password, "$")

	verify := password.Verify("123456", passwordInfo[2], passwordInfo[3], option)
	fmt.Println(verify)
}

func genMD5(code string) string {
	//md5加密
	hash := md5.New()
	_, _ = io.WriteString(hash, code)

	return hex.EncodeToString(hash.Sum(nil))
}
