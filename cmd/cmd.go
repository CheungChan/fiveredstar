package cmd

import (
    "bufio"
    "context"
    "crypto/md5"
    "encoding/hex"
    "github.com/cheungchan/fiveredstar/io"
    "github.com/cheungchan/fiveredstar/logger"
    "os/exec"
    "os/user"
    "strings"
    "sync"
)

func GetShellOutput(ctx context.Context, cmdStr string, debug bool) chan string {
    // 获取shell输出的标准输出和标准错误
    outputChanel := make(chan string, 10)
    go runShellAsync(ctx, cmdStr, outputChanel, debug)
    return outputChanel

}
func GetShellOutputOnce(cmdStr string, debug bool) (string, error) {
    // 同步的获取shell输出
    c := exec.Command("bash", "-c", cmdStr)
    if debug {
        logger.Logger.Debug().Msgf("执行命令 %s", cmdStr)

    }
    // 此处是windows版本
    // c := exec.Command("cmd", "/C", cmd)
    output, err := c.CombinedOutput()
    return string(output), err
}
func HandleOutputChannel(ch chan string, debug bool, op func(line string)) {
    if debug {
        logger.Logger.Debug().Msgf("处理output channel")
    }
    for line := range ch {
        line = strings.Trim(line, "\n")
        line = strings.TrimSpace(line)
        if line != "" {
            op(line)
        }
    }
    if debug {
        logger.Logger.Debug().Msgf("处理output channel完成")
    }
}
func runShellAsync(ctx context.Context, cmdStr string, outputChanel chan string, debug bool) {
    defer func() {
        if debug {
            logger.Logger.Debug().Msgf("关闭outputChanel,%s", cmdStr)
        }
        close(outputChanel)
    }()
    logger.Logger.Info().Msgf("执行%s", cmdStr)
    cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()

    err := cmd.Start()
    if err != nil {
        logger.Logger.Error().Msgf("cmd start失败 %s, err: %+v", cmdStr, err)
        return
    }
    readerOut := bufio.NewScanner(stdout)
    readerErr := bufio.NewScanner(stderr)
    // 必须要等到两个goroutine把标准输出和标准错误消费完再return
    wg := sync.WaitGroup{}
    wg.Add(2)
    //实时循环读取输出流中的一行内容
    go func() {
        defer func() {
            if err := recover(); err != nil {
                logger.Logger.Error().Msgf("cmd 失败 %s, err: %+v", cmdStr, err)
            }
            wg.Done()
        }()
        if debug {
            logger.Logger.Debug().Msgf("标准输出重定向到outputChanel,%s", cmdStr)
        }
        for readerOut.Scan() {
            line := readerOut.Text()
            outputChanel <- line

        }
        if debug {
            logger.Logger.Debug().Msgf("标准输出重定向到outputChanel完成,%s", cmdStr)
        }
    }()
    go func() {
        defer func() {
            if err := recover(); err != nil {
                logger.Logger.Error().Msgf("cmd 失败 %s, err: %+v", cmdStr, err)
            }
            wg.Done()
        }()
        if debug {
            logger.Logger.Debug().Msgf("标准错误重定向到outputChanel,%s", cmd)
        }
        for readerErr.Scan() {
            line := readerErr.Text()
            outputChanel <- line
        }
        if debug {
            logger.Logger.Debug().Msgf("标准错误重定向到outputChanel完成,%s", cmdStr)
        }
    }()
    if debug {
        logger.Logger.Debug().Msgf("标准输出，标准错误启动goroutine等待,%s", cmdStr)
    }
    _ = cmd.Wait()
    wg.Wait()
}
func FileExists(path string) bool {
    return io.FileIsExisted(path)
}

func IsDev(whiteUserName ...string) (r bool) {
    // 是否是开发环境，username传入白名单
    current, _ := user.Current()
    defer func() {
        logger.Logger.Debug().Msgf("当前用户：%s,IsDev:%t", current.Username, r)
    }()
    for _, u := range whiteUserName {
        if u == current.Username {
            r = true
            return
        }
    }
    r = false
    return
}
func GetMd5(str string) string {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}
