# redis-helper
## 这是一个批量处理redis中数据的工具，可以批量删除key、批量为key设置失效时间。
### redis配置：
`"Redis":{\n`
`    "Address":"127.0.0.1",	//redis地址\n`
`    "Port":"6379",			//redis端口\n`
`    "Password":"",			//密码\n`
`    "DBIndex":4				//redis数据库index\n`
`}`
### 处理规则配置：
`"Rule":{\n`
`	"Pattern":["*"],				//key的匹配方式，支持多个匹配规则\n`
`	"ScanCount":1000,			//scan函数每次扫描的key的数量\n`
`	"Method":"expire",			//处理方法（包括删除delete、设置失效时间expire）\n`
`	"ExpireTime":259200			//设置失效时间（单位：秒）\n`
`}`