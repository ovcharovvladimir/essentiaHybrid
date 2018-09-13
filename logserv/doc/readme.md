## For use send command in programm

### Usage  

* add librray "net/http"
* add procedure Set_info
* Call procedure in needed place to program  

### Module

```golang
// *************************************************************
// Title   : Send to log server information
// Date    : 12-02-2018 21:20
// Library : "net/http"
// Usage   :  Send_Info("Hybrid","Worker","Add blockchain","Info", "Blockid","Accountid",time.Now().Format("2006-01-02"))
// *************************************************************
func Send_Info(Project, Module, Opertion, Status, BlockId, AccountID, CreateTime string ){
    url     := "http://18.223.111.231:5898/api/add/"+Project+"*"+Module+"*"+Opertion+"*"+Status+"*"+BlockId+"*"+AccountID+"*"+CreateTime
	re,err  := http.NewRequest("GET", url, nil)
	
	if err!=nil{
       Inf("Send Info", "Error request.", "e")           
	}

	res, erd := http.DefaultClient.Do(re)
	if erd!=nil{
       Inf("Send Info", "Error client connection.", "e")           
	}

	defer res.Body.Close()
}
```
