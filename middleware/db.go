package middleware

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var dbPingInterval = 90 * time.Second
var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

// 获取数据库的引擎
func getDBEngineDSN(database *Database) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", database.Username, database.Password, database.Host, database.Port, database.DbName, database.Charset)
	return dsn
}

func InitDB(database *Database) *gorm.DB {
	dsn := getDBEngineDSN(database)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
			//NoLowerCase:   true,  // 以大写的方式建表
		},
		// 执行任何sql时都创建并缓存预编译语句，提高后续调用速度
		PrepareStmt:          true,
		DisableAutomaticPing: false,
		// 对于写操作（创建，更新，删除），为了保证数据的完整性，GORM将它们封装在事务内运行
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
		AllowGlobalUpdate:      false,
	})
	if err != nil {
		logrus.Errorf("fail to connect database:%v\n", err)
		os.Exit(-1)
	}
	sqlDb, sqlErr := db.DB()
	if sqlErr != nil {
		logrus.Errorf("fail to connect database:%v\n", err)
		os.Exit(-1)
	}
	// 连接池空闲连接数量
	sqlDb.SetMaxIdleConns(10)
	// 打开数据库最大连接数量
	sqlDb.SetMaxOpenConns(100)
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	// 保持连接
	go KeepAliveDb(sqlDb)
	DB = db
	return db
}

func KeepAliveDb(db *sql.DB) {
	t := time.Tick(dbPingInterval)
	for {
		<-t
		err := db.Ping()
		if err != nil {
			logrus.Errorf("database ping error: %v\n", err.Error())
		}
	}
}
