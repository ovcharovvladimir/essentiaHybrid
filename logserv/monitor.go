/**********************************************************************************************************************************************************************************************
 * ESSENTIA MONITOR
 * (c) Copyright Essentia (2018)
 * https://github.com/DisposaBoy/GoSublime
 * DATE : 
 ***********************************************************************************************************************************************************************************************/
package main

import (
	"net/http"
	"fmt"
	"flag"
	"io/ioutil"
	"time"
	"strings"
	"encoding/json"
	"html/template"
    "path"
    "os/exec"
   r "github.com/dancannon/gorethink"
)

/***************************************************************************************************************************************
 *   Title        : Connection to DB
 *   Initialisation Service
 * 	 Date         : 2018-10-11
 *	 Description  : Initialization DB Connect
 *   Author       : Savchenko Arthur
 ****************************************************************************************************************************************/
func init() {
	session, err := r.Connect(r.ConnectOpts{Database: "wrk"})

  	// Обработка ошибок
	if err != nil {
	   Inf("init","Error connection to database.", "f")
	   return
	}
	
	// Settig for connection
	session.SetMaxOpenConns(200)
	session.SetMaxIdleConns(200)
	sessionArray = append(sessionArray, session)
    Inf("init", "Session is created.",  "i")
}

/*******************************************************************************************************************************
 * DATETIME         : 28-07-2015 12:44
 * DESCRIPTION      : Стартовая процедура
 * NOTES            : Запуск сервиса с параметрами
 *******************************************************************************************************************************/
func main() {
    // Flags
    Port:=flag.String("Port",":5898", "Input Pleas Port for service. By default 5898.")	
    flag.Parse()	
    
    // Route
	http.HandleFunc("/",                      StartPage)        // Start Page
    http.HandleFunc("/login/",                Login)            // Registartion
    http.HandleFunc("/static/",               StaticPage)       // Link to static page

    // DB
	http.HandleFunc("/db/start/",             Db_Prepea)        // Created database
	http.HandleFunc("/db/del/",               Db_Delete)        // Clear базы

    // Test
	http.HandleFunc("/tst/add/",              Db_testlog)       // Test add record
	http.HandleFunc("/tst/cli/",              Test_os)          // Test call client monitor
	
	// Admin panel
	http.HandleFunc("/api/admin/",            Admin_panel)      // Admin panel 
	http.HandleFunc("/api/test/add/",         AddInf)           // Add inf to log test
	http.HandleFunc("/api/add/",              AddInfStr)        // Add inf to log

	// Reports 
	http.HandleFunc("/rep/test/",             Rep_log)          // Test login operation
	http.HandleFunc("/rep/log/",              Rep_log_journal)  // View in HTML report
    http.HandleFunc("/rep/json/",             Rep_log_json)     // Export to json format
    http.HandleFunc("/rep/graph/",            Rep_graph)        // Export to json format

    // Cli
    http.HandleFunc("/cli/send/",             Cli_send)         // Export to json format    
        
    // Info
    Inf("main", "Server is started on the port " + *Port, "i")

	err:=http.ListenAndServe(*Port, nil)
	if err!=nil{
	   Inf("main", err.Error(), "w")
	   Inf("main","Error start service!" , "f")
	}
}

/********************************************************************************************************************************
 *   Static Page
 *
 *   /static/....
 *********************************************************************************************************************************/
func StaticPage(w http.ResponseWriter, r *http.Request) {
	// Allows
	w.Header().Set("Access-Control-Allow-Origin", "*") 
	/* Allows
	   if origin := r.Header().Get("Origin"); origin != "" {
	    w.Header().Set("Access-Control-Allow-Origin", origin)
	    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	    w.Header().Set("Access-Control-Allow-Headers",  "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	    w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
	   }
	*/
	//  File static page
	http.ServeFile(w, r, r.URL.Path[1:])
}

/********************************************************************************************************************************
 *  TITLE             : Start page
 *  DESCRIPTION       : Start page with description
 *  DATE              : 05.09.2018 15:37
 **********************************************************************************************************************************/
func StartPage(w http.ResponseWriter, req *http.Request) {
     s:=Html
     w.Write([]byte(s))
}

/********************************************************************************************************************************
 *  TITLE             : Log test
 *  DESCRIPTION       : Вход в сервис
 *  DATE              : 05.09.2018 17:37
 **********************************************************************************************************************************/
func Login(w http.ResponseWriter, req *http.Request) {
     fmt.Println("Ok")
     w.Write([]byte("OK"))
}

/********************************************************************************************************************************
 *  TITLE             : Подготовка базы данных 
 *  DESCRIPTION       : Вход в сервис
 *  DATE              : 05.09.2018 15:37
 **********************************************************************************************************************************/
