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
