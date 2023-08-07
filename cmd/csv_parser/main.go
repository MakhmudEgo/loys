package main

import (
	"bytes"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dz1/internal/domain"
)

var config = flag.String("config", "./configs/people.csv", "csv file with users")

func main() {
	file, err := os.Open(*config)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	reader.Comment = '#'

	biography := make([]byte, 15)
	// count := 31241

	var count int64 = 130000
	time.Sleep(time.Second * 20)
	log.Println("skipping")
	for i := 0; i < int(count); i++ {
		_, err = reader.Read()
		if err != nil {
			fmt.Println("Ошибка при чтении файла:", err)
			return
		}
	}
	// 2023/07/13 18:52:46 *********** {"error":"failed to connect to `host=localhost user=postgres database=loys`: dial error (dial tcp 127.0.0.1:5432: connect: can't assign requested address)"}
	log.Println("skipping success")

	countWorkers := 30
	chRec := make(chan []string, countWorkers)
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < countWorkers; i++ {
		go func(i int) {
			log.Println("worker number", i+1)
			defer wg.Done()
			var e error
			for record := range chRec {
				parseInt, err := strconv.Atoi(record[1])
				if err != nil {
					log.Fatal(err)
				}
				date := time.Now().AddDate(-parseInt, 0, 0)
				gender := "male"
				if record[1][len(record[1])-1] == 'a' {
					gender = "female"
				}
				name := strings.Split(record[0], " ")
				_, _ = rand.Read(biography)
				user := domain.UserCreateReq{
					FirstName:  name[0],
					SecondName: name[1],
					Birthdate:  date.Format(time.DateOnly),
					Gender:     gender,
					Biography:  string(biography),
					City:       record[2],
					Password:   record[0] + record[1],
				}
				body, _ := json.Marshal(user)
				var resp *http.Response
				resp, e = http.Post("http://localhost:8080/user/register", "", bytes.NewReader(body))
				if e != nil || resp.StatusCode != http.StatusCreated {
					if resp != nil {
						all, _ := io.ReadAll(resp.Body)
						resp.Body.Close()
						log.Println("***********", string(all), "***********")
						log.Println("***********", e, "***********")
					}
					time.Sleep(time.Second * 10)
					for e != nil || resp.StatusCode != http.StatusCreated {
						resp, e = http.Post("http://localhost:8080/user/register", "", bytes.NewReader(body))
						if e != nil || resp.StatusCode != http.StatusCreated {
							if resp != nil {
								all, _ := io.ReadAll(resp.Body)
								log.Println("***********", string(all), "***********")
								log.Println("***********", e, "***********")
							}
							time.Sleep(time.Second * 10)
						}
					}
					log.Println(e)
				}

				log.Println("========", atomic.LoadInt64(&count), user, resp.StatusCode, "========")
				atomic.AddInt64(&count, 1)
			}
		}(i)
	}
	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			break
		}
		chRec <- record
	}
	close(chRec)
	wg.Wait()
}
