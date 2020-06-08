package unsigned

import (
	"fmt"
	"log"
	"testing"
)

func TestUbsi(t *testing.T) {
	mp := New(64)
	{
		xs := []uint64{
			10, 3, 7, 9, 0, 1, 9, 8, 2, 12, 35435, 6545654, 2332, 2,
		}
		for i, x := range xs {
			if err := mp.Set(uint32(i), x); err != nil {
				log.Fatal(err)
			}
		}
		{
			if err := mp.Del(uint32(9)); err != nil {
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
			mq, err := mp.Eq(uint64(3))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(= 3) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Lt(uint64(10))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(< 10) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Le(uint64(10))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(<= 10) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Gt(uint64(7))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(> 7) -> %v\n", mq.ToArray())
		}
		{
			mq, err := mp.Ge(uint64(7))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(>= 7) -> %v\n", mq.ToArray())
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
		xs := []uint64{
			10, 3, 7, 9, 0, 1, 9, 8, 2, 12, 35435, 6545654, 2332, 2,
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
			m, err := mq.Eq(uint64(3))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(= 3) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Lt(uint64(10))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(< 10) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Le(uint64(10))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(<= 10) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Gt(uint64(7))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(> 7) -> %v\n", m.ToArray())
		}
		{
			m, err := mq.Ge(uint64(7))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\t(>= 7) -> %v\n", m.ToArray())
		}
	}
}
