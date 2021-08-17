package cmd

import (
	"context"
	"fmt"
	"testing"
)

func TestGetShellOutput(t *testing.T) {
	ch := GetShellOutput(context.TODO(),"tail /tmp/adobegc.log",true)
	HandleOutputChannel(ch,true, func(line string) {
		fmt.Println(line)
	})
}

func TestGetShellOutputOnce(t *testing.T) {
	out, err := GetShellOutputOnce(context.TODO(),"ls -l",true)
	if err!=nil{
		t.Error(err)
		return
	}
	t.Log(out)
}