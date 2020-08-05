package main

import (
	"log"
	"time"
	"strconv"
	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
	eos "github.com/eoscanada/eos-go"
)

/*
var contract_name = "alanogoddess"
var connection_db = "root:root@/alanogoddess"
var connection_eos = "http://103.45.157.202:8888"
*/
var contract_name = "alanogoddess"
var connection_db = "root:root@/alanogoddess"
var connection_eos = "https://nodes.get-scatter.com:443"

type betaccount struct {
	Id uint64 `json:"id"`
	User eos.AccountName `json:"user"`
	Quantity eos.Asset `json:"quantity"`
	Inviter eos.AccountName `json:"inviter"`
}

type globalindex struct {
	Id uint64 `json:"id"`
	Gindex uint64 `json:"gindex"`
}

type betstate struct {
	Round uint64 `json:"round"`
	Users uint64 `json:"users"`
	Bets uint64 `json:"bets"`
	Quantity eos.Asset `json:"quantity"`
	Inviteamount eos.Asset `json:"inviteamount"`
	Tpoolamount eos.Asset `json:"tpoolamount"`
	Apoolamount eos.Asset `json:"apoolamount"`
	Blockprefix uint64 `json:"blockprefix"`
	Blocknum uint64 `json:"blocknum"`
	Result uint64 `json:"result"`
	Btime uint64 `json:"btime"`
	Otime uint64 `json:"otime"`
	Prizeamount eos.Asset `json:"prizeamount"`
}

type betoffer struct {
	Id uint64 `json:"id"`
	Round uint64 `json:"round"`
	User eos.AccountName `json:"user"`
	Quantity eos.Asset `json:"quantity"`
	Bets uint64 `json:"bets"`
	Number uint64 `json:"number"`
}

type innerledger struct {
	User eos.AccountName `json:"user"`
	Round uint64 `json:"round"`
	Bounds uint64 `json:"bounds"`
}

