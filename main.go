package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var arr = []string{"."}

func main() {
	var mainJSON map[string][]string = make(map[string][]string)

	for a := 0; a < len(arr); a++ {

		files, err := ioutil.ReadDir(arr[a])
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fName := f.Name()
			fi, err := os.Stat(arr[a] + "/" + fName)
			if err == nil {
				mode := fi.Mode()

				if mode.IsDir() {
					arr = append(arr, fName)
					continue
				}
			}
			res1 := strings.Index(fName, ".")
			if res1 != -1 && (fName[res1:] == ".ipset" || fName[res1:] == ".netstat") {
				file, err := os.Open(arr[a] + "/" + fName)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				mainCat := ""
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					txt := scanner.Text()
					catlen := len("# Category")
					if txt[0] != []byte("#")[0] {
						// fmt.Println(scanner.Text())
						mainJSON[mainCat] = append(mainJSON[mainCat], scanner.Text())
					} else if len(txt) >= catlen && txt[:catlen] == "# Category" {
						cat := txt[catlen:]
						catTrim := strings.Trim(cat, " ")
						if catTrim[0] == []byte(":")[0] && len(catTrim) > 2 {
							catTrim = catTrim[2:]
							_, ok := mainJSON[catTrim]
							if !ok {
								mainJSON[catTrim] = []string{}
							}
							mainCat = catTrim
						}
					}
				}

				if err := scanner.Err(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	d1, _ := json.Marshal(mainJSON)
	ioutil.WriteFile("output.json", d1, 0644)
	// fmt.Println(json.Marshal(mainJSON))
}
