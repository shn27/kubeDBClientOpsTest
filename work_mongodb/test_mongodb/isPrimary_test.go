package test_mongodb

import (
	"fmt"
	work "github.com/shn27/Test/work_mongodb"
	"testing"
)

func Test_isPrimary(t *testing.T) {
	primary, err := work.IsPrimary()
	if err != nil {
		return
	}
	fmt.Println(primary)
}
