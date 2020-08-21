package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// time.Time formats.
const (
	TimestampOutputFormat = "2006-01-02 15:04:05"
)

const charset = "0123456789"

var sexs []string
var cities []string

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func init() {
	sexs = []string{"男", "女"}
	cities = []string{"河北", "山西", "内蒙古", "辽宁", "吉林", "黑龙江", "江苏", "浙江", "安徽",
		"福建", "江西", "山东", "河南", "湖北", "湖南", "广东", "上海", "北京", "天津", "重庆"}
}

func main() {
	var uid, pid uint64

	fp, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	ts, err := csv.NewReader(fp).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		for _, t := range ts {
			if len(t) != 513 {
				continue
			}
			vec, err := getVector(t[1:])
			if err != nil {
				continue
			}
			fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%s\n",
				uid, (pid | uid<<34), t[0], String(17), sexs[uid%2], cities[int(uid)%len(cities)], time.Now().Add(time.Duration(pid|uid<<34)).Format(TimestampOutputFormat), vec)
			uid++
			pid++
		}
	}
}

func getVector(t []string) (string, error) {
	v := make([]float32, 512)
	for i, x := range t {
		f, err := strconv.ParseFloat(x, 32)
		if err != nil {
			return "", err
		}
		v[i] = float32(f)
	}
	data, err := json.Marshal(v)
	return string(data), err
}

/*
func main() {
	var uid, pid uint64

	fs, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		for _, f := range fs {
			fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n",
				uid, (pid | uid<<34), f.Name(), String(17), sexs[uid%2], cities[int(uid)%len(cities)], time.Now().Add(time.Duration(pid|uid<<34)).Format(TimestampOutputFormat))
			uid++
			pid++
		}
	}
}
*/
