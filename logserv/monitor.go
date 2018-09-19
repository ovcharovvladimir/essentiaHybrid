/*****************************************************************************
 * ESSENTIA MONITOR
 * (c) Copyright Essentia (2018)
 * https://github.com/DisposaBoy/GoSublime
 * DATE : 
 *****************************************************************************/
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

/*****************************************************************************
 *   Title        : Connection to DB
 *   Initialisation Service
 * 	 Date         : 2018-10-11
 *	 Description  : Initialization DB Connect
 *   Author       : Savchenko Arthur
 *****************************************************************************/
func init() {
    
	session, err := r.Connect(r.ConnectOpts{Database: "wrk", ReadTimeout: time.Second * 250, WriteTimeout: time.Second *10 })

  	// Обработка ошибок
	if err != nil {
	   Inf("init","Error connection to database.", "f")
	   return
	}
	
	// Settig for connection
	session.SetMaxOpenConns(200)
	// session.SetMaxIdleConns(200)
	sessionArray = append(sessionArray, session)

	fmt.Println("Log server is started....")
    Inf("init", "Session is created.",  "i")
}

/*****************************************************************************
 * DATETIME         : 28-07-2015 12:44
 * DESCRIPTION      : Стартовая процедура
 * NOTES            : Запуск сервиса с параметрами
 *****************************************************************************/
func main() {
    // Flags
    Port:=flag.String("Port",":5898", "Input Pleas Port for service. By default 5898.")	
    flag.Parse()	
    
    // Route
	http.HandleFunc("/",                      StartPage)        // Start Page
    http.HandleFunc("/login/",                Login)            // Registartion
    http.HandleFunc("/static/",               StaticPage)       // Link to static page
    http.HandleFunc("/info/",                 InfoPage)         // Link to static page
    http.HandleFunc("/about/",                AboutPage)        // Link to static page

    
    // DB
	http.HandleFunc("/db/start/",             Db_Prepea)        // Created database
	http.HandleFunc("/db/del/",               Db_Delete)        // Clear базы
	http.HandleFunc("/db/drop/",              Table_log_drop)   // Drop log table

    // Test
	http.HandleFunc("/tst/add/",              Db_testlog)       // Test add record
	http.HandleFunc("/tst/cli/",              Test_os)          // Test call client monitor
	
	// Admin panel
	http.HandleFunc("/api/add/",              AddInfStr)        // Add inf to log !
	http.HandleFunc("/api/admin/",            Admin_panel)      // Admin panel 
	http.HandleFunc("/api/test/add/",         AddInf)           // Add inf to log test
	

	// Reports 
	http.HandleFunc("/rep/test/",             Rep_log_test)     // Test login operation
	http.HandleFunc("/rep/log/",              Rep_log_journal)  // View in HTML report
    http.HandleFunc("/rep/json/",             Rep_log_json)     // Export to json format
    http.HandleFunc("/rep/graph/",            Rep_graph)        // Simple graph page
    http.HandleFunc("/rep/count/",            Rep_count)        // Count rec in table

    // Admin
    http.HandleFunc("/adm/idx/",              Admi_create_index)        // Count rec in table

    // Cli
    http.HandleFunc("/cli/send/",             Cli_send)         // Export to json format    
        
    // Information about load server...
    Inf("main", "Server is started on the port: " + *Port, "i")

    srvhttp := &http.Server{Addr:*Port, ReadTimeout:10*time.Minute,WriteTimeout:10*time.Minute}
    err:=srvhttp.ListenAndServe()

	if err!=nil{
	   Inf("main", err.Error(), "w")
	   Inf("main","Error start service!" , "f")
	}

	// err:=http.ListenAndServe(*Port, nil)
	// if err!=nil{
	//    Inf("main", err.Error(), "w")
	//    Inf("main","Error start service!" , "f")
	// }
}


//************************************************************
//  Name    : Static pages and library CSS, Javascrip and etc... 
//  Date    : 12-09-2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : /static/....
//************************************************************
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

//************************************************************
//  Name    : Start page
//  Date    : 05.09.2018 15:37
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : /
//************************************************************
func StartPage(w http.ResponseWriter, req *http.Request) {
     // w.Write([]byte(Html))
      PG("index.html", "Log Journal", "View report journal log.", nil, w, req)	
}

