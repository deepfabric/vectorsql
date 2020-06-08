package bm

import (
	"fmt"
	"testing"
)

func TestBm(t *testing.T) {
	{
		var bs []Bm

		bs = append(bs, Bm{false, nil, "bm0"})
		bs = append(bs, Bm{true, nil, "bm1"})
		bs = append(bs, Bm{false, nil, "bm2"})
		bs = append(bs, Bm{false, nil, "bm3"})
		fmt.Printf("%v\n", Gen(bs))
	}
	{
		var bs []Bm

		bs = append(bs, Bm{false, nil, "bm0"})
		{
			var bt []Bm

			bt = append(bt, Bm{true, nil, "bm1"})
			bt = append(bt, Bm{false, nil, "bm2"})
			bs = append(bs, Bm{false, bt, ""})
		}
		bs = append(bs, Bm{false, nil, "bm3"})
		fmt.Printf("%v\n", Gen(bs))
	}

}
