# redis-helper
## 这是一个批量处理redis中数据的工具，可以批量删除key、批量为key设置失效时间。
### redis配置：
``` json
"Redis":{  
"Address":"127.0.0.1",//redis地址  
"Port":"6379",        //redis端口  
"Password":"",        //密码  
"DBIndex":4           //redis数据库index  
}  
```
### 处理规则配置：
``` json
"Rule":{  
"Pattern":["*"],      //key的匹配方式，支持多个匹配规则  
"ScanCount":1000,     //scan函数每次扫描的key的数量  
"Method":"expire",    //处理方法（包括删除delete、设置失效时间expire）  
"ExpireTime":259200   //设置失效时间（单位：秒）  
}  
```
