package signed

import (
	"fmt"
	"log"
	"testing"
)

func TestBsi(t *testing.T) {
	mp := New(64)
	{
		xs := []int64{
			10, 3, -7, 9, 0, 1, 9, -8, 2, -1, 12, -35435, 6545654, 2332, 2,
		}
		for i, x := range xs {
			if err := mp.Set(uint32(i), x); err != nil {
				log.Fatal(err)
			}
		}
		{
			fmt.Printf("\tlist\n")
			for i, x := range xs {
				v, ok := mp.Get(uint32(i))
				if ok {
					fmt.Printf("\t\t[%v] = %v, %v\n", i, x, v)
				}
			}
		}
		{
			mq, err := mp.Eq(int64(3))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(= 3) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Lt(int64(-1))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(< -1) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Le(int64(-1))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(<= -1) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Gt(int64(-1))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(> -1) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Ge(int64(-1))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(>= -1) -> %v\n", mq.ToArray())
		}
	}
	data, err := mp.Show()
	if err != nil {
		log.Fatal(err)
	}
	mq := New(0)
	if err := mq.Read(data); err != nil {
		log.Fatal(err)
	}
	{
		xs := []int64{
			10, 3, -7, 9, 0, 1, 9, -8, 2, -1, 12, -35435, 6545654, 2332, 2,
		}
		{
			fmt.Printf("\tlist\n")
			for i, x := range xs {
				v, ok := mq.Get(uint32(i))
				if ok {
					fmt.Printf("\t\t[%v] = %v, %v\n", i, x, v)
				}
			}
		}
		{
			m, err := mq.Eq(int64(-35435))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(= -35435) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Lt(int64(-7))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(< -7) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Le(int64(-8))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(<= -8) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Gt(int64(0))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(> 0) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Ge(int64(0))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(>= 0) -> %v\n", m.ToArray())
		}
	}
}
