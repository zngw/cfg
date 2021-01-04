# 游戏配置转换工具

　　在游戏开发过程中，肯定会需要用到一些配置，如果联网游戏还会涉及到服务器和客户端。相同的功能不同的属性可能会在服务器和客户端共用，有的可能只有服务器使用，有的只有客户端来使用。在配置表的时候如果服务器和客户端分开来配容易出错。可以把数据配在一个表里，然后用工具拆分是服务器的还是客户端的数据。  

　　此工具是使用Go语言开发的可以实现配置一表xlsx表，按标志分别输入为客户端或服务器的表，输出格式有 **json**、**ts**、**js**、**mongodb**、**csv**的小工具。

# 安装

### 源码安装

安装好go开发环境，直接编译

```bash
git clone https://github.com/zngw/cfg.git
```

### 直接下载编译好的可执行文件

下载地址：[cfg](https://github.com/zngw/cfg/releases)

# 配置

### 直接运行

运行后会装**./excel**目录下的xlsx文件，服务器转成 **mongodb**,客户端转成 **json**

### 配置说明

#### 文件配置 

配置文件为json格式。可以用命令行参数 -c 加载配置文件，默认为同目录下的 **./conf.json**

```json
{
    "path":"./excel",
    "pre":"Table",
    "type": "all",
    "client":["json","js","ts"],
    "server":["mdb","csv"]
}
```

|  参数  | 说明             |
|-------|-----------------|
| path  | Excel所在目录     | 
| pre   | 转换表前缀        |
| type  | 输出类型          |
| client| 客户端输出文件类型  |
| server| 服务器输出文件类型  |

#### 命令行配置

```bash
cfg -c ./cfg.json -pre Table -p ./excel -t all -cli json|ts -ser mdb
```

|  参数  | 说明             |
|-------|-----------------|
| c     | 配置文件          | 
| path  | Excel所在目录     | 
| pre   | 转换表前缀        |
| type  | 输出类型          |
| client| 客户端输出文件类型  |
| server| 服务器输出文件类型  |

如果同时配置了文件和命令行参数，则命令行参数覆盖配置文件中的

# 表格配置

一般的配置可分为二种，一种是属性值，就是一个key对应一个值的；另一种就是数组性质的表格。

## 配置说明

#### 配置表中表明该配置是属性哪端

* Null  : 不加入配置
* Common: 共同配置
* Client: 客户端配置
* Server: 服务器配置

#### 编译输入类型

* all: 编译服务和客户端
* cli: 编译客户端
* ser: 编译服务器

#### 输入文件类型

* json:输出Json格式
* js:输出js格式
* ts:输出ts格式
* mdb:输出Mongodb的js格式
* csv:输出csv

#### 数据类型说明

* BOOL: bool类型
* INT: int类型
* LONG: long类型
* FLOAT: 浮点类型
* STRING: 字符串
* OBJ: json数据
* ARRAY: json数组

## 解析数组表

表格前四行为固定格式。  
第一行标识为变量为服务器、客户端使用  
第二行为变量名  
第三行为变量数据类型  
第四行为变量描述，不记录转后的数据中  

|Common|Common|Common|Server|Server|Client|
|------|------|------|------|------|------|
|id	|type|	name|	attr|	rand|	desc|
|INT|	INT|	STRING|	OBJ|	ARRAY|	STRING|
|道具ID|	类型	|名字	|固定属性|	随机属性|	描述|
|1	|1	|刀	|{"atk":100}	|[{"weight":80,"attr":"atk","min":10,"max":100},{"weight":20,"attr":"atk","min":100,"max":500}]|	这是一把好刀|
|2	|1	|剑	|{"atk":100}	|[{"weight":80,"attr":"atk","min":10,"max":100},{"weight":20,"attr":"atk","min":100,"max":500}]|	这是一把好剑|
|3	|2	|衣服|{"def":100,"hp":55}|	[{"weight":50,"attr":"def","min":10,"max":100},{"weight":50,"attr":"hp","min":10,"max":100}]|	这是衣服|

### 输出

#### 客户端js

```js
var TableItem = [{"desc":"这是一把好刀","id":1,"name":"刀","type":1},{"desc":"这是一把好剑","id":2,"name":"剑","type":1},{"desc":"这是衣服","id":3,"name":"衣服","type":2}]
modules.export = TableItem
```

#### 客户端ts

```ts
export let TableItem = [{"desc":"这是一把好刀","id":1,"name":"刀","type":1},{"desc":"这是一把好剑","id":2,"name":"剑","type":1},{"desc":"这是衣服","id":3,"name":"衣服","type":2}]
```

#### 客户端json

```json
[{"desc":"这是一把好刀","id":1,"name":"刀","type":1},{"desc":"这是一把好剑","id":2,"name":"剑","type":1},{"desc":"这是衣服","id":3,"name":"衣服","type":2}]
```

#### 服务器csv

```csv
id,type,name,attr,rand
1,1,刀,100,atk=100=10=80|atk=500=100=20
2,1,剑,100,atk=100=10=80|atk=500=100=20
3,2,衣服,100=55,def=100=10=50|hp=100=10=50
```

#### 服务器mdb

```mdb
db.getCollection("TableItem").drop();db.createCollection("TableItem");db.getCollection("TableItem").insert({id:NumberInt("1"),type:NumberInt("1"),name:"刀",attr:{atk:100},rand:[{attr:"atk",min:10,max:100,weight:80},{attr:"atk",min:100,max:500,weight:20}],_id:"1"});db.getCollection("TableItem").insert({id:NumberInt("2"),type:NumberInt("1"),name:"剑",attr:{atk:100},rand:[{weight:80,attr:"atk",min:10,max:100},{weight:20,attr:"atk",min:100,max:500}],_id:"2"});db.getCollection("TableItem").insert({type:NumberInt("2"),name:"衣服",attr:{def:100,hp:55},rand:[{weight:50,attr:"def",min:10,max:100},{min:10,max:100,weight:50,attr:"hp"}],_id:"3",id:NumberInt("3")});
```

## 解析属性表

第一行第一列为空，用来区分表格类型。  
第一行为描述，为配置固定格式，不会记录转换后的数据中。  
第二行开始为正式配置  

|	   | 变量名      |数据类型 |	值	                | 描述  |
|------|------------|-------|-----------------------|------|
|Common|down_url    |STRING |https://zengwu.com.cn  |下载地址|
|Server|init_coin   |INT    |10000                  |初始金币|
|Client|default_head|STRING |icon_0                 |默认头像|

### 输出

#### 客户端js

```js
var TableConst = {"default_head":"icon_0","down_url":"https://zengwu.com.cn"}
modules.export = TableConst
```

#### 客户端ts

```ts
export let TableConst = {"default_head":"icon_0","down_url":"https://zengwu.com.cn"}
```

#### 客户端json

```json
{"default_head":"icon_0","down_url":"https://zengwu.com.cn"}
```

#### 服务器csv

```csv
down_url,init_coin
https://zengwu.com.cn,10000
```

#### 服务器mdb

```mdb
db.getCollection("TableConst").drop();db.createCollection("TableConst");db.getCollection("TableConst").insert({down_url:"https://zengwu.com.cn",init_coin:NumberInt("10000")});
```
