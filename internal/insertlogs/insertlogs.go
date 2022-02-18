package insertlogs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"logsCollector/internal/db"
	"logsCollector/pkg/crud"
	"logsCollector/pkg/mwutils"
	"net/http"

	tbl "logsCollector/pkg/tables"
)

//ppp

func InsLogs(w http.ResponseWriter, r *http.Request) {
	log.Println("InsLogs", "started")
	defer log.Println("InsLogs", "finished")

	mwutils.SetupResponse(&w, r)
	if r.Method == "OPTIONS" {
		log.Println("InsLogs", "OPTIONS")
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		log.Println("InsLogs", http.StatusText(405))
		return
	}

	// Read body
	b, err := ioutil.ReadAll(r.Body)
	log.Println("InsLogs", string(b))

	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println("InsLogs", err.Error())

	}

	var logs tbl.Logs
	err = json.Unmarshal(b, &logs)

	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println(r, "insLogs", err.Error())
		return
	}

	//log.Println(r, "crud insert", crud.CreateInsert(logs))
	var lastInsertId int16

	txn, _ := db.Db.Begin()
	insq := crud.CreateInsert(logs)

	err = db.Db.QueryRow(insq).Scan(&lastInsertId)

	if err != nil {
		http.Error(w, http.StatusText(500)+" "+err.Error(), 500)
		log.Println("InsLogs", err.Error())
		txn.Rollback()
		return
	}

	err = txn.Commit()
	log.Println("InsLogs", "aft insert commit")

	/* SQL
	err = db.Db.QueryRow("INSERT INTO monit_sch.logs (app, fname,level,msg,reqID,time,userID)"+
		" VALUES($1, $2, $3, $4, $5, $6, $7) returning id", logs.App, logs.Fname, logs.Level, logs.Msg, logs.ReqID, logs.Time, logs.UserID).Scan(&lastInsertId)
	log.Println("aft insert err", err)
	if err != nil {
		http.Error(w, http.StatusText(500)+" "+err.Error(), 500)
		txn.Rollback()
		return
	}

	err = txn.Commit()
	*/

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		log.Println("aft insert commit err", err.Error())
		return
	}

	//log.Println("lastInsertId", lastInsertId)
	//fmt.Fprintf(w, "%d", lastInsertId)
}