//************************************************************
//  Name    : Loging page (Draft)
//  Date    : 05.09.2018 17:37
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : 
//************************************************************
func Login(w http.ResponseWriter, req *http.Request) {
     fmt.Println("Ok")
     w.Write([]byte("OK"))
}


//************************************************************
//  Name    : Prepea data for work in databale
//  Date    : 06.09.2018 15:37
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : db/start/
//************************************************************
func Db_Prepea(w http.ResponseWriter, req *http.Request) {
   
    r.DBCreate("wrk").Exec(sessionArray[0])      
    Table_drop()
    Inf("db prepea","Database and log table was created!" , "i")
}


//************************************************************
//  Name    : Delete data in databale
//  Date    : 07.09.2018 19:37
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : /db/del/
//************************************************************
func Db_Delete(w http.ResponseWriter, req *http.Request) {
	 rd := r.DeleteOpts{Durability: "soft", ReturnChanges: false}
     r.DB("wrk").Table("log").Delete(rd).Exec(sessionArray[0])
     Inf("db delete","Database log was be clear!" , "i")
     s:=[]byte("Таблица полность очищена.")
     w.Write(s)
}

//************************************************************
//  Name    : Test insert record to database
//  Date    : 07.09.2018 11:37
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : /tst/add/
//************************************************************
func Db_testlog(w http.ResponseWriter, req *http.Request) {
     
     var Dat LogStruct

     // charge of data
     Dat.Operation  = "Test initial operation"	
     Dat.Project    = "Block-chain-beacon"
     Dat.Module     = "RewardCount"
     Dat.Datetime   =  CurTime 
     Dat.Status     = "Info"
     Dat.BlockId    = "000015783b764259d382017d91a36d206d0600e2cbb3567748f46a33fe9297cf"
     Dat.AccountId  = "AccountID/ContractID"
     Dat.CreateTime = time.Now().Format("2006-01-02 14:55") 
     
     // Test add in database
     Db_LogAdd(Dat)      
     Inf("db test", Dat.BlockId,  "i")
}

//************************************************************
//  Name    : Добавление одной записи в лог таблицу 
//  Date    : 05-09-2018 15:37
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : 
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

//************************************************************
//  Name    : Добавление информации в таблицу в формате JSON
//  Date    : 05.09.2018 15:37
//  Company : Essentia
//  Number  : 
//  Module  : Monitoring
//  Path    : 
//************************************************************
func AddInf(w http.ResponseWriter, req *http.Request) {
    
    // Data
    m := make(map[string]interface{})

    // Conflict rule for dpuble index field
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

	// Adding unix time in automatically mode
	u     := time.Now().UnixNano() / 1000000 
	u_str := time.Unix(0, u*1000000).Format("2006-01-02T15:04:05.000")

	m["Unxtime"] = fmt.Sprintf("%d", u) 
	m["Timestr"] = u_str

	// Check error
	if errj != nil {
	   Inf("Add inf", "Document don't have body.",  "e")
	}

    go func(){ 
	    // Add new document
		erri:=r.DB("wrk").Table("log").Insert(m, Conflictrule).Exec(sessionArray[0])

	    // Check error
		if erri != nil {
		   Inf("Add doc", "Error insert message to log table",  "e")
		}
	}()

    // Ok insert
	Inf("Add inf", "Succeseful adding record to log table.",  "i")
}

