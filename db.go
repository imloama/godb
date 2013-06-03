package godb
import (
  "fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	//"log"
)

var MySQLPool chan mysql.Conn

//配置内容
var DbConfig DBConfig

type DBConfig struct {
	DBUsername    string
	DBPassword    string
	DBURL         string
	DBName        string
	MAX_POOL_SIZE int
}

func init() {
	//初始化配置文件
	DbConfig.DBUsername = "root"
	DbConfig.DBName = "test"
	DbConfig.DBPassword = ""
	DbConfig.DBURL = "127.0.0.1:3306"
	DbConfig.MAX_POOL_SIZE = 5

}

func GetConnection() mysql.Conn {
	if MySQLPool == nil {
		MySQLPool = make(chan mysql.Conn, DbConfig.MAX_POOL_SIZE)
	}
	if len(MySQLPool) == 0 {
		go func() {
			for i := 0; i < DbConfig.MAX_POOL_SIZE/2; i++ {
				db := mysql.New("tcp", "", DbConfig.DBURL, DbConfig.DBUsername, DbConfig.DBPassword, DbConfig.DBName)
				err := db.Connect()
				if err != nil {
					panic(err)
				}
				putMySQLConn(db)
			}
		}()
	}
	return <-MySQLPool
}

func putMySQLConn(conn mysql.Conn) {
	if MySQLPool == nil {
		MySQLPool = make(chan mysql.Conn, DbConfig.MAX_POOL_SIZE)
	}
	if len(MySQLPool) == DbConfig.MAX_POOL_SIZE {
		conn.Close()
		return
	}
	MySQLPool <- conn

}
