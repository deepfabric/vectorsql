package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/RoaringBitmap/roaring"
	"github.com/deepfabric/vectorsql/pkg/sql/client"
)

func main() {
	cli, err := client.New("tcp://172.19.0.17:9000?username=cdp_user&password=infinivision2019")
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	if ts, err := cli.Query(os.Args[1]); err != nil {
		log.Fatal(err)
	} else {
		for i, t := range ts {
			n := len(t) - 1
			ioutil.WriteFile(fmt.Sprintf("%v.jpg", i), []byte(t[n]), os.FileMode(0664))
		}
	}
	fmt.Printf("process: %v\n", time.Now().Sub(t))
}

func UnmarshalMap(data []byte) (*roaring.Bitmap, error) {
	switch data[0] {
	case 0:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		mp := roaring.New()
		mp.AddMany(decodeVector(data[n:]))
		return mp, nil
	case 1:
		data = data[1:]
		_, n := binary.Uvarint(data)
		if n < 0 {
			return nil, errors.New("overflow")
		}
		data = data[n:]
		mp := roaring.New()
		if err := mp.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return mp, nil
	}
	return nil, nil
}

func decodeVector(v []byte) []uint32 {
	hp := *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	hp.Len /= 4
	hp.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&hp))
}
