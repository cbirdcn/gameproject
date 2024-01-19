package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"context"
	"github.com/go-redis/redis/v8"
	"encoding/json"
	"strconv"
	"sync"
	"strings"
	"time"
)

type Account struct {
	Id int `gorm:"id"`
	Name string `gorm:"name"`
	Password string `gorm:"password"`
}

func (Account) TableName() string {
	return "account"
}

func main() {
	dsn := "root:root@tcp(host.docker.internal:3307)/db_account?charset=utf8"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("db open err")
	}
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "host.docker.internal:6379",
		Password:"",
		DB:0,
	})


	var accounts []Account
	res := db.Find(&accounts)
	if res.Error != nil {
		fmt.Println(res.Error)
	} else {
		for _,account := range accounts {
			idString := strconv.Itoa(account.Id)
			key := fmt.Sprintf("account_%s", idString)
			mapv := make(map[string]interface{})
			ma, _ := json.Marshal(account)
			value := make([]interface{}, 0)
			_ = json.Unmarshal(ma, &mapv)
			for k, v := range mapv {
				value = append(value, k)
				value = append(value, v)
			}
			err = rdb.HMSet(ctx, key, value...).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println(accounts)
	}
	
	wg := sync.WaitGroup{}
	wg.Add(1)
	go SaveAccountCoroutine(ctx, rdb, db, &wg)
	wg.Wait()
}

func SaveAccountCoroutine(ctx context.Context, rdb *redis.Client, db *gorm.DB, wg *sync.WaitGroup) {
//	for {
		key := "changed_account"
		length, _ := rdb.LLen(ctx,key).Result()
		if length > 0 {
			fmt.Println(length)
			duration, _ := time.ParseDuration("60s")
			val, _ := rdb.BRPop(ctx, duration, key).Result()
			fmt.Println(val)
			mapv := make(map[string]interface{})
			_ = json.Unmarshal([]byte(val[0]), &mapv)
			tb := mapv["tb"]
			id := mapv["id"]
			op := mapv["op"]
			where := mapv["where"]
			// test
			where = "id = 1"
			sql := "test"
			fmt.Println(mapv)
			return

			switch op {
				case "soft_delete":
					sql = fmt.Sprintf("update %s set is_delete = 1 where %s", tb, where)
				case "upsert":
					hkey := fmt.Sprintf("account_%s", id)
					hkeys, _ := rdb.HKeys(ctx, hkey).Result()
					keys := strings.Join(hkeys, ",")
					hvals, _ := rdb.HVals(ctx, hkey).Result()
					vals := strings.Join(hvals, ",")
					sql = fmt.Sprint("replace into %s(%s) values(%s)", tb, keys, vals)
				default :
					break
			}
			fmt.Println(sql)

		}
//	}
	wg.Done()
}
