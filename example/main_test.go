package example_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/opentoys/ledisdb/config"
	"github.com/opentoys/ledisdb/ledis"
)

func TestDB(t *testing.T) {
	// opt.D
	rdb, e := ledis.Open(config.NewConfigDefault())
	if e != nil {
		t.Fatal(e)
	}
	defer rdb.Close()

	db, e := rdb.Select(0)
	if e != nil {
		t.Fatal(e)
	}

	// if e = db.Set([]byte("hello"), []byte("3s")); e != nil {
	// 	t.Fatal(e)
	// }

	if v, e := db.Get([]byte("hello")); e != nil {
		t.Fatal(e)
	} else {
		fmt.Println(string(v))
	}

	<-time.After(time.Second * 4)

	if v, e := db.Get([]byte("hello")); e != nil {
		t.Fatal(e)
	} else {
		fmt.Println(string(v))
	}

	db.Expire([]byte("hello"), 3)
}
