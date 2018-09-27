// 
// Copyright 2018 Essentia
// https://github.com/ovcharovvladimir/essentiaHybrid/tree/master/waletproxy
// 
package main
import (
	"fmt"
    "io/ioutil"
	"time"
    "log"
    "net/http"
    "encoding/json"
    "net/http/httputil"
    "net/url"
)

// Node structure
type Node  struct {
             Ip string     `json:"Ip"`
             Port string   `json:"Port"`
             Note string   `json:"Note"`
             Status string `json:"Status"`
} 

// Sdetting structure
type Sett struct {
    Mainport string `json:"Mainport"`
    Version  string `json:"Version"`
    Nodes   []Node
}

//************************************************************
//  Name    : Main  
//************************************************************
func main() {

    Host  := ActiveNode()
    Sett  := ReadSettingFile()
    Port  := ":"+ Sett.Mainport

    // For test with other resources
    // Host="localhost:5555"
    // Host="www.youtube.com"
    proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme:"http", Host:Host})

    // Proxy
    proxy.Director = func(req *http.Request) {
        fmt.Println("Host redirect:", Host) 
        req.Host       = Host
        req.URL.Host   = Host
        req.URL.Scheme = "http"    
    }

    // Proxy 
     http.Handle("/",   proxy)
    
    // Routing
    http.HandleFunc("/nodes/",        ShowNodes)              // Show all nodes
    http.HandleFunc("/active/",       Active_node)            // Show active first node for work 
    http.HandleFunc("/actives/",      Active_nodes)           // Show active nodes 
    http.HandleFunc("/down/",         Down_nodes)             // Show down nodes 
    http.HandleFunc("/test/",         Api_test)               // Test service response 
    http.HandleFunc("/admin/",        Api_admin)              // Test service response 
       
    // View settings
    fmt.Println("\n Start Service Port ", Port, "\n","Active Work Node    :", Host)
    err := http.ListenAndServe(Port, nil)
    
    // Error
    if err != nil {
       log.Println("Error service",err.Error())
    }
}

//************************************************************
//  Name    : Show All Node 
//************************************************************
func ShowNodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()
    
    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         
         if ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Active"
         }else{
            nd.Status="Not available"
         }
         WPr(nd,w)
    }
}

//************************************************************
//  Name    : Show All Node 
//************************************************************
func Active_node(w http.ResponseWriter, req *http.Request) {
     ft:="Active note :" + ActiveNode()
     Wprn(ft,w)
}

//************************************************************
//  Name    : Show All Node 
//************************************************************
func Active_nodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         if ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Active"
            WPr(nd, w)
         }
    }
}

//************************************************************
//  Name    : Text  output rep
//************************************************************
func Wprn(Txt string, w http.ResponseWriter){
     w.Write([]byte(Txt))
}

//************************************************************
//  Name    : Format  output rep
//************************************************************
func WPr(nd Node, w http.ResponseWriter){
     ft:=fmt.Sprintf("Node :%-15s Ip : %-15s Port : %-6s Status : %-20s \n",nd.Note,nd.Ip,nd.Port,nd.Status)
     w.Write([]byte(ft))
}

//************************************************************
//  Name    : Show All Node 
//************************************************************
func Down_nodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()
    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         
         if !ChekNodeWork(nd.Ip, nd.Port){
             nd.Status ="Is Down"
             WPr(nd,w)
          }
    }
}

//************************************************************
//  Name    : Test
//************************************************************
func Api_test (wr http.ResponseWriter, rq *http.Request) {
     t:="Test service " + time.Now().Format("02/01/2006 15:45")
     fmt.Println(t)
     Wprn(t,wr)
}

//************************************************************
// Set the proxied request's host to the destination host (instead of the
// source host).  e.g. http://foo.com proxying to http://bar.com will ensure
// that the proxied requests appear to be coming from http://bar.com
//
// For both this function and queryCombiner (below), we'll be wrapping a
// Handler with our own HandlerFunc so that we can do some intermediate work
//************************************************************
func sameHost(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.Host = r.URL.Host
        handler.ServeHTTP(w, r)
    })
}

//************************************************************
// Allow cross origin resource sharing
//************************************************************
func addCORS(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
        handler.ServeHTTP(w, r)
    })
}

//************************************************************
//  Name    : Get info about first active note 
//************************************************************
func ActiveNode() string {
    Rip:=""
    Nd :=ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
        if ChekNodeWork(nd.Ip, nd.Port){
            Rip=nd.Ip+":"+ nd.Port                
            break
          }
    }
    return Rip
}

//************************************************************
//  Name    : Reading  setting from file json 
//************************************************************
func ReadSettingFile() Sett {
    // Setting structure
    var m Sett
    
    // Open file with settings
    // (in Unix ./config.json)
    file, err := ioutil.ReadFile("./setting.json")

    // Error
    if err != nil {
       log.Println("Error reading setting.") 
       return Sett{}
    }
    
    errj:=json.Unmarshal([]byte(file), &m)
    if errj!=nil{
       log.Println("Error ummarshaling.") 
       return Sett{}
    }
    return m
}

//************************************************************
//  Name    : Check node 
//************************************************************
func ChekNodeWork(Ip,Port string) bool {
     timeout := time.Duration(700 * time.Millisecond )
     client  := http.Client{Timeout: timeout}
     resp, err := client.Get("http://"+Ip+":"+Port)
    
    if err!=nil {
       return false      
    }

    if resp.StatusCode==200{
       return true
    }else{
       return false 
    }
}


// *********************************************************************
// Admin panel
// *********************************************************************
func Api_admin(w http.ResponseWriter, req *http.Request) {

html:=`
<!DOCTYPE html>
<html lang="en">
<head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Admin</title>

        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css">
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
        <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
        <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"></script>

        <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Roboto|Varela+Round">
        <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">

        <style>
              body {color: #060f21b5; font-family: 'Roboto';  font-size: 16px;}
              .card-header1 {background-color: #108c84; color:white; font-weight:bold;}
        </style>
</head>

<body>
    <div class="container">
        <h2>Proxy inform panel</h2>

        <!-- Grey with black text -->
        <nav class="navbar navbar-expand-sm bg-dark navbar-dark">
              <ul class="navbar-nav">
                <li class="nav-item active"> <a class="nav-link" href="#">Active</a></li>
                <li class="nav-item"> <a class="nav-link" href="#">Link</a> </li>
                <li class="nav-item"> <a class="nav-link" href="#">Link</a> </li>
                <li class="nav-item"> <a class="nav-link disabled" href="#">Disabled</a> </li>
              </ul>
        </nav>
        

        <!--
            <div class="card bg-info text-white">
                <div class="card-body">Proxy inform panel</div>
            </div>
        -->

        <br>
        <div class="card">
                    <div class="card-header" style="background-color:#red;">
                         <h5>Admin panel</h5>
                    </div>

                    <div class="card-body">
                          <a href="/nodes/">Preview all nodes</a><br>
                          <a href="/active/">Preview <b>first</b> active node</a><br>
                          <a href="/actives">Preview all <b>active</b> nodes</a><br>
                          <a href="/down/">Preview all <b>Disabled</b>nodes</a><br>
                          <a href="/test/">Test service</a><br>
                    </div>      
                    <div class="card-footer">Essentia</div>           
                </div>

        </div>    
    </div>    
</body>
</html>
 `
Wprn(html,w)
}


