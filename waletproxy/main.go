// Copyright 2018 Essentia
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
//  Date    : 25/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func main() {
     Host  := ActiveNode()
     Sett  := ReadSettingFile()
     Port  := ":"+ Sett.Mainport
     proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host:Host})
    
    proxy.Director = func(req *http.Request) {
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
//  Date    : 27/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func ShowNodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()
    
    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         
         if !ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Activate"
         }else{
            nd.Status="Disabled"
         }
         WPr(nd,w)
    }
}

//************************************************************
//  Name    : Show All Node 
//  Date    : 27/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Active_node(w http.ResponseWriter, req *http.Request) {
     ft:="Active note :" + ActiveNode()
     Wprn(ft,w)
}

//************************************************************
//  Name    : Show All Node 
//  Date    : 27/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Active_nodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         if ChekNodeWork(nd.Ip, nd.Port){
            WPr(nd, w)
         }
    }
}

//************************************************************
//  Name    : Text  output rep
//  Date    : 26/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Wprn(Txt string, w http.ResponseWriter){
     w.Write([]byte(Txt))
}
//************************************************************
//  Name    : Format  output rep
//  Date    : 26/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func WPr(nd Node, w http.ResponseWriter){
     ft:="Node : "+nd.Note +"Ip : "+ nd.Ip+" Port : "+ nd.Port +" Status : " + nd.Status +" \n" 
     w.Write([]byte(ft))
}

//************************************************************
//  Name    : Show All Node 
//  Date    : 27/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
//************************************************************
func Down_nodes(w http.ResponseWriter, req *http.Request) {
    Nd:=ReadSettingFile()
    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         if !ChekNodeWork(nd.Ip, nd.Port){
            ft:="Node :"+nd.Note +"Ip:"+ nd.Ip+" Port:"+ nd.Port +" Active \n"
            Wprn(ft,w)
     }
         
    }
}

//************************************************************
//  Name    : Test
//  Date    : 27/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
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
//  Date    : 25.09.2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
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
//  Date    : 26/09/2018
//  Author  : Svachenko Arthur
//  Company : Essentia
//  Number  : 
//  Module  : 
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
//  Date    : is a snippet.
//  Author  : Svachenko Arthr
//  Company : Essentia
//  Number  : 
//  Module  : 
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
