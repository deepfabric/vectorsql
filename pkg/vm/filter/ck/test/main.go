package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/deepfabric/vectorsql/pkg/sql/client"
	"github.com/deepfabric/vectorsql/pkg/vm/extend"
	"github.com/deepfabric/vectorsql/pkg/vm/extend/overload"
	"github.com/deepfabric/vectorsql/pkg/vm/filter/ck"
	"github.com/deepfabric/vectorsql/pkg/vm/value"
)

func main() {
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := time.Now()
			mp, err := ck.New(cli, genCondition()).Bitmap()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("process: %v\n", time.Now().Sub(t))
			fmt.Printf("\t%v\n", mp.ToArray())
		}()
	}
	wg.Wait()
}

func genCondition() map[string]extend.Extend {
	mp := make(map[string]extend.Extend)
	{
		mp["people"] = &extend.BinaryExtend{
			Op:    overload.EQ,
			Left:  &extend.Attribute{"gender"},
			Right: value.NewInt(0),
		}
	}
	{
		mp["people_events"] = &extend.BinaryExtend{
			Op:    overload.EQ,
			Left:  &extend.Attribute{"area"},
			Right: value.NewString("shanghai"),
		}
	}
	return mp
}
