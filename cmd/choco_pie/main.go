package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/jackc/pgx/v5/pgxpool"
)

const countRecord = 100_000
const workerNum = 10

func main() {
	ctx := context.Background()
	pgMasterURL := fmt.Sprintf("postgres://%s:%s@%s/%s",
		"postgres",
		"postgres",
		"localhost:5432",
		"postgres",
	)
	config, err := pgxpool.ParseConfig(pgMasterURL)
	if err != nil {
		log.Fatal(err)
	}
	config.MinConns = 10
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var recorted int64

	wg := sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < countRecord/workerNum; i++ {
				_, err = db.Exec(ctx, `insert into choco_pie(value) values($1)`, i+1)
				if err != nil {
					log.Fatal(err, recorted)
				}
				r := atomic.AddInt64(&recorted, 1)
				log.Println(r)
			}
		}()
	}

	wg.Wait()

	log.Println("record", countRecord)
}
