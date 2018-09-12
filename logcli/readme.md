# Test  


Для тестирования можно вызвать файл с обязательными параметрами

|#|Arg|Description|
|--|--|--|
|1.| **Project**|Имя проекта|
|2.| **Module** |Имя модуля|
|3.| "**Operation**"|Краткое сообщение для мониторинга|
|4.| **Status**|Статус сообщения (Info, Warn, Error)|
|5.| **BlockId**|Ид блока|
|6.| **AccountId**|Ид счета|
|7.| **CreateTime**|Дата и время создания блока|


## Building
### Start.sh
```sh
#!/bin/bash


clear
echo Start building...

# Enviroument
export GOPATH=$PWD
export GOROOT=$HOME/go
export PATH=$PATH:$GOROOT/bin

# Start
go build -o clmon
./clmon Prrysm Modul Text info
```


##  Пример тестирования
### Test.sh

```sh
#!/bin/bash

clear
echo Test work clien for log server

# Start
./clmon Project Module "Short send to srver" Status
```
## Использование в коде прогаммы

```golang

// 
// Test вызова клиента для монитора
// 
func Test_os(){
	 
	 cmd:= exec.Command("./clmon",  "Prizm", "Module", "Пример передачи сообщения на сервер", "info", "0x37cc62924b876a043bad996399d3bc15f0f629f01e1ef6c457b1b486681a568c","Acc092","2018-12-09 12:45")
	 err:= cmd.Run()

	 if err!=nil{
	 	fmt.Println(err.Error())
	    Inf("send", "Error message to server", "w")		
	 }

     Inf("send", "Send to log server ", "w")	
}
```



## Примечание :
При вызове третьего параметра если в нем содержится несколько слов - его необходимо брать в кавычки.

```
./clmon Project Module "Этот параметр нужно брать в кавычки" Status
```


