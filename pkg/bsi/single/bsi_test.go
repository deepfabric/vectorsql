package single

import (
	"fmt"
	"log"
	"testing"
)

func TestBsi(t *testing.T) {
	mp := New()
	{
		xs := []float32{
			1.9e4, 1.8e5, 1.9e6, 1.4e8, 1.3e10, 1.5e11,
		}
		for i, x := range xs {
			{
				fmt.Printf("set [%v] = %v\n", i, x)
			}
			if err := mp.Set(uint32(i), x); err != nil {
				log.Fatal(err)
			}
		}
		{
			m0, err := mp.Lt(float32(1.3e10))
			if err != nil {
				log.Fatal(err)
			}
			m1, err := mp.Gt(float32(1.8e5))
			if err != nil {
				log.Fatal(err)
			}
			m0.And(m1)
			fmt.Printf("%v\n", m0.ToArray())
		}
		/*
			{
				mq, err := mp.Lt(float32(3.0))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\t(< 3.0) -> %v\n", mq.ToArray())
			}
			{
				mq, err := mp.Gt(float32(3.0))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\t(> 3.0) -> %v\n", mq.ToArray())
			}
		*/
	}
}
