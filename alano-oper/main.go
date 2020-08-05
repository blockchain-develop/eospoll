package main

import (
	"log"
	"time"
	"strconv"
	
	eos "github.com/eoscanada/eos-go"
)

type eosnode_state struct {
	Online      int
	EosAPI      eos.API
}

var eosnodes_addr = [2]string{"https://nodes.get-scatter.com:443","https://nodes.get-scatter.com:443"}
var eosnodes_state = [2]eosnode_state{}
var eosnodes_number = 2

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

func eosnodes_init() {
	for i := 0;i < eosnodes_number;i ++ {
		eosnodes_state[i].Online = 0
		eosnodes_state[i].EosAPI = eos.New(eosnodes_addr[i])
	}
}

func eosnodes_check(index int) {
	log.Prinf(index)
	if index < 0 || index >= node_number {
		return;
	}
	
	//
	info, err := eosnodes_state[index].EosAPI.GetInfo()
	if err != nil {
		eosnodes_state[i].Online = 0
		eosnodes_state[i].EosAPI = eos.New(eosnodes_addr[i])
	} else {
		eosnodes_state[i].Online = 1
	}
}

func eosnodes_listen() {
	for i := 0;i < node_number;i ++ {
		time.Sleep(1 * time.Second)
		eosnodes_check(i)
	}
}

func select_node() int {
	for i := 0;i < node_number;i ++ {
		if eosnodes_state[i].Online == 1 {
			return i
		}
	}
	return -1
}

func get_round() int {
	index := select_node()
	if index == -1 {
		return -1
	}
	
	gettable_request := eos.GetTableRowRequest {
		Code : contract_name,
		Scope : contract_name,
		Table : "globalindex",
	}
	
	gettable_reeponse, err := eosnodes_state[index].EosAPI.GetTableRows(gettable_request)
	if err != nil {
		return -1
	}
	
	var globalindexs []*globalindex
	err = gettable_response.BinaryToStructs(&globalindexs)
	if err != nil {
		return -1
	}
	
	if len(globalindexs) == 0 {
		return -1
	}
	
	return globalindexs[0].Gindex
}

func get_last_betstate(round int) (open int,time uint32) {
	index := select_node()
	if index == -1 {
		return -1,0
	}
	
	gettable_request := eos.GetTableRowRequest {
		Code : contract_name,
		Scope : contract_name,
		Table : "betstate",
		LowerBound : strconv.FormatInt(round, 10),
		Limit : 1,
	}
	
	gettable_reeponse, err := eosnodes_state[index].EosAPI.GetTableRows(gettable_request)
	if err != nil {
		return -1,0
	}
	
	var betstates []*betstate
	err = gettable_response.BinaryToStructs(&betstates)
	if err != nil {
		return -1,0
	}
	
	if len(betstates) == 0 {
		return -1,0
	}
	
	result := betstates[0].Result
	if result & 0xff00000000000000 {
		open = 1
	}
	btime := betstates[0].Btime {
		time = btime / 1000000
	}
	return open,time
}

func will_time(time uint32) (nexttime uint32) {
	nexttenminite = ((time / 600) + 1) * 600
	curtime = uint32(Time.Unix())
	if curtime < nexttenminite {
		nexttime = nexttenminite
	} else {
		nexttime = ((curtime / 600) + 1) * 600
	}
}

func open_bet() {
	
}

func close_bet() {
	
}


func main() {
	eosnodes_init()
	go eosnodes_listen()
	
	for {
		round := get_round()
		if round == -1 {
			time.Sleep(1 * time.Second)
			continue
		}
		
		open,time := get_last_betstate(round)
		if open == -1 || open == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		
		oritenminite := ((time / 600) + 1) * 600
		curtime := uint32(Time.Unix())
		nexttenminite := ((curtime / 600)) * 600
		var willtime = nexttenminite
		if oritenminite > nexttenminite {
			willtime = oritenminite
		}
		
		if curtime > willtime && (curtime - willtime) < 10 {
			close_bet()
			time.Sleep(5 * time.Second)
			open_bet()
			time.Sleep(5 * time.Second)
			continue
		} else if curtime < willtime {
			time.Sleep((willtime - curtime) * time.Second)
			continue
		} else {
			time.Sleep((600 - curtime + willtime) * time.Second)
			continue
		}
	}
}
