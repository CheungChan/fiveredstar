package zip

import "testing"

func TestZip(t *testing.T) {
    t.Log("开始压缩")
    err := Zip("./test_zip/raw_files.zip", "test_raw_files", true)
    if err != nil {
        t.Error(err)
    }
    t.Log("开始解压")
    err = UnZip("./test_unzip", "./test_zip/raw_files.zip", true)
    if err != nil {
        t.Error(err)
    }

}
