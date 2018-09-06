# Test  


Для тестирования можно вызвать файл с обязательными параметрами
1. **Project**       - имя проекта
2. **Module**        - Имя модуля
3. "**Short send**"  - краткое сообщение для мониторинга
4. **Status**        - Статус сообщения (Info, Warn, Error)

## Start.sh
```sh
#!/bin/bash


clear
echo Start building...

# Enviroument
export GOPATH=$HOME/app
#export GOPATH=$PWD
export GOROOT=$HOME/go
export PATH=$PATH:$GOROOT/bin

# Start
go build -o clmon
./clmon Prrysm Modul Text info
```


## Пример проверки
## test.sh

```sh
#!/bin/bash

clear
echo Test work clien for log server

# Start
./clmon Project Module "Short send to srver" Status
```

