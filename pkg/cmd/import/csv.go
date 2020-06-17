package main

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/deepfabric/vectorsql/pkg/request"
)

func loadCsv(id uint32, lines [][]string, name string) ([][]interface{}, []map[string]interface{}) {
	args := make([][]interface{}, len(lines))
	ts := make([]map[string]interface{}, len(lines))
	for i, j := 0, len(lines); i < j; i++ {
		ts[i] = make(map[string]interface{})
	}
	var wg sync.WaitGroup
	step := len(lines) / 8
	for x := 0; x <= 8; x++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for i, j := idx*step, idx*step+step; i < j; i++ {
				if i >= len(lines) {
					break
				}
				line := lines[i]
				arg := make([]interface{}, 0, len(line))
				{
					ts[i]["seq"] = uint32(i) + id
					arg = append(arg, uint32(i)+id)
				}
				{
					sex, err := strconv.ParseUint(line[0], 0, 8)
					if err != nil {
						log.Fatalf("Failed to parse file '%s': %s is not uint8", name, line[1])
					}
					ts[i]["sex"] = uint8(sex)
					arg = append(arg, uint8(sex))
				}
				{
					age, err := strconv.ParseUint(line[1], 0, 8)
					if err != nil {
						log.Fatalf("Failed to parse file '%s': %s is not uint8", name, line[2])
					}
					ts[i]["age"] = uint8(age)
					arg = append(arg, uint8(age))
				}
				{
					line[2] = strings.ReplaceAll(line[2], " ", "")
					line[2] = strings.ReplaceAll(line[2], "\t", "")
					ts[i]["area"] = string(line[2])
					arg = append(arg, string(line[2]))
				}
				{
					{
						line[3] = strings.ReplaceAll(line[3], "[", "")
						line[3] = strings.ReplaceAll(line[3], "]", "")
						line[3] = strings.ReplaceAll(line[3], " ", "")
						line[3] = strings.ReplaceAll(line[3], "\t", "")
					}
					fs := make(map[string]*request.Part)
					xs := strings.Split(line[3], ",")
					for _, x := range xs {
						data, err := ioutil.ReadFile(x)
						if err != nil {
							log.Fatalf("Failed to read image file '%s': %v\n", x, err)
						}
						fs[x] = &request.Part{Typ: "image", Data: data}
					}
					_, err := getVector(fs)
					if err != nil {
						log.Fatalf("Failed to get vector: %v\n", err)
					}
				}
				args[i] = arg
			}
		}(x)
	}
	wg.Wait()
	/*
		for i, line := range lines {
			arg := make([]interface{}, 0, len(line))
			{
				ts[i]["seq"] = uint32(i) + id
				arg = append(arg, uint32(i)+id)
			}
			{
				sex, err := strconv.ParseUint(line[0], 0, 8)
				if err != nil {
					log.Fatalf("Failed to parse file '%s': %s is not uint8", name, line[1])
				}
				ts[i]["sex"] = uint8(sex)
				arg = append(arg, uint8(sex))
			}
			{
				age, err := strconv.ParseUint(line[1], 0, 8)
				if err != nil {
					log.Fatalf("Failed to parse file '%s': %s is not uint8", name, line[2])
				}
				ts[i]["age"] = uint8(age)
				arg = append(arg, uint8(age))
			}
			{
				line[2] = strings.ReplaceAll(line[2], " ", "")
				line[2] = strings.ReplaceAll(line[2], "\t", "")
				ts[i]["area"] = string(line[2])
				arg = append(arg, string(line[2]))
			}
			{
				{
					line[3] = strings.ReplaceAll(line[3], "[", "")
					line[3] = strings.ReplaceAll(line[3], "]", "")
					line[3] = strings.ReplaceAll(line[3], " ", "")
					line[3] = strings.ReplaceAll(line[3], "\t", "")
				}
				fs := make(map[string]*request.Part)
				xs := strings.Split(line[3], ",")
				for _, x := range xs {
					data, err := ioutil.ReadFile(x)
					if err != nil {
						log.Fatalf("Failed to read image file '%s': %v\n", x, err)
					}
					fs[x] = &request.Part{Typ: "image", Data: data}
				}
				_, err := getVector(fs)
				if err != nil {
					log.Fatalf("Failed to get vector: %v\n", err)
				}
			}
			args[i] = arg
		}
	*/
	return args, ts
}

func readFile(name string) [][]string {
	fp, err := os.Open(name)
	if err != nil {
		log.Fatalf("Failed to open file '%s': %v\n", name, err)
	}
	defer fp.Close()
	r := csv.NewReader(bufio.NewReader(fp))
	r.TrimLeadingSpace = true
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read file '%s': %v\n", name, err)
	}
	return lines
}
