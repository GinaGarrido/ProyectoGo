package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Results struct {
	Data []Clave `json:"results"`
}
type Clave struct {
	Name                   string `json:"name"`
	Model                  string `json:"model"`
	Manufacturer           string `json:"manufacturer"`
	Cost_in_credits        string `json:"cost_in_credits"`
	Length                 string `json:"length"`
	Max_atmosphering_speed string `json:"max_atmosphering_speed"`
	Crew                   string `json:"crew"`
	Passengers             string `json:"passengers"`
	Cargo_capacity         string `json:"cargo_capacity"`
	Consumables            string `json:"consumables"`
	Hyperdrive_rating      string `json:"hyperdrive_rating"`
	Mglt                   string `json:"MGLT"`
	Starship_class         string `json:"starship_class"`
	Pilots                 string `json:"pilots"`
	Films                  string `json:"films"`
	Created                string `json:"created"`
	Edited                 string `json:"edited"`
	URL                    string `json:"url"`
}

//Variable cliente global
var client *redis.Client

//Creación de cliente
func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

//Consumir API
func getAPI() (responseObject Results) {
	response, err := http.Get("https://swapi.co/api/starships/?page=1")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &responseObject)
	return
}

func main() {
	r := gin.Default()
	r.POST("/starships", loadStarships)
	r.GET("/starships/:val", getStarship)

	client = newClient()

	r.Run(":3000")
}

func loadStarships(c *gin.Context) {
	responseObject := getAPI()

	for i := 0; i < len(responseObject.Data); i++ {
		var nombre = responseObject.Data[i].Name
		b, _ := json.Marshal(responseObject.Data[i])
		err := client.Set(nombre, string(b), 0).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getStarship(c *gin.Context) {
	val, err := client.Get(c.Param("val")).Result()
	if val != "" {
		ship := Clave{}
		json.Unmarshal([]byte(val), &ship)
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(200, ship)
	} else {
		c.JSON(404, gin.H{
			"error": "Ship Not Found",
		})
	}
}
