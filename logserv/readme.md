# Test  


Для тестирования можно вызвать файл с обязательными параметрами
1. **Project**       - имя проекта
2. **Module**        - Имя модуля
3. "**Short send**"  - краткое сообщение для мониторинга
4. **Status**        - Статус сообщения (Info, Warn, Error)

## Пример проверки

```sh
#!/bin/bash

clear
echo Test work clien for log server

# Start
./clmon Project Module "Short send to srver" Status
```

