/****************************************************************************************
 *	ESSENTIA CLIENT MONITOR
 *  Copyright ESSENTIA  (2018)
 *  12-09-2018 
 *  Ver 0.001 
 ****************************************************************************************/
package main


import (
	"net/http"
	"fmt"
	"os"
)

var(
    ipserv="http://18.223.111.231:5898"
)


/****************************************************************************************
 * DATETIME         : 06-09-2018 17:44
 * NOTES            : Send log information to server
 ****************************************************************************************/
func main() {
     
     // Args
     arg       :=os.Args
     Project   :=arg[1]
     Module    :=arg[2]
     Operation :=arg[3]
     Status    :=arg[4]
     BlockId   :=arg[5]
     AccountId :=arg[6]
  
     Send_Info(Project, Module, Operation, Status, BlockId,AccountId)
}

// *************************************************************
//  Send to log server information
//  Send_Info("EssentiaHybrid", "worker", "", "Info","ccc34554zxcxzcaddfdf3445cvdv","acc12344")
// *************************************************************
func Send_Info(Project, Module, Operation, Status, BlockId, AccountId string ){
    ipserv := "http://18.223.111.231:5898"
    url    :=  ipserv+"/api/add/"+Project+"*"+Module+"*"+Operation+"*"+Status+"*"+BlockId+"*"+AccountId
	re,errr:=  http.NewRequest("GET", url, nil)
	
	if errr!=nil{
       fmt.Println(errr.Error()) 
       return
	}

	res, erd := http.DefaultClient.Do(re)
	if erd!=nil{
       fmt.Println(erd.Error()) 
       return
	}

	defer res.Body.Close()
}
