package zip

import (
    logger2 "github.com/cheungchan/fiveredstar/logger"
    "testing"
)

func TestZip(t *testing.T) {
    _ = logger2.New("logs","zip_test.log",true,3,30*1024*1024,1)
    t.Log("开始测试带文件夹的压缩")
    err := Zip("./test_zip_dir/raw_files.zip", "test_raw_files", true,false)
    if err != nil {
        t.Error(err)
    }
    t.Log("开始测试带文件夹的解压")
    err = UnZip("./test_unzip_dir", "./test_zip_dir/raw_files.zip", true)
    if err != nil {
        t.Error(err)
    }
    t.Log("开始测试不带文件夹的压缩")
    err = Zip("./test_zip_flat/3.zip", "test_raw_files/3", true,true)
    if err != nil {
        t.Error(err)
    }
    t.Log("开始测试不带文件夹的解压")
    err = UnZip("./test_unzip_flat", "./test_zip_flat/3.zip", true)
    if err != nil {
        t.Error(err)
    }
}