func init_db() {
	
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	{
		sql := `
			create table if not exists numberindex(
			id int
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
		
		//
		sql = "delete from numberindex"
		stmt_delete, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt_delete.Close()
		
		_, err = stmt_delete.Exec()
		if(err != nil) {
			panic(err.Error())
		}
		
		//
		sql = "insert into numberindex values(?)"
		for i := 1;i <= 12;i++ {
			stmt_insert, err := db.Prepare(sql)
			if err != nil {
				panic(err.Error())
			}
			defer stmt_insert.Close()
			
			_, err = stmt_insert.Exec(i)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	
	{
		sql := `
			create table if not exists globalindex(
			id bigint,
			gindex bigint
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
	}
		
	{
		sql := `
			create table if not exists betaccount(
			id bigint,
			user varchar(16),
			quantity varchar(16),
			inviter varchar(16),
			quantity_a bigint
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
	}
	
	{
		sql := `
			create table if not exists betstate(
			round bigint,
			users bigint,
			bets bigint,
			quantity varchar(32),
			inviteamount varchar(32),
			tpoolamount varchar(32),
			apoolamount varchar(32),
			blockprefix bigint,
			blocknum bigint,
			result bigint,
			btime varchar(32),
			otime varchar(32),
			prizeamount varchar(32),
			quantity_a bigint,
			inviteamount_a bigint,
			tpoolamount_a bigint,
			apoolamount_a bigint,
			prizeamount_a bigint,
			btime_a bigint,
			otime_a bigint
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
	}
	
	
	{
		sql := `
			create table if not exists betoffer(
			round bigint,
			id bigint,
			user varchar(16),
			quantity varchar(32),
			bets bigint,
			number bigint,
			quantity_a bigint
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
	}
		
	{
		sql := `
			create table if not exists innerledger(
			round bigint,
			user varchar(16),
			bounds bigint
			);`
			
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec()
		if(err != nil) {
			panic(err.Error())
		}
	}
}

func sync_globalindex() {
	
	//
	eos_api := eos.New(connection_eos)
	_, err := eos_api.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	
	//
	//
	gettable_request := eos.GetTableRowsRequest {
		Code: contract_name,
		Scope: contract_name,
		Table: "globalindex",
	}
	
	var index uint64
	gettable_response, err := eos_api.GetTableRows(gettable_request)
	if err != nil {
		panic(err.Error())
	}
	
	var globalindexs []*globalindex
	err = gettable_response.BinaryToStructs(&globalindexs)
	if err != nil {
		panic(err.Error())
	}
	
	if len(globalindexs) == 0 {
		return ;
	}
	index = globalindexs[0].Gindex;
	
		
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	// exist?
	sql := "select count(*) from globalindex where id = ?"
	stmt_exist, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_exist.Close()
	
	var counter int
	err = stmt_exist.QueryRow(4229443000054317056).Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	
	//
	if counter == 1 {
		var sql = "update globalindex set gindex = ? where id = ?"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec(index, 4229443000054317056)
		if err != nil {
			panic(err.Error())
		}
	} else {
		var sql = "insert into globalindex values(?,?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		_, err = stmt.Exec(4229443000054317056,index)
		if err != nil {
			panic(err.Error())
		}
	}
	
	log.Printf("sync glocalindex successful")	
}

func sync_betaccount() {
			
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	// exist?
	sql := "select count(*) from betaccount"
	stmt_exist, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_exist.Close()
	
	var counter int64
	err = stmt_exist.QueryRow().Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	log.Println(counter)
	
	//
	eos_api := eos.New(connection_eos)
	_, err = eos_api.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	
	//
	//
	gettable_request := eos.GetTableRowsRequest {
		Code: contract_name,
		Scope: contract_name,
		Table: "betaccount",
		LowerBound: strconv.FormatInt(counter, 10),
		Limit: 100,
		//LowerBound: strconv.FormatInt(counter, 10),
	}
	
	gettable_response, err := eos_api.GetTableRows(gettable_request)
	if err != nil {
		panic(err.Error())
	}
	
	var betaccounts []*betaccount
	err = gettable_response.BinaryToStructs(&betaccounts)
	if err != nil {
		panic(err.Error())
	}
	
	for i := 0;i < (len(betaccounts));i++ {
		sql = "insert into betaccount values(?,?,?,?,?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		//log.Printf(betplayers[i].User)
		//log.Printf(betplayers[i].Proxyer)
		_, err = stmt.Exec(
			(betaccounts[i].Id),
			string(betaccounts[i].User),
			(betaccounts[i].Quantity.String()),
			string(betaccounts[i].Inviter),
			(betaccounts[i].Quantity.Amount))
		if err != nil {
			panic(err.Error())
		}	
	}
	
	log.Printf("sync betaccount successful")	
}


func sync_betstate() {
			
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	// exist?
	sql := "select count(*) from betstate"
	stmt_exist, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_exist.Close()
	
	var counter int64
	err = stmt_exist.QueryRow().Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	log.Println(counter)
	
	//
	eos_api := eos.New(connection_eos)
	_, err = eos_api.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	
	var all_index,update_index,add_index int64
	if counter == 0 {
		all_index = 0
		update_index = -1
		add_index = 0
	} else {
		all_index = counter - 1
		update_index = 0
		add_index = 1
	}
	
	//
	//
	gettable_request := eos.GetTableRowsRequest {
		Code: contract_name,
		Scope: contract_name,
		Table: "betstate",
		LowerBound: strconv.FormatInt(all_index, 10),
		Limit: 100,
		//LowerBound: strconv.FormatInt(counter, 10),
	}
	
	gettable_response, err := eos_api.GetTableRows(gettable_request)
	if err != nil {
		panic(err.Error())
	}
	
	var betstates []*betstate
	err = gettable_response.BinaryToStructs(&betstates)
	if err != nil {
		panic(err.Error())
	}
	
	for i := add_index;i < int64(len(betstates));i++ {
		sql = "insert into betstate values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		//log.Printf(betplayers[i].User)
		//log.Printf(betplayers[i].Proxyer)
		//g.Println(betstates[i].Result)
		bt := time.Unix(int64(betstates[i].Btime/1000000),0)
		bt_str := bt.Format("2006-01-02 15:04:05")
		ot := time.Unix(int64(betstates[i].Otime/1000000),0)
		ot_str := ot.Format("2006-01-02 15:04:05")
		_, err = stmt.Exec(
			(betstates[i].Round),
			(betstates[i].Users),
			(betstates[i].Bets),
			(betstates[i].Quantity.String()),
			(betstates[i].Inviteamount.String()),
			(betstates[i].Tpoolamount.String()),
			(betstates[i].Apoolamount.String()),
			(betstates[i].Blockprefix),
			(betstates[i].Blocknum),
			(betstates[i].Result),
			(bt_str),
			(ot_str),
			(betstates[i].Prizeamount.String()),
			(betstates[i].Quantity.Amount),
			(betstates[i].Inviteamount.Amount),
			(betstates[i].Tpoolamount.Amount),
			(betstates[i].Apoolamount.Amount),
			(betstates[i].Prizeamount.Amount),
			(betstates[i].Btime),
			(betstates[i].Otime))
		if err != nil {
			panic(err.Error())
		}	
	}
	
	if update_index != -1 && len(betstates) > 0 {
		sql = "update betstate set users = ?,bets = ?,quantity = ?,inviteamount = ?,tpoolamount = ?,apoolamount = ?,blockprefix = ?,blocknum = ?,result = ?,btime = ?,otime = ?,prizeamount = ?,quantity_a = ?,inviteamount_a = ?,tpoolamount_a = ?,apoolamount_a = ?,prizeamount_a = ?,btime_a = ?,otime_a = ? where round = ?"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		//log.Printf(betplayers[i].User)
		//log.Printf(betplayers[i].Proxyer)
		//log.Println(betstates[counter - 1].Counter)
		//log.Println(betstates[counter - 1].Round)
		//log.Println(betstates[update_index].Otime)
		bt := time.Unix(int64(betstates[update_index].Btime/1000000),0)
		bt_str := bt.Format("2006-01-02 15:04:05")
		ot := time.Unix(int64(betstates[update_index].Otime/1000000),0)
		ot_str := ot.Format("2006-01-02 15:04:05")
		_, err = stmt.Exec(
			(betstates[update_index].Users),
			(betstates[update_index].Bets),
			(betstates[update_index].Quantity.String()),
			(betstates[update_index].Inviteamount.String()),
			(betstates[update_index].Tpoolamount.String()),
			(betstates[update_index].Apoolamount.String()),
			(betstates[update_index].Blockprefix),
			(betstates[update_index].Blocknum),
			(betstates[update_index].Result),
			(bt_str),
			(ot_str),
			(betstates[update_index].Prizeamount.String()),
			(betstates[update_index].Quantity.Amount),
			(betstates[update_index].Inviteamount.Amount),
			(betstates[update_index].Tpoolamount.Amount),
			(betstates[update_index].Apoolamount.Amount),
			(betstates[update_index].Prizeamount.Amount),
			(betstates[update_index].Btime),
			(betstates[update_index].Otime),
			(betstates[update_index].Round))
		if err != nil {
			panic(err.Error())
		}
	}
	
	log.Printf("sync betstate successful")	
}


func sync_betoffer() {
			
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	eos_api := eos.New(connection_eos)
	_, err = eos_api.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	
	//
	//
	gettable_request := eos.GetTableRowsRequest {
		Code: contract_name,
		Scope: contract_name,
		Table: "betoffer",
		Limit: 1000,
		//LowerBound: strconv.FormatInt(counter, 10),
	}
	
	gettable_response, err := eos_api.GetTableRows(gettable_request)
	if err != nil {
		panic(err.Error())
	}
	
	var betoffers []*betoffer
	err = gettable_response.BinaryToStructs(&betoffers)
	if err != nil {
		panic(err.Error())
	}
	
	if len(betoffers) == 0 {
		return;
	}
	gindex := betoffers[0].Round;
	
	// exist?
	sql := "select count(*) from betoffer where round = ?"
	stmt_exist, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_exist.Close()
	
	var counter int64
	err = stmt_exist.QueryRow(gindex).Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	log.Println(counter)
	
	
	for i := counter;i < int64(len(betoffers));i++ {
		sql = "insert into betoffer values(?,?,?,?,?,?,?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		//log.Printf(betplayers[i].User)
		//log.Printf(betplayers[i].Proxyer)
		_, err = stmt.Exec(
			gindex,
			(betoffers[i].Id),
			string(betoffers[i].User),
			betoffers[i].Quantity.String(),
			betoffers[i].Bets,
			betoffers[i].Number,
			betoffers[i].Quantity.Amount)
		if err != nil {
			panic(err.Error())
		}	
	}
	
	log.Printf("sync betoffer successful")	
}

func sync_innerledger() {
			
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	eos_api := eos.New(connection_eos)
	_, err = eos_api.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	
	//
	//
	gettable_request := eos.GetTableRowsRequest {
		Code: contract_name,
		Scope: contract_name,
		Table: "innerledger",
		Limit: 1000,
		//LowerBound: strconv.FormatInt(counter, 10),
	}
	
	gettable_response, err := eos_api.GetTableRows(gettable_request)
	if err != nil {
		panic(err.Error())
	}
	
	var innerledgers []*innerledger
	err = gettable_response.BinaryToStructs(&innerledgers)
	if err != nil {
		panic(err.Error())
	}
	
	if len(innerledgers) == 0 {
		return;
	}
	
	gindex := innerledgers[0].Round
	
	// remove all
	sql := "delete from innerledger where round = ?"
	stmt_remove, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_remove.Close()
	
	_, err = stmt_remove.Exec(gindex)
	if err != nil {
		panic(err.Error())
	}
	
	for i := 0;i < (len(innerledgers));i++ {
		sql = "insert into innerledger values(?,?,?)"
		stmt, err := db.Prepare(sql)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		
		//log.Printf(betplayers[i].User)
		//log.Printf(betplayers[i].Proxyer)
		_, err = stmt.Exec(
			gindex,
			string(innerledgers[i].User),
			innerledgers[i].Bounds)
		if err != nil {
			panic(err.Error())
		}	
	}
	
	log.Printf("sync innerledger successful")	
}

func sync_db() {
	sync_globalindex()
	sync_betaccount()
	sync_betstate()
	sync_betoffer()
	sync_innerledger()
	
}

func main() {
	// init
	init_db()
	
	for {
		time.Sleep(10 * time.Second)
		sync_db()
	}
}
