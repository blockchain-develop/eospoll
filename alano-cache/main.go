package main

import (
	"log"
	"strconv"
	"net/http"
	"github.com/emicklei/go-restful"
	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var connection_db = "root:root@/alanogoddess"

type User struct {
	Name         string
	Members      []User
}

type UserResource struct {
	Users        map[string]User
}

type betstate struct {
	Round uint64 `json:"round"`
	Users uint64 `json:"users"`
	Bets uint64 `json:"bets"`
	Quantity string `json:"quantity"`
	Inviteamount string `json:"inviteamount"`
	Apoolamount string `json:"apoolamount"`
	Result uint64 `json:"result"`
	Otime string `json:"otime"`
	Prizeamount string `json:"prizeamount"`
	Btime_a uint64 `json:"btime_a"`
}

type numberpool struct {
	Number uint64 `json:"number"`
	Bets uint64 `json:"bets"`
	Precent string `json:"precent"`
	Prize string `json:"prize"`
	Prize_a uint64 `json:"prize_a"`
}

type usernumber struct {
	User string `json:"user"`
	Number int `json:"number"`
	CanBet int `json:"canbet"`	
}

func (u UserResource) Register(container *restful.Container) {
	{
		ws := new(restful.WebService)
		ws.Path("/proxyermembers").
			Consumes(restful.MIME_JSON, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_JSON)
			
		ws.Route(ws.GET("/{user-id}").To(u.ProxyerMembers))
		container.Add(ws)
	}
	
	{
		ws := new(restful.WebService)
		ws.Path("/historybetstate").
			Consumes(restful.MIME_JSON, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_JSON)
			
		ws.Route(ws.GET("/all").To(u.BetState))
		container.Add(ws)
	}
	
	{
		ws := new(restful.WebService)
		ws.Path("/numberpool").
			Consumes(restful.MIME_JSON, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_JSON)
			
		ws.Route(ws.GET("/all").To(u.NumberPool))
		container.Add(ws)
	}
	
	{
		ws := new(restful.WebService)
		ws.Path("/betoffercheck").
			Consumes(restful.MIME_JSON, restful.MIME_JSON).
			Produces(restful.MIME_JSON, restful.MIME_JSON)
			
		ws.Route(ws.GET("/").To(u.BetOffer))
		container.Add(ws)
	}	
}

func (u UserResource) FindProxyerMembers(id string) []string {
	
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	sql := "select count(*) from betaccount where inviter = ?"
	stmt_counter, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_counter.Close()
	
	var counter int
	err = stmt_counter.QueryRow(id).Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	
	//
	sql = "select user from betaccount where inviter = ?"
	stmt_query, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_query.Close()
	
	rows, err := stmt_query.Query(id)
	if err != nil {
		panic(err.Error())
	}
	
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("query successful")
	
	//
	var members = make([]string, counter)
	var col1 []byte
	var i = 0
	
	for rows.Next() {
		if i >= counter {
			break;
		}
		
		err = rows.Scan(&col1)
		if err != nil {
			log.Printf(err.Error())
		}
		
		//
		members[i] = string(col1)
		i++
	}
	
	//
	return members[:i]
}

func (u UserResource) ProxyerMembers(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	log.Printf(id)
	
	members := u.FindProxyerMembers(id)
	response.WriteEntity(members)
}

func (u UserResource) BetState(request *restful.Request, response *restful.Response) {
	
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	sql := "select count(*) from betstate"
	stmt_counter, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_counter.Close()
	
	var counter int
	err = stmt_counter.QueryRow().Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	
	//
	sql = "select round,users,bets,quantity,inviteamount,apoolamount,result,otime,prizeamount,btime_a from betstate order by round desc limit 50"
	stmt_query, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_query.Close()
	
	rows, err := stmt_query.Query()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("query successful")
	
	//
	var members = make([]betstate, counter)
	var round, users, bets, result,btime_a uint64
	var quantity, inviteamount, apoolamount, otime, prizeamount []byte
	var state betstate
	var index = 0
	
	for rows.Next() {
		if index >= counter {
			break;
		}
		
		err = rows.Scan(&round, &users, &bets, &quantity, &inviteamount, &apoolamount,
			&result, &otime, &prizeamount, &btime_a)
		if err != nil {
			log.Printf(err.Error())
		}
		
		if result > 10000000 {
			continue;
		}
		
		state.Round = round
		state.Users = users
		state.Bets = bets
		state.Quantity = string(quantity)
		state.Inviteamount = string(inviteamount)
		state.Apoolamount = string(apoolamount)
		state.Result = result
		state.Otime = string(otime)
		state.Prizeamount = string(prizeamount)
		state.Btime_a = btime_a
		//
		members[index] = state
		index ++
	}
	
	//
	ret_result := members[:index]
	response.WriteEntity(ret_result)
}

func (u UserResource) NumberPool(request *restful.Request, response *restful.Response) {
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	sql := "select betoffer.number as number,sum(betoffer.bets) as numberbets,sum(betoffer.bets)/betstate.bets as precent,(betstate.quantity_a+betstate.apoolamount_a/2)/sum(betoffer.bets) as amount from betoffer left join betstate on betoffer.round = betstate.round where betstate.round in (select gindex from globalindex) group by betoffer.number"
	stmt_query, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_query.Close()
	
	rows, err := stmt_query.Query()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("query successful")
	
	//
	var members = make([]numberpool, 13)
	for i := 0;i < 13;i ++ {
		members[i].Number = uint64(i)
		members[i].Bets = 0
		members[i].Precent = "0.00%"
		members[i].Prize = "0.0000 EOS"
	}
	var number, numberbets uint64
	var precent,amount float64
	var all_bets uint64 = 0
	var all_amount uint64 = 0
	
	for rows.Next() {
		err = rows.Scan(&number, &numberbets, &precent, &amount)
		if err != nil {
			log.Printf(err.Error())
		}
		
		members[number].Number = number
		members[number].Bets = numberbets
		members[number].Precent = strconv.FormatFloat((precent*100), 'f', 2, 64) + "%"
		members[number].Prize = strconv.FormatFloat((amount/10000), 'f', 4, 64) + " EOS"
		members[number].Prize_a = uint64(amount)
		
		all_bets += numberbets
		all_amount = numberbets * uint64(amount)
	}
	
	members[0].Bets = all_bets;
	members[0].Precent = "0.00%";
	members[0].Prize = strconv.FormatFloat((float64(all_amount)/10000), 'f', 4, 64) + " EOS"
	members[0].Prize_a = all_amount
	
	//
	response.WriteEntity(members)
}

func (u UserResource) BetOffer(request *restful.Request, response *restful.Response) {
	user := request.QueryParameter("user")
	number, err := strconv.Atoi(request.QueryParameter("number"))
	
	log.Printf(user)	
	log.Println(number)
	
	//
	db, err := sql.Open("mysql",connection_db)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("open successful")
	defer db.Close()
	
	//
	sql := "select count(distinct user) as counter from betoffer where number = ?"
	stmt_query, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_query.Close()
	
	var counter int
	err = stmt_query.QueryRow(number).Scan(&counter)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("query successful")
	
	
	//
	sql = "select count(*) as counter from betoffer where number = ? and user = ?"
	stmt_query1, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt_query1.Close()
	
	var counter1 int
	err = stmt_query1.QueryRow(number,user).Scan(&counter1)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("query successful")
	
	//
	var state usernumber
	state.User = user
	state.Number = number
	bret := checkusernumber(number, counter, counter1)
	if bret == true {
		state.CanBet = 1
	} else {
		state.CanBet = 0
	}
	
	//
	response.WriteEntity(state)
}

func checkusernumber(number int,counter int,exist int) bool {
	if exist != 0 {
		return true
	}
	
	if number == 11 || number == 12 {
		if counter == 1 {
			return false
		}
	} else if number == 9 || number == 10 {
		if counter == 2 {
			return false
		}
	} else if number == 7 || number == 8 {
		if counter == 4 {
			return false
		}
	} else if number == 5 || number == 6 {
		if counter == 8 {
			return false
		}
	} else if number == 3 || number == 4 {
		if counter == 16 {
			return false
		}
	} else if number == 1 || number == 2 {
		if counter == 32 {
			return false
		}
	}
	return true
}

func main() {
	wsContainer := restful.NewContainer()	
	wsContainer.Router(restful.CurlyRouter{})
	u := UserResource{map[string]User{}}
	u.Register(wsContainer)
	
	//
	cors := restful.CrossOriginResourceSharing {
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST"},
		CookiesAllowed: false,
		Container:      wsContainer}

	wsContainer.Filter(cors.Filter)
	wsContainer.Filter(wsContainer.OPTIONSFilter)
	
	log.Printf("start listening on localhost:7780")
	server := &http.Server{Addr:":7780",Handler:wsContainer}
	log.Fatal(server.ListenAndServe())
}
