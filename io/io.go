package io

import (
    "github.com/pkg/errors"
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func IsDir(name string) bool {
    if info, err := os.Stat(name); err == nil {
        return info.IsDir()
    }
    return false
}

func FileIsExisted(filename string) bool {
    existed := true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        existed = false
    }
    return existed
}
func MakeDir(dir string) error {
    if !FileIsExisted(dir) {
        if err := os.MkdirAll(dir, 0777); err != nil { //os.ModePerm
            return errors.Wrap(err, "MakeDir failed")
        }
    }
    return nil
}

// CopyFile 使用io.Copy
func CopyFile(src, des string) (written int64, err error) {
    srcFile, err := os.Open(src)
    if err != nil {
        return 0, err
    }
    defer func(srcFile *os.File) {
        err = srcFile.Close()
    }(srcFile)

    //获取源文件的权限
    fi, _ := srcFile.Stat()
    perm := fi.Mode()

    //desFile, err := os.Create(des)  //无法复制源文件的所有权限
    desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm) //复制源文件的所有权限
    if err != nil {
        return 0, err
    }
    defer func(desFile *os.File) {
        err = desFile.Close()
    }(desFile)

    return io.Copy(desFile, srcFile)
}

func CopyDir(srcPath, desPath string) error {
    //检查目录是否正确
    if srcInfo, err := os.Stat(srcPath); err != nil {
        return err
    } else {
        if !srcInfo.IsDir() {
            return errors.New("源路径不是一个正确的目录！")
        }
    }

    if desInfo, err := os.Stat(desPath); err != nil {
        return err
    } else {
        if !desInfo.IsDir() {
            return errors.New("目标路径不是一个正确的目录！")
        }
    }

    if strings.TrimSpace(srcPath) == strings.TrimSpace(desPath) {
        return errors.New("源路径与目标路径不能相同！")
    }

    err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
        if f == nil {
            return err
        }

        //复制目录是将源目录中的子目录复制到目标路径中，不包含源目录本身
        if path == srcPath {
            return nil
        }

        //生成新路径
        destNewPath := strings.Replace(path, srcPath, desPath, -1)

        if !f.IsDir() {
            _, err = CopyFile(path, destNewPath)
            if err != nil {
                return err
            }
        } else {
            if !FileIsExisted(destNewPath) {
                return MakeDir(destNewPath)
            }
        }

        return nil
    })

    return err
}

// ListDir  获取指定路径下的所有文件，只搜索当前路径，不进入下一级目录，可匹配后缀过滤（suffix为空则不过滤）*/
func ListDir(dir, suffix string) (files []string, err error) {
    files = []string{}

    _dir, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    suffix = strings.ToLower(suffix) //匹配后缀

    for _, _file := range _dir {
        if _file.IsDir() {
            continue //忽略目录
        }
        if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(_file.Name()), suffix) {
            //文件后缀匹配
            files = append(files, path.Join(dir, _file.Name()))
        }
    }

    return files, nil
}

// WalkDir  获取指定路径下以及所有子目录下的所有文件，可匹配后缀过滤（suffix为空则不过滤）*/
func WalkDir(dir, suffix string) (files []string, err error) {
    files = []string{}

    err = filepath.Walk(dir, func(name string, fi os.FileInfo, err error) error {
        if fi.IsDir() {
            //忽略目录
            return nil
        }

        if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
            //文件后缀匹配
            files = append(files, name)
        }

        return nil
    })

    return files, err
}
func RemoveFile(filename string) error {
    return os.Remove(filename)
}
func RemoveDirAll(dir string) error {
    return os.RemoveAll(dir)
}
