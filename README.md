# SDSL Interpreter

## 总览

SDSL(Simple Domain Specific Language) Interpreter是一个简单领域定义语言应答服务脚本解释器。

您可以使用下面的SDSL语法设计一份脚本，并保存在script.txt文件中，然后运行程序。您的用户将通过浏览器或其他支持WebSocket的方式与本解释器交互。

## SDSL语法

SDSL的每条语句都形如`[关键字] [符号]`，且`[符号]`部分可以是0-2个。`[关键字]`指明了该语句的含义，`[符号]`部分用于构成完整的语句含义。

### 关键字用法

#### Step

用法：`Step [StepId]`

定义一个步骤。

#### Speak

用法：`Speak [Expression]`

将`Expression`的内容告诉给用户。

#### Listen

用法：`Listen [beginTimer] [endTimer]`

听取用户的命令，要求用户在`beginTimer`时间内开始说，说话的时间不超过`endTimer`。本程序将这部分抽象为要求用户在`beginTimer`时间内发送消息，消息长度不超过`endTimer`字节。

#### Branch

用法：`Branch [Answer] [StepId]`

设置当用户的消息为`[Answer]`时，下一步转移到`[StepId]`对应的步骤。

#### Silence

用法：`Silence [StepId]`

设置当无法听到用户的消息时，下一步转移到`[StepId]`对应的步骤。

#### Default

用法：`Default [StepId]`

设置当用户的消息为不在`Branch`设置转移的消息列表中时，下一步转移到`[StepId]`对应的步骤。

#### Exit

用法：`Exit`

设置当前步骤为结束步骤，**一个脚本只能有一个结束步骤**。

### 变量

变量为一个`$`开头的字符串，用于将`Expression`中对应位置替换为与用户有关的字符串。

**注意：** 变量仅能在`Expression`中使用，在其他位置使用将被视为一般的字符串。

### 注释

输入`#`以编写注释，`#`后到其所在行末尾均被视为注释。

## 常量、函数与数据结构

```go
const (
	// DEFAULT 代表初始化控制块的默认模式（填充随机浮点数）
	DEFAULT = 0
	// FILE 代表使用文件初始化控制块
	FILE = 1
	// SQL 代表使用数据库初始化控制块
	SQL = 2
)
```

#### func  ControlBlockInit

```go
func ControlBlockInit(controlBlock *UserControlBlock, mode int, address string) error
```

ControlBlockInit 初始化用户控制块

#### func  CreateLog

```go
func CreateLog() error
```

CreateLog 创建日志文件

#### func  HandleMessage

```go
func HandleMessage(conn *websocket.Conn, controlBlock *UserControlBlock)
```

HandleMessage 处理前端命令

#### func  HandleSpeak

```go
func HandleSpeak(controlBlock *UserControlBlock) []byte
```

HandleSpeak 完成脚本中Speak语句的功能

#### func  Handler

```go
func Handler(w http.ResponseWriter, r *http.Request)
```

Handler 处理前后端通信事务

#### func  Log

```go
func Log(message string)
```

Log 写入日志

#### func  ParseFile

```go
func ParseFile(fileName string) error
```

ParseFile 处理脚本文件

#### func  ParseLine

```go
func ParseLine(line []string) error
```

ParseLine 处理脚本文件的每一行

#### func  ProcessBranch

```go
func ProcessBranch(answer, nextStepId string)
```

ProcessBranch 处理Branch分支

#### func  ProcessDefault

```go
func ProcessDefault(nextStepId string)
```

ProcessDefault 处理Default分支

#### func  ProcessListen

```go
func ProcessListen(start, end string) error
```

ProcessListen 处理Listen分支

#### func  ProcessSilence

```go
func ProcessSilence(nextStepId string)
```

ProcessSilence 处理Silence分支

#### func  ProcessSpeak

```go
func ProcessSpeak(token []string)
```

ProcessSpeak 处理Speak分支

#### func  ProcessStep

```go
func ProcessStep(stepId StepId)
```

ProcessStep 处理Step分支

#### func  ProcessTokens

```go
func ProcessTokens(tokens []string) error
```

ProcessTokens 对每一行的token进行处理

#### func  Server

```go
func Server(port int) error
```

Server 提供Web服务

#### type Answer

```go
type Answer string
```

Answer 表示用户的命令

#### type Expression

```go
type Expression struct {
}
```

Expression 表示Speak的表达式

#### func  ProcessExpression

```go
func ProcessExpression(tokens []string) Expression
```

ProcessExpression 将token转为表达式返回给Speak

#### type Listen

```go
type Listen struct {
}
```

Listen 是用于控制Listen的变量

#### type Script

```go
type Script struct {
}
```

Script 是脚本的数据结构

#### type Step

```go
type Step struct {
}
```

Step 是步骤的数据结构

#### type StepId

```go
type StepId string
```

StepId 表示步骤ID

#### func  HandleListen

```go
func HandleListen(conn *websocket.Conn, controlBlock *UserControlBlock) StepId
```

HandleListen 完成脚本中Listen语句的功能

#### type UserControlBlock

```go
type UserControlBlock struct {
}
```

UserControlBlock 是控制块的数据结构

#### type VarName

```go
type VarName string
```

VarName 表示脚本中的变量名

## 使用

编写main.go，下面给出了一个例子。

```go
//main.go
package main

import (
   "main/package/sdsl"
   "os"
)

var port int
var scriptPath string

func main() {
   //初始化
   port = 7777 //设置Web服务端口
   scriptPath = "script.txt" //将您编写的脚本放入这个文件
   err := sdsl.CreateLog() //创建日志
   if err != nil {
      os.Exit(1)
   }
   err = sdsl.ParseFile(scriptPath) //处理脚本文件
   if err != nil {
      os.Exit(1)
   }

   //运行服务
   err = sdsl.Server(port)
   if err != nil {
      os.Exit(1)
   }
}
```

访问localhost:[port]测试程序功能是否正常。如需使用其他支持Websocket的方式，请连接至ws://localhost:[port]/ws来测试。
