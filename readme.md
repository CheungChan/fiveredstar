## golang小工具

### 目前支持的操作
#### logger
```go
package main

import (
    "github.com/cheungchan/fiveredstar/logger"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "fmt"
)

var gLogger zerolog.Logger

// 脚本日志初始化
gLogger = logger.New(path, "myProject.log", true, 20, 10*1024*1024, 1)
// 打印日志
logger.Logger.Info().Msgf("执行%s,超时时间%ds", cmdStr, timeout)

// 在gin里使用日志
r := gin.Default()
r.Use(GinMiddleware("test_gin"))
```
#### cmd
```go
import (
"github.com/cheungchan/fiveredstar/cmd"
)

// 执行持续输出的命令  命令的标准输出和标准错误会重定向的out这个channel里，通过cmd.HandleOutputChannel来按行迭代channel里面的line
out := cmd.GetShellOutput(ctx, cmdStr, false)
cmd.HandleOutputChannel(out, false, func (line string) {
    fmt.Println(line)
})

// 执行很快执行完的命令
s, err:= cmd.GetShellOutputOnce("ls -l", false)

// 判断是否是开发者，可以传入一系列开发者白名单，如果当前操作系统的用户名字了，则是dev模式
cmd.IsDev("chenzhang","liuhan")

// 获取字符串的md5值
cmd.GetMd5("balabala")
```
#### io
```go
import (
"github.com/cheungchan/fiveredstar/io"

// 判断是否是文件夹
b := io.IsDir("/root")

// 判断文件是否存在
b = io.FileIsExist("/root/a.txt")

// 创建文件夹
io.MakeDir("/root/tmp")

// 拷贝文件
io.CopyFile("/root/a.txt","/root/b.txt")

// 递归拷贝文件夹
io.CopyDir("/root/a/","/root/b/")

// 列出当前文件夹下的所有文件，不递归
io.ListDir("/root")

// 列出当前文件夹下的所有文件和子文件夹下的所有文件
io.WalkDir("/root")

// 删除文件
io.RemoveFile("/root/a.txt")

// 删除文件夹
io.RemoveDirAll("/root/a/")
)

```