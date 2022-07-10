package ghost

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

type BaseDBModel struct {
	Id        int       `gorm:"primaryKey"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"autoUpdateTime"`
}

type dbModel interface {
	TableName() string
}

var alias2db = make(map[string]*gorm.DB)
var alias2dbModels = make(map[string][]dbModel)

type dbConfig struct {
	engine    string
	dsn       string
	host      string
	port      string
	user      string
	pwd       string
	dbname    string
	debugMode bool

	maxConns       int // 最大连接数
	maxIdleConns   int
	maxIdleTimeout int
}

func (this *dbConfig) GetDsn() string {
	if this.dsn != "" {
		return this.dsn
	}
	switch this.engine {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			this.user,
			this.pwd,
			this.host,
			this.port,
			this.dbname,
		)
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			this.user,
			this.pwd,
			this.host,
			this.port,
			this.dbname,
		)
	default:
		return this.dsn
	}
}

func (this *dbConfig) GetDbName() string {
	if this.dbname != "" {
		return this.dbname
	}
	return strings.Split(strings.Split(this.dsn, "?")[0], "/")[1]
}

func (this *dbConfig) IsDebugMode() bool {
	return this.debugMode
}

func NewDbConfigFromDSN(engine, dsn string) *dbConfig {
	instance := new(dbConfig)
	instance.engine = engine
	instance.dsn = dsn
	return instance
}

func NewDbConfig(engine, host, port, user, pwd, dbname string, debugMode bool, args ...int) *dbConfig {
	instance := new(dbConfig)
	instance.engine = engine
	instance.host = host
	instance.port = port
	instance.user = user
	instance.pwd = pwd
	instance.dbname = dbname
	instance.debugMode = debugMode
	l := len(args)
	if l >= 1 {
		instance.maxConns = args[0]
	} else {
		instance.maxConns = 10
	}
	if l >= 2 {
		instance.maxIdleConns = args[1]
	}
	if l >= 3 {
		instance.maxIdleTimeout = args[2]
	}
	return instance
}

func GetDB(args ...string) *gorm.DB {
	nsp := "default"
	switch len(args) {
	case 1:
		nsp = args[0]
	}
	if d, ok := alias2db[nsp]; ok && d != nil {
		return d
	}
	return nil
}

func GetDBFromCtx(ctx context.Context, args ...string) *gorm.DB {
	db := GetDB(args...)
	if idb := ctx.Value("db"); idb != nil {
		db = idb.(*gorm.DB)
	}
	return db
}

func ConnectDB(dbconfig *dbConfig, args ...string) *gorm.DB {
	alias := "default"
	switch len(args) {
	case 1:
		alias = args[0]
	}
	logLevel := logger.Silent
	if dbconfig.IsDebugMode() {
		logLevel = logger.Info
	}
	var dial gorm.Dialector
	dsn := dbconfig.GetDsn()
	log.Printf(dsn)
	switch dbconfig.engine {
	case "mysql":
		dial = mysql.Open(dsn)
	case "postgres":
		dial = postgres.Open(dsn)
	default:
		log.Panic("unsupport db driver", dbconfig.engine)
	}
	gdb, err := gorm.Open(dial, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel: logLevel, // Log level
				Colorful: true,
			},
		),
	})
	if err != nil {
		log.Panic(err)
	}
	db, err := gdb.DB()
	if err != nil {
		log.Panic(err)
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(dbconfig.maxIdleTimeout))
	db.SetMaxIdleConns(dbconfig.maxIdleConns)
	db.SetMaxOpenConns(dbconfig.maxConns)
	alias2db[alias] = gdb
	log.Printf("connecting %s db: %s ...", dbconfig.engine, dbconfig.GetDbName())
	return gdb
}

func RegisterDBModel(dm dbModel, args ...string) {
	alias := "default"
	switch len(args) {
	case 1:
		alias = args[0]
	}
	if _, ok := alias2dbModels[alias]; ok {
		alias2dbModels[alias] = append(alias2dbModels[alias], dm)
	} else {
		alias2dbModels[alias] = []dbModel{dm}
	}
}

// SyncDB 将定义的model同步到数据库
func SyncDB(args ...string) {
	alias := "default"
	switch len(args) {
	case 1:
		alias = args[0]
	}
	dbModels, ok := alias2dbModels[alias]
	if ok && len(dbModels) > 0 {
		creatingList := make([]interface{}, 0)
		updatingList := make([]interface{}, 0)
		gdb := GetDB(alias)
		mig := gdb.Migrator()
		for _, dbModel := range alias2dbModels[alias] {
			if mig.HasTable(dbModel) {
				updatingList = append(updatingList, dbModel)
			} else {
				creatingList = append(creatingList, dbModel)
			}
		}
		if len(updatingList) > 0 {
			err := mig.AutoMigrate(updatingList...)
			if err != nil {
				panic(err)
			}
		}
		if len(creatingList) > 0 {
			err := mig.CreateTable(creatingList...)
			if err != nil {
				panic(err)
			}
		}
	}
}
