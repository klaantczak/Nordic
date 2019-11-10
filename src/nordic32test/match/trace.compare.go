package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"strings"
)

func jsonKeys(obj map[string]interface{}) []string {
	keys := []string{}
	for k, _ := range obj {
		keys = append(keys, k)
	}
	return keys
}

func jsonCompareKeys(path string, keys1, keys2 []string) bool {
	allOk := true

	for _, k1 := range keys1 {
		ok := false

		for _, k2 := range keys2 {
			if k1 == k2 {
				ok = true
				break
			}
		}

		if !ok {
			allOk = false
			fmt.Printf("%s: %s is in keys1, but not in keys2\n", path, k1)
		}
	}

	for _, k2 := range keys2 {
		ok := false

		for _, k1 := range keys1 {
			if k1 == k2 {
				ok = true
				break
			}
		}

		if !ok {
			allOk = false
			fmt.Printf("%s: %s is in keys2, but not in keys1\n", path, k2)
		}
	}

	return allOk
}

func jsonDiffObj(path string, obj1 map[string]interface{}, obj2 map[string]interface{}) bool {
	keys1 := jsonKeys(obj1)
	keys2 := jsonKeys(obj2)

	if !jsonCompareKeys(path, keys1, keys2) {
		return false
	}

	for _, k := range keys1 {
		p := fmt.Sprintf("%s[%s]", path, k)

		v1 := obj1[k]
		t1 := reflect.TypeOf(v1)

		v2 := obj2[k]
		t2 := reflect.TypeOf(v2)

		if t1 != t2 {
			fmt.Printf("%s: Type mismatch\n", p)
			return false
		}

		if fmt.Sprintf("%s/%s", path, k) == "{...}/time" {
			time1, ok1 := v1.(float64)
			time2, ok2 := v2.(float64)
			if !ok1 || !ok2 {
				fmt.Printf("%s: %s: %v != %v\n", path, k, v1, v2)
				return false
			}
			d := math.Abs(time1 - time2)
			if d > 1e-12 {
				fmt.Printf("%s: %s: %v !~ %v\n", path, k, time1, time2)
				return false
			}
			continue
		}

		if t1.String() == "map[string]interface {}" {
			if !jsonDiffObj(p, v1.(map[string]interface{}), v2.(map[string]interface{})) {
				return false
			}
		} else {
			switch v := v1.(type) {
			case int:
				if v != v2.(int) {
					fmt.Printf("%s: %s: %v != %v\n", path, k, v1, v2)
					return false
				}
			case float64:
				delta := math.Abs(v - v2.(float64))
				diff := delta / math.Abs(v)
				if diff > 1e-7 {
					fmt.Printf("%s: %s: %v != %v\n", path, k, v1, v2)
					return false
				}
			case string:
				if v != v2.(string) {
					fmt.Printf("%s: %s: %v != %v\n", path, k, v1, v2)
					return false
				}
			}
		}
	}

	return true
}

func main() {
	var file1 = flag.String("file1", "", "first file")
	var file2 = flag.String("file2", "", "second file")

	flag.Parse()

	if *file1 == "" || *file2 == "" {
		fmt.Println("Usage: trace.compare.go -file1 ... -file2 ...")
		return
	}

	fmt.Println(*file1, *file2)

	f1, _ := os.Open(*file1)
	r1 := bufio.NewReader(f1)
	defer f1.Close()

	f2, _ := os.Open(*file2)
	r2 := bufio.NewReader(f2)
	defer f2.Close()

	for {
		var err error

		line1, err1 := r1.ReadBytes('\n')
		line2, err2 := r2.ReadBytes('\n')
		if err1 == io.EOF && err2 == io.EOF {
			log.Print("Files are identical")
			return
		}

		if err1 != nil {
			log.Fatal(err)
		}
		if err2 != nil {
			log.Fatal(err)
		}

		var obj1 map[string]interface{}
		err = json.Unmarshal(line1, &obj1)
		if err != nil {
			log.Fatalf("Error parsing data from trace.567483.jslog: %v", err)
		}

		var obj2 map[string]interface{}
		err = json.Unmarshal(line2, &obj2)
		if err != nil {
			log.Fatal(err)
		}

		if jsonDiffObj("{...}", obj1, obj2) {
			if _, ok := obj1["time"]; ok {
				//fmt.Println("==", strings.TrimSpace(string(line1)))
			}
		} else {
			if len(line1) > 1024 {
				fmt.Println("file1:", strings.TrimSpace(string(line1))[0:1024]+"...")
			} else {
				fmt.Println("file1:", strings.TrimSpace(string(line1)))
			}
			if len(line2) > 1024 {
				fmt.Println("file2:", strings.TrimSpace(string(line2))[0:1024]+"...")
			} else {
				fmt.Println("file2:", strings.TrimSpace(string(line2)))
			}
			return
		}
	}
}
