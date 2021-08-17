package zip

import (
    "archive/zip"
    "fmt"
    "github.com/cheungchan/fiveredstar/cmd"
    "github.com/cheungchan/fiveredstar/logger"
    "io"
    "os"
    "path"
    "path/filepath"
    "strings"
)

// Zip 压缩文件
// dst 要压缩为的文件名
// src 要压缩的目录
// debug 是否开启debug模式
// flat  压缩完里面是否要包含文件夹
func Zip(dst, src string, debug bool, flat bool) (err error) {
    dd := path.Dir(dst)
    if !cmd.FileExists(dd) {
        os.MkdirAll(dd, os.ModePerm)
        if debug {
            logger.Logger.Debug().Msgf("%s文件夹不存在，创建文件夹", dd)
        }
    }
    if !cmd.FileExists(src) {
        return fmt.Errorf("%s文件不存在，无法压缩", src)
    }
    // 创建准备写入的文件
    fw, err := os.Create(dst)
    defer fw.Close()
    if err != nil {
        return err
    }

    // 通过 fw 来创建 zip.Write
    zw := zip.NewWriter(fw)
    defer func() {
        // 检测一下是否成功关闭
        if err := zw.Close(); err != nil {
            logger.Logger.Fatal().Msgf("%+v", err)
        }
    }()

    // 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
    return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
        if errBack != nil {
            return errBack
        }

        // 通过文件信息，创建 zip 的文件信息
        fh, err := zip.FileInfoHeader(fi)
        if err != nil {
            return
        }

        // 替换文件信息中的文件名
        fh.Name = strings.TrimPrefix(path, string(filepath.Separator))
        if flat && fi.IsDir(){
            logger.Logger.Debug().Msgf("%s是文件夹，不压缩",fi.Name())
            return nil
        }
        if flat {
            fh.Name = filepath.Base(fh.Name)
        }

        // 这步开始没有加，会发现解压的时候说它不是个目录
        if fi.IsDir() {
            fh.Name += "/"
        }

        // 写入文件信息，并返回一个 Write 结构
        w, err := zw.CreateHeader(fh)
        if err != nil {
            return
        }

        // 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
        // 如目录，也没有数据需要写
        if !fh.Mode().IsRegular() {
            return nil
        }

        // 打开要压缩的文件
        fr, err := os.Open(path)
        defer fr.Close()
        if err != nil {
            return
        }

        // 将打开的文件 Copy 到 w
        n, err := io.Copy(w, fr)
        if err != nil {
            return
        }
        // 输出压缩的内容
        if debug {
            logger.Logger.Debug().Msgf("成功压缩文件： %s, 共写入了 %d 个字符的数据", path, n)
        }

        return nil
    })
}
func UnZip(dst, src string, debug bool) (err error) {
    if !cmd.FileExists(dst) {
        os.MkdirAll(dst, os.ModePerm)
        if debug {
            logger.Logger.Debug().Msgf("%s文件夹不存在，创建文件夹", dst)
        }
    }
    if !cmd.FileExists(src) {
        return fmt.Errorf("%s文件不存在，无法解压缩", src)
    }
    // 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
    // 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
    zr, err := zip.OpenReader(src)
    defer zr.Close()
    if err != nil {
        return
    }

    // 如果解压后不是放在当前目录就按照保存目录去创建目录
    if dst != "" {
        if err := os.MkdirAll(dst, 0755); err != nil {
            return err
        }
    }

    // 遍历 zr ，将文件写入到磁盘
    for _, file := range zr.File {
        path := filepath.Join(dst, file.Name)

        // 如果是目录，就创建目录
        if file.FileInfo().IsDir() {
            if err := os.MkdirAll(path, file.Mode()); err != nil {
                return err
            }
            // 因为是目录，跳过当前循环，因为后面都是文件的处理
            continue
        }

        // 获取到 Reader
        fr, err := file.Open()
        if err != nil {
            return err
        }

        // 创建要写出的文件对应的 Write
        fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
        if err != nil {
            return err
        }

        n, err := io.Copy(fw, fr)
        if err != nil {
            return err
        }

        // 将解压的结果输出
        if debug {
            logger.Logger.Debug().Msgf("成功解压 %s ，共写入了 %d 个字符的数据", path, n)
        }

        // 因为是在循环中，无法使用 defer ，直接放在最后
        // 不过这样也有问题，当出现 err 的时候就不会执行这个了，
        // 可以把它单独放在一个函数中，这里是个实验，就这样了
        fw.Close()
        fr.Close()
    }
    return nil
}
