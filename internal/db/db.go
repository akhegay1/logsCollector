package db

import (
	"bufio"
	"database/sql"
	"fmt"

	"logsCollector/internal/crypto"
	"logsCollector/pkg/mwutils"
	"os"
	"strings"

	"log"
)

var Db *sql.DB

var (
	host     = ""
	port     = ""
	user     = ""
	password = ""
	dbname   = ""
)
var words []string

func Connect() string {
	log.Println("Connect", "started")
	key := mwutils.Key
	///////////!!!!!!!ENCRYPT PASSW!!!!!!!/////

	//passw := "admin123"
	//foo := Encrypt(key, passw)
	//fmt.Println("foo", foo)

	////////////////FILE/////
	conf, err := os.Open("conf")
	if err != nil {
		log.Println("connect", "failed opening file conf: "+err.Error())
		return fmt.Sprintf("failed opening file conf: %s", err)
	}
	defer conf.Close()

	sc := bufio.NewScanner(conf)

	for sc.Scan() {
		str := sc.Text() // GET the line string
		words = strings.Fields(str)
	}

	if err := sc.Err(); err != nil {
		log.Println("Connect", "scan file error: "+err.Error())
		return fmt.Sprint("Connect", "scan file error: "+err.Error())
	}

	host = words[0]
	port = words[1]
	user = words[2]
	password = words[3]
	dbname = words[4]

	user = crypto.Decrypt(key, user)
	password = crypto.Decrypt(key, password)
	//fmt.Println("password", password)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer Db.Close()

	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	return "Successfully connected!"
}