func Db_Prepea(w http.ResponseWriter, req *http.Request) {
     r.DBCreate("wrk").Exec(sessionArray[0])      
     r.DB("wrk").TableCreate("log").Exec(sessionArray[0])
     Inf("db prepea","Database and log table was created!" , "i")
}

/********************************************************************************************************************************
 *  TITLE             : Подготовка базы данных 
 *  DESCRIPTION       : Delete table
 *  DATE              : 05.09.2018 15:37
 **********************************************************************************************************************************/
func Db_Delete(w http.ResponseWriter, req *http.Request) {
     go r.DB("wrk").Table("log").Delete().Exec(sessionArray[0])
     Inf("db delete","Database log was be clear!" , "i")
     s:=[]byte("Таблица полность очищена.")
     w.Write(s)
}

/********************************************************************************************************************************
 *  TITLE             : Test insert
 *  DESCRIPTION       : Add to table test information
 *  DATE              : 05.09.2018 15:37
 *  URL               :                                 
 **********************************************************************************************************************************/
func Db_testlog(w http.ResponseWriter, req *http.Request) {
     
     var Dat LogStruct

     Dat.Operation  = "Test initial operation"	
     Dat.Project    = "Block-chain-beacon"
     Dat.Module     = "RewardCount"
     Dat.Datetime   =  CurTime 
     Dat.Status     = "Info"
     Dat.BlockId    = "000015783b764259d382017d91a36d206d0600e2cbb3567748f46a33fe9297cf"
     Dat.AccountId  = "AccountID/ContractID"
     Dat.CreateTime = time.Now().Format("2006-01-02 14:55") 
    
     Inf("db test", Dat.BlockId,  "w")
     
     // Test add in database
     Db_LogAdd(Dat)      
}

//************************************************************
//  Name    : Добавление одной записи в лог таблицу 
//  Date    : 05-09-2018 15:37
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Db_LogAdd(Dat LogStruct){
	 Conflictrule := r.InsertOpts{Conflict: "replace", Durability:"soft", ReturnChanges: false}
	 defer func(){
        recover()
	  }()

	 go func(){
	    err:=r.DB("wrk").Table("log").Insert(Dat, Conflictrule).Exec(sessionArray[0])
	    if err!=nil{
           Inf("Db log Add","Error insert to log." , "e")
	    }
    }() 
}

/********************************************************************************************************************************
 *  TITLE             : Delete all records 
 *  DESCRIPTION       : Clear log
 *  DATE              : 05.09.2018 18:37
 **********************************************************************************************************************************/
func DbLogDelete(Dat LogStruct){
	 r.DB("wrk").Table("log").Delete().Exec(sessionArray[0])
	 fmt.Println("Information was deleted...")
}

/********************************************************************************************************************************
 *  TITLE             : Добавление информации в таблицу в формате JSON
 *  DESCRIPTION       : Вход в сервис
 *  DATE              : 05.09.2018 15:37
 *                                 
 **********************************************************************************************************************************/
func AddInf(w http.ResponseWriter, req *http.Request) {
    
    // Data
    m := make(map[string]interface{})

    // Conflict rule 
    Conflictrule := r.InsertOpts{Conflict: "replace", Durability:"soft", ReturnChanges: false}

	// Чтение тела документа
	reads, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	// Error
	if len(string(reads)) == 0 {
	   Inf("AddInf", "Document Dont have body ....",  "e")
	}
	
	// Load data 
	errj := json.Unmarshal([]byte(reads), &m)

	// Замена значений полей в структуре новыми сформироваными значениями даты времени
	// Замена счетчика текущим значением
	u     := time.Now().UnixNano() / 1000000 
	u_str := time.Unix(0, u*1000000).Format("2006-01-02T15:04:05.000")

	m["Unxtime"] = fmt.Sprintf("%d", u) 
	m["Timestr"] = u_str

	// Check error
	if errj != nil {
	   Inf("Add inf", "Document Dont have body ....",  "e")
	}

    // Add new document
	erri:=r.DB("wrk").Table("log").Insert(m,Conflictrule).Exec(sessionArray[0])

    // Check error
	if erri != nil {
	   Inf("Add doc", "Error insert to log table",  "e")
	}

	Inf("Add inf", "Succeseful adding record to database",  "i")
}


