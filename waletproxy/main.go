// Copyright 2018 
// https://github.com/ovcharovvladimir/essentiaHybrid/tree/master/waletproxy

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
    "strconv"
    "bytes"
    "html/template"
    "path"
    "github.com/fatih/color"
)

// Node structure
type Node  struct {
     Ip        string   `json:"Ip"`
     Port      string   `json:"Port"`
     Note      string   `json:"Note"`
     Status    string   `json:"Status"`
     Disabled  string   `json:"Disabled"`
     Datetime  string   `json:"Datetime"`
} 

// Sdetting structure
type Sett struct {
     Mainport string `json:"Mainport"`
     Version  string `json:"Version"`
     Nodes   []Node
}

// Transport structure
type transport struct {
     http.RoundTripper
}

var _ http.RoundTripper = &transport{}
type Mst map[string]interface{}                              



//**********************************************************
// Main 
//**********************************************************
func main() {

    // Set color
    Whites  := color.New(color.FgWhite).PrintlnFunc()
    FgMag   := color.New(color.FgHiYellow)
    Cyan    := color.New(color.Bold, color.FgCyan)
    FgRed   := color.New(color.Bold, color.FgRed).PrintlnFunc()
    FgGreen := color.New(color.Bold, color.FgGreen).PrintlnFunc()
    
    FgGreen("Proxy server")
    Whites("Version: 1.01 (Testing)")

    // Active node
    Host  := ActiveNode()      
    Sett  := ReadSettingFile()
    Port  := ":"+ Sett.Mainport

    FgMag.Print("Active node : ")    
    FgRed(Host)
    Cyan.Printf("Listen proxy port %s \n", Port)    

    // Check active node
    if Host== "" {
       Host="Error: All nodes disabled."
    }

    // For test with other resources
    // Host="localhost:5555"
    // Host="www.youtube.com"
    proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme:"http", Host:Host})
    proxy.Transport = &transport{http.DefaultTransport}

    // Proxy
    proxy.Director = func(req *http.Request) {
        // Allows
        req.Header.Set("Content-Type","application/json;charset=utf-8")
        req.Header.Set("Access-Control-Allow-Origin","*")
        req.Header.Set("Access-Control-Allow-Headers","X-Requested-With")
        req.Header.Set("X-Forwarded-For",Host)
       
        req.Host       = Host
        req.URL.Host   = Host
        req.URL.Scheme = "http"    
    }

    // Proxy 
    http.Handle("/",                  proxy)                  // Proxy 
    
    // Routing
    http.HandleFunc("/nodes/",        ShowNodes)              // Show all nodes
    http.HandleFunc("/active/",       Active_node)            // Show active first node for work 
    http.HandleFunc("/actives/",      Active_nodes)           // Show active nodes 
    http.HandleFunc("/down/",         Down_nodes)             // Show down nodes 
    http.HandleFunc("/test/",         Api_test)               // Test service response 
    http.HandleFunc("/admin/",        Api_admin)              // Admin panel
    http.HandleFunc("/report/node/",  Nodes_report)           // Node report
  
    err := http.ListenAndServe(Port, nil)
    
    // Error
    if err != nil {
       log.Println("Error start service.",err.Error())
    }
}

//************************************************************
//  Transport
//************************************************************
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

  resp, err = t.RoundTripper.RoundTrip(req)
  if err != nil {
     return nil, err
  }

  b, err := ioutil.ReadAll(resp.Body)
  if err != nil {
     return nil, err
  }

  err = resp.Body.Close()
  if err != nil {
     return nil, err
  }

  b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
  body := ioutil.NopCloser(bytes.NewReader(b))
  resp.Body = body
  resp.ContentLength = int64(len(b))
  resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
  resp.Header.Set("Content-Type", "application/json; charset=utf-8")
  return resp, nil
}

//************************************************************
// Show all nodes 
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
// Show first active node
//************************************************************
func Active_node(w http.ResponseWriter, req *http.Request) {
     An:=""
     if ActiveNode()!=""{
        An = ActiveNode()
     }else{
        An = "No active node."
     }
     ft:="Active note :" + An
     Wprn(ft,w)
}

//************************************************************
// Show all active nodes
//************************************************************
func Active_nodes(w http.ResponseWriter, req *http.Request) {
    Nd := ReadSettingFile()

    // Looop for nodes
    for _,nd:=range Nd.Nodes{
         if ChekNodeWork(nd.Ip, nd.Port){
            nd.Status="Active"
            WPr(nd, w)
         }
    }
}

//************************************************************
// Text  output rep
//************************************************************
func Wprn(Txt string, w http.ResponseWriter){
     w.Write([]byte(Txt))
}

//************************************************************
// Format  output rep
//************************************************************
func WPr(nd Node, w http.ResponseWriter){
     ft:=fmt.Sprintf("Node :%-15s Ip : %-15s Port : %-6s Status : %-20s \n",nd.Note,nd.Ip,nd.Port,nd.Status)
     w.Write([]byte(ft))
}


//************************************************************
// Show all disabled node 
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
// Test
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
        w.Header().Set("Content-Type", "application/json;charset=utf-8")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
        handler.ServeHTTP(w, r)
    })
}

//************************************************************
// Get info about first active note 
//************************************************************
func ActiveNode() string {
    Rip := ""
    Nd  := ReadSettingFile()

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
// Reading  setting from file json 
//************************************************************
func ReadSettingFile() Sett {
    // Setting structure
    var m Sett
    
    // Open file with settings
    file, err := ioutil.ReadFile("./setting.json")
    if err != nil {
       log.Println("Error reading setting file.") 
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
// Check node 
//************************************************************
func ChekNodeWork(Ip,Port string) bool {
     timeout := time.Duration(500 * time.Millisecond )
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
    
    fp := path.Join("tmp", "node.html")                                 
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")                  
    Err(err, "Error template execute.")

    errf := tmpl.Execute(w, nil)
    Err(errf, "Error template execute.")
}


// *********************************************************************
// Node journal 
// /report/node/
// *********************************************************************
func Nodes_report(w http.ResponseWriter, req *http.Request) {
     Nd:=ReadSettingFile()
     
    
    // Looop for nodes
    for i, nd:=range Nd.Nodes{
          if ChekNodeWork(nd.Ip, nd.Port){
             Nd.Nodes[i].Status="Active"
             Nd.Nodes[i].Disabled="success"
          }else{
             Nd.Nodes[i].Status="Disabled"
             Nd.Nodes[i].Disabled="danger"
          }
          
          Nd.Nodes[i].Datetime=time.Now().Format("02.01.2006 15:04:05")

    }

    Dt:= Mst{"Dts": Nd.Nodes, "Title": "Active Nodes ", "Descript": "Serach" }


    fp := path.Join("tmp", "journal.html")                                 
    tmpl, err := template.ParseFiles(fp, "tmp/main.html")                  
    Err(err, "Error template execute.")

    errf := tmpl.Execute(w, Dt)
    Err(errf, "Error template execute.")


}

/***************************************************************
  Check Eror
 ****************************************************************/
func Err(Er error, Txt string) {
    if Er != nil {
       log.Println("ERROR : " + Txt)
       return
    }
}

















