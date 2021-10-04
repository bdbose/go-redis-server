package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Response struct {
	Detail []struct {
		Title     string   `json:"title"`
		Descp     string   `json:"descp"`
		Source    string   `json:"source"`
		MetaExtra []string `json:"metaExtra"`
		MetaTag   []string `json:"metaTag"`
		NewsURL   string   `json:"newsUrl"`
		ImgURL    string   `json:"imgUrl"`
	} `json:"detail"`
}

func GetNewsDat() (string){
	res,err := http.Get("https://lit-woodland-31639.herokuapp.com/world")
	if err != nil{
		fmt.Println(err.Error())
	}
	responseData, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil{
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	return string(responseData)
}

func CheckCache(c redis.Client )(string,error){
	val, err := c.Get("id1234").Result()
	if err != nil{
		fmt.Println("FETCHING FROM API")
		val = GetNewsDat()
		err := c.Set("id1234", val, 10*time.Minute).Err()
		if err != nil {
			fmt.Println(err.Error())
			return val,err
		}
		return val,err
	}
	
	fmt.Println("FETCHING FROM REDIS")
	return val,err
}

func main() {
	r := gin.Default()
	redisAddr := os.Getenv("REDISADDRESS")
	redisPass := os.Getenv("REDISPASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: redisPass,
		DB: 0,
	})
	
	r.GET("/", func(c *gin.Context) {
		val,err := CheckCache(*client)
		if err != nil {
			c.JSON(400, gin.H{
				"err": err,
			})
			return
		}
		var response Response
		json.Unmarshal([]byte(val),&response)
		c.JSON(200, gin.H{
			"data": response.Detail,
		})
	})
	r.Run() 
}