//************************************************************
//  Name    : Основная процедура добавления записи в лог таблицу 
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
      if len(t)==1 {
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
// Name       : Report view log journal (test)
// Date       : 12-09-2018
// Company    : Essentia
// Author     : Svachenko Arthur
// Module     : monitor 
// URL        : /rep/log/  
// Usage      : /rep/log/10   - manula set get records  
// By Default : Get 100 records in reverse order
// **********************************************************

func Rep_log_test(w http.ResponseWriter, req *http.Request) {
    
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
// Name       : Report view log journal
// Date       : 14-09-2018
// Company    : Essentia
// Author     : Svachenko Arthur
// Module     : monitor 
// URL        : /rep/log/  
// Usage      : /rep/log/10   - manula set get records  
// By Default : Get 100 records in reverse order
// **********************************************************
func Rep_log_journal(w http.ResponseWriter, req *http.Request) {
	   
	   ro:= r.RunOpts{ArrayLimit: 20000000}
	   // ord:=r.OrderByOpts{Index: "Datetime"}

	   p := req.URL.Path[len("/rep/log/"):]

       l := 100
    
       // Check param
       if p!=""{
          l=Sti(p) 
       } 
	
	   var Data []Mst
	   tb:=r.DB("wrk").Table("log")
	   // Rk, er := r.DB("wrk").Table("log").Without("id","Id").OrderBy(r.Desc("Datetime")).Limit(l).Run(sessionArray[0],ro)
       
       // Created Index
       // tb.IndexCreate("Datetime").Exec(sessionArray[0])
       // tb.IndexWait().Exec(sessionArray[0])
       // Inf("Rep-log", "Index created.",  "i")
       

	    Rk, er := tb.OrderBy(r.Desc("Datetime")).Limit(l).Run(sessionArray[0],ro)
	   // Rk, er := tb.Without("id","Id").OrderBy(ord).Limit(l).Run(sessionArray[0],ro)
   
	   // Error
	   if er != nil {
	   	  fmt.Println("Error open table", er.Error())
	      // Inf("Rep-log", "Error open table log.",  "e")
	   }
   
	   defer Rk.Close()

	   err:=Rk.All(&Data)

	   if err!=nil{
	   	fmt.Println("Error read log table :", err.Error())
	   }
      
       // Get page
       PG("journal.html", "Log Journal", "View report journal log.", Data, w, req)
}

// **********************************************************
// Name       : Report view log journal in format JSON
// Date       : 14-09-2018
// Company    : Essentia
// Author     : Svachenko Arthur
// Module     : monitor 
// URL        : /rep/json/  
// Usage      : /rep/json/10   - manula set get records  
// By Default : Get 100 records in reverse order
// **********************************************************
func Rep_log_json(w http.ResponseWriter, req *http.Request) {
    ro:= r.RunOpts{ArrayLimit: 20000000}
    p := req.URL.Path[len("/rep/json/"):]
    l := 100
    
    if p!=""{
       l=Sti(p) 
    } 
	
	var response []Mst
    res, er := r.DB("wrk").Table("log").Without("id","Id").OrderBy(r.Desc("Datetime")).Limit(l).Run(sessionArray[0],ro)

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
// Name       : Admin panel page
// Date       : 14-09-2018 12:01
// Company    : Essentia
// Author     : Svachenko Arthur
// Module     : monitor 
// URL        : /rep/admin/
// Usage      : 
// **********************************************************
func Admin_panel(w http.ResponseWriter, req *http.Request) {
     PG("adm.html", "Admi Panel", "Administrtion and monitoring", nil, w,req)
}

//************************************************************
//  Name    : Test call from client
//  Date    : 10-09-2018 15:56
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : monitor
//  Module  : 
//  Params  : 
//************************************************************
func Cli_send(w http.ResponseWriter, req *http.Request) {
     go Send_Info("Prysm","Reward","Samples text","Infotest","X0afdfdfdsxzcvdfgffffdgfgfdgdfgdfdfqerert","AccountID",CTM())
     Inf("Cli send", "Send to server ", "i")     
}

//************************************************************
//  Name    : Send to log server information
//  Date    : 12-09-2018 21:26
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//  Usage   :  Send_Info("Hybrid","Worker","Add blockchain","Info", "Blockid","Accountid",time.Now().Format("2006-01-02"))
//************************************************************
func Send_Info(Project, Module, Opertion, Status, BlockId, AccountID, CreateTime string ){
    url        := "http://18.223.111.231:5898/api/add/"+Project+"*"+Module+"*"+Opertion+"*"+Status+"*"+BlockId+"*"+AccountID+"*"+CreateTime
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

//************************************************************
//  Name    : Sample call procedure in program mode
//  Date    : 10-09-2018 21:26
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : monitoring
//************************************************************
func Test_os(w http.ResponseWriter, req *http.Request){
	 cmd:= exec.Command("./clmon",  "Призма", "Modul", "Send message to log file", "info")
	 err:= cmd.Run()

	 if err!=nil{
	 	fmt.Println(err.Error())
	    Inf("send", "Error", "w")		
	 }
     Inf("send", "Send to server ", "w")	
}

//************************************************************
//  Name    : Graph page test 
//  Date    : 10-09-2018 16:44
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : /rep/graph/
//  Module  : 
//************************************************************
func Rep_graph(w http.ResponseWriter, req *http.Request){
     p := req.URL.Path[len("/rep/graph/"):]
     PG("grap"+p+".html", "Admi Panel", "Administrtion and monitoring", nil, w,req)
}

//************************************************************
//  Name    : Graph page test 
//  Date    : 14-09-2018 16:44
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : /rep/count/
//  Module  : 
//************************************************************
func Rep_count(w http.ResponseWriter, req *http.Request){
	var cnt int64
	res, er := r.DB("wrk").Table("log").Count().Run(sessionArray[0])

	if er != nil {
       Inf("Rep JSON", "Error read table",  "e") 
	}

	defer res.Close()
	er = res.One(&cnt)

    // Check error
	if er != nil {
		Inf("Rep JSON", "Error read data form table log",  "w") 
    }		
     resp:=Int64toStr(cnt)
     fmt.Println(resp)
    // w.Header().Set("Content-Type", "application/text; charset=utf-8")
	w.Write([]byte(resp))

    // Call generator page
    // PG("cnt.html", "Count data", "Count rec in log table", cnt, w, req)

}

//************************************************************
//  Name    : PG (Page Generator)
//  Date    : 10-09-2018 15:56
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : Creator page from template
//  Module  : GPage
//  Usage  :  PG("page.html", "Title", "Description", Data, w, req)
//************************************************************
func PG(PageNameHtml, Title, Description string, Data []Mst,  w http.ResponseWriter, req *http.Request) {
	 Dt        := Mst{"Dts": Data, "Title": Title, "Descript": Description, "Datrep": CTM()}
	 fp        := path.Join("tmp", PageNameHtml)                  
	 tmpl, err := template.ParseFiles(fp, "tmp/main.html")  
	 
	 defer func(){
	 	recover()
	 }()

	 // Err(err, "Error events template execute.")
	 if err!=nil{
   	    Inf("Cli send", "Bad read MAIN template or absent file.", "w")     	
   	 } else {
   	     errf:= tmpl.Execute(w, Dt)
   	 
   	     if errf!=nil{
   	        Inf("Cli send", "Bad read template or absent file.", "e")     	
   	     }	
   	 }
}

//************************************************************
//  Name    : Info Page
//  Date    : 14-09-2018 20:04
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : /info
//  Module  : 
//************************************************************
func InfoPage(w http.ResponseWriter, req *http.Request){
	defer func(){
		recover()
	}()
    PG("info.html", "Info", "Info", nil, w,req)
}

//************************************************************
//  Name    : About Page
//  Date    : 14-09-2018 20:44
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : /info
//  Module  : 
//************************************************************
func AboutPage(w http.ResponseWriter, req *http.Request){
	defer func(){recover()}()
    PG("about.html", "About", "About programm", nil, w,req)
}


//************************************************************
//  Name    : Created index in table 
//  Date    : 17-20-2018
//  Author  : Svachenko Arthr
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Admi_create_index(w http.ResponseWriter, req *http.Request){
resp:=`  
<html>
     <h1>Index was created successful! </h1>
</html>`
	err:=r.DB("wrk").Table("log").IndexCreate("Datetime").Exec(sessionArray[0])
	if err!=nil{
		Inf("Admin", "Error created index in table log.", "e")     	
	}
	 Inf("Admin", "Index was created susseccfully.", "i")     	
     w.Write([]byte(resp))	
}


//************************************************************
//  Name    : Table drop is API
//  Date    : 18-20-2018 22:22
//  Author  : Svachenko Arthr
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Table_log_drop(w http.ResponseWriter, req *http.Request){
     Table_drop()
}

//************************************************************
//  Name    : Table drop is API
//  Date    : 18-20-2018 22:22
//  Author  : Svachenko Arthr
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Table_drop(){	
    tc:= r.TableCreateOpts{Durability: "soft"} 
    r.DB("wrk").TableDrop("log").Exec(sessionArray[0])
    Inf("db prepea","Table was dropted!", "i")
    r.DB("wrk").TableCreate("log",tc).Exec(sessionArray[0])
    r.DB("wrk").Table("log").IndexCreate("Datetime").Exec(sessionArray[0])
    Inf("Db", "Table log was created susseccfully.", "i")     	
}