//************************************************************
//  Name    : Основная процедурадобавления записи в лог таблицу 
//  Date    : 05-10-2018 12:57
//  Company : Essentia
//  Number  : A01
//  Module  : API
//  Usage   : .../api/add/proj*module*operation*status*BlockId*AccountId*CreateTime
//  URl     : /api/add/
//************************************************************
func AddInfStr(w http.ResponseWriter, req *http.Request) {
      var Dat LogStruct
      var Chk bool

      p := req.URL.Path[len("/api/add/"):]
      t := strings.Split(p, "*")

      // Chek count parameters
      if len(t)==1{
         Inf("Add Rec to log.", "Bad parameters",  "f")     	
         return
      }

      // Load data
      Dat.Project     = t[0]         
      Dat.Module      = t[1]
      Dat.Operation   = t[2]
      Dat.Status      = t[3]
      Dat.BlockId     = t[4]
      Dat.AccountId   = t[5]
      Dat.CreateTime  = t[6]
      Dat.Sys()

      s:=time.Now()

      // Add to database log
      go Db_LogAdd(Dat)     

      
      // Check time insert
      if Chk {      
         // go Db_LogAdd(Dt)     
         f:=time.Now()
         r:=f.Sub(s)
         fmt.Println("Time in insert :",r)
         Inf("Add Str", "Succeseful adding record to database",  "i") 
      }
}

// **********************************************************
// Report view log
// **********************************************************
func Rep_log(w http.ResponseWriter, req *http.Request) {
    
    // Cтруктура
	type Inventory struct {
		 Country   string
		 Index     string
	}

	// Данные для зарядки
	sweaters := []Inventory{
	            {"Добавление нового документа",             "/api/"},
	            {"Работа с документом по ID документу",     "/api/id/"},
	            {"Получение набора документов по фильтру",  "/api/filter/"},
	            {"Получение максимального сиквенса",        "/api/seq/"},
	            {"Информация по сервису",                   "/docs/info/"},
	}

	// {"Описание действий", "Путь к сервисам"},
	sss := `<!DOCTYPE html>
	            <htmL>
	            <head>
	                  <title>Head Office</title>
                      <meta http-equiv="Content-Type" content="text/html; charset=windows 1251">
                      <meta http-equiv="Content-Language" content="Ru-ru">
                      <meta name="viewport" content="width=device-width, initial-scale=1">

					  <link rel="stylesheet" href="http://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css">
					  <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
					  <script src="http://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>

					   <style type="text/css" >
					          body      {margin:10px;}
					          html      {color: #3B3B3B; font-size:10px; font-family:Calibri;}
					          td1       {border: 1px gray solid; padding:5px;}
					         .ddblock   {padding:10px;}
					         .hhd1      {border:1px gray solid; padding:5px; background-color: #E7E7E7; font-weight: bold;}
					         .hdsf1     {color: #3366FF;	font-weight: bold;	font-size:17px;}
					         .btn-group {width:350px}
					         .mb        {width:1000px; left:50%; margin-left:-500px; position: absolute;}
					   </style>
                     </head>

	                <body>
                    <h3>Описание процедур </h3>
                    <br>

	                <table class="table">
	                    <thead>
                                <tr>
                                <th class="hhd1">Описание процедур</th>
                                <th class="hhd1">Путь к функциям</th>
                                </tr>
                        </thead>

                        <tbody>
				               {{ range. }}
				               <tr>
							       <td> <b>{{printf "%s" .Country}}</b> </td>
							       <td> <a href="{{printf "%s" .Index }}">{{printf "%s" .Index }}</a> </td>
					           </tr>
					           {{ end }}
		                 </tbody>
		           </table>
            </body>
            </html>`

	tmpl, err := template.New("test").Parse(sss)

	// Error
	if err != nil {
		Inf("Init", "Bad read template",  "w")
		w.WriteHeader(208)
	}


	err = tmpl.Execute(w, sweaters)

	// Error
	if err != nil {
	   Inf("Rep", "Bad execute template.",  "e")
	   w.WriteHeader(208)
	}
}	

// **********************************************************
// Report view log
// URL :/rep/log/
// **********************************************************
func Rep_log_journal(w http.ResponseWriter, req *http.Request) {
	  p := req.URL.Path[len("/rep/log/"):]
      l := 100
    
       // Check param
       if p!=""{
          l=Sti(p) 
       } 
	
	   var Data []Mst
	   Rk, er := r.DB("wrk").Table("log").Without("id","Id").OrderBy(r.Desc("Datetime")).Limit(l).Run(sessionArray[0])
   
	   // Error
	   if er != nil {
	      Inf("Rep-log", "Error open table log.",  "e")
	   }
   
	   defer Rk.Close()

	   Rk.All(&Data)
   
	   Dt        := Mst{"Dts": Data, "Title": "Log journal", "Descript": "Log", "Datrep": CTM()}
	   fp        := path.Join("tmp", "journal.html")                  
	   tmpl, err := template.ParseFiles(fp, "tmp/main.html")  
	   Err(err, "Error events template execute.")
   
	   errf       := tmpl.Execute(w, Dt)
	   Err(errf, "Error events template execute.")
}

