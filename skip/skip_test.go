package skip

import (
	. "github.com/badgerodon/collections"
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	sl := New(func(a,b Any)bool {
		return a.(int) < b.(int)
	})
	for i := 0; i < 20; i++ {
		sl.Insert(i, 1)
	}
	fmt.Println(sl)
	sl.Remove(15)
	fmt.Println(sl)
}
