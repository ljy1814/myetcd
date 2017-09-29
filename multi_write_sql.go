package main

import (
	_ "ac-common-go/mysql"
	crand "crypto/rand"
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	user     = "root"
	password = "root"
	host     = "127.0.0.1:3306"
	database = "test"
)

func main1() {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database)
	driver, err := sql.Open("mysql", dns)
	if err != nil {
		log.Errorf("open MySQL Driver failed:dns[%s], err[%v]", dns, err)
		return
	}
	err = driver.Ping()
	if err != nil {
		log.Errorf("MySQL ping failed:dns[%s], err[%v]", dns, err)
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for count := 0; count < 20000; count++ {
				SQL := "INSERT INTO news(title,author, keywords,description) VALUES("
				SQL += "'" + randomStr() + "',"
				SQL += "'" + randomStr() + "',"
				SQL += "'" + randomStr() + "',"
				SQL += "'" + randomStr() + "')"
				_, err := driver.Exec(SQL)
				if err != nil {
					log.Errorf("err[%v]", err)
				}
			}
		}()
	}
	//	time.Sleep(20 * time.Second)
	wg.Wait()
}

func randomStr() string {
	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	var bytes = make([]byte, 5+(rn.Int()%59))
	charTab := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
	crand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = charTab[b%byte(len(charTab))]
	}
	return string(bytes)
}