// **********************************************************
// Report view log
// URL :/rep/json/
// **********************************************************
func Rep_log_json(w http.ResponseWriter, req *http.Request) {

    p := req.URL.Path[len("/rep/json/"):]
    l :=100
    
    if p!=""{
       l=Sti(p) 
    } 
		
	var response []Mst
	res, er := r.DB("wrk").Table("log").Without("id","Id").OrderBy(r.Desc("Datetime")).Limit(l).Run(sessionArray[0])

	if er != nil {
       Inf("Rep JSON", "Error read table",  "e") 
	}

	defer res.Close()
	er = res.All(&response)

    // Check error
	if er != nil {
		Inf("Rep JSON", "Error read data form table log",  "w") 
	} else {
		data, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	}
}

// **********************************************************
// Report view log
// URL : /rep/admin/
// **********************************************************
func Admin_panel(w http.ResponseWriter, req *http.Request) {
     PG("adm.html", "Admi Panel", "Administrtion and monitoring", nil, w,req)
}

//************************************************************
//  Name    : PG PAge Generator
//  Date    : 10-09-2018 15:56
//  Company : Essentia
//  Number  : Creator page from template
//  Module  : GPage
//************************************************************
func PG(PageNameHtml, Title, Description string, Data []Mst,  w http.ResponseWriter, req *http.Request) {
	   Dt        := Mst{"Dts": Data, "Title": Title, "Descript": Description, "Datrep": CTM()}
	   fp        := path.Join("tmp", PageNameHtml)                  
	   tmpl, err := template.ParseFiles(fp, "tmp/main.html")  
	   Err(err, "Error events template execute.")
   	   errf      := tmpl.Execute(w, Dt)
	   Err(errf, "Error events template execute.")
}

// **********************************************************
//  Call from client
// **********************************************************
func Cli_send(w http.ResponseWriter, req *http.Request) {

     go Send_Info("Prysm","Reward","Samples text","Infotest","X0afdfdfdsxzcvdfgffffdgfgfdgdfgdfdfqerert","AccountID",CTM())
     Inf("Cli send", "Send to server ", "i")     
}


// *************************************************************
// Title   : Send to log server information
// Date    : 12-02-2018 21:20
// Library : "net/http"
// Usage   :  Send_Info("Hybrid","Worker","Add blockchain","Info", "Blockid","Accountid",time.Now().Format("2006-01-02"))
// *************************************************************
func Send_Info(Project, Module, Opertion, Status, BlockId, AccountID, CreateTime string ){
    url      := "http://18.223.111.231:5898/api/add/"+Project+"*"+Module+"*"+Opertion+"*"+Status+"*"+BlockId+"*"+AccountID+"*"+CreateTime
	req,  err  := http.NewRequest("GET", url, nil)
	
	if err!=nil{
       Inf("Send Info", "Error request.", "e")           
	}

	res, erd := http.DefaultClient.Do(req)
	if erd!=nil{
       Inf("Send Info", "Error client connection.", "e")           
	}

	defer res.Body.Close()
}

// *************************************************************
// Title        : Send to log server information
// Description  : Old version wit return data form server
// *************************************************************
func Send_Info_old(Project, Module, Opertion, Status,BlockId,AccountID,CreateTime string  ){
    url    := "http://18.223.111.231:5898/api/add/"+Project+"*"+Module+"*"+Opertion+"*"+Status+"*"+BlockId+"*"+AccountID+"*"+CreateTime
	req,_  := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Service-token", "d41aee26-cc94-E9ff-e9a5-f0701845624b")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

//************************************************************
//  Name    : Test вызова клиента для монитора
//  Date    : 
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Test_os(w http.ResponseWriter, req *http.Request){
	 cmd:= exec.Command("./clmon",  "Призма", "Modul", "Пример передачи примера", "info")
	 err:= cmd.Run()

	 if err!=nil{
	 	fmt.Println(err.Error())
	    Inf("send", "Error", "w")		
	 }
     Inf("send", "Send to server ", "w")	
}

//************************************************************
//  Name    : Grap page 
//  Date    : 10-09-2018 16:44
//  Company : Essentia
//  Number  : /rep/graph/
//  Module  : 
//************************************************************
func Rep_graph(w http.ResponseWriter, req *http.Request){
     p := req.URL.Path[len("/rep/graph/"):]
     PG("grap"+p+".html", "Admi Panel", "Administrtion and monitoring", nil, w,req)
}
