package main
 
import (
"fmt"
"strconv"
"bytes"
"io/ioutil"
"net/http"
"time"
"github.com/gofiber/fiber/v2"
)
 
// global structures 

type ErrorResults struct {
    Error string `json:"Error"`
}
 
// global constants
const URL = "http://162.55.235.87/?id="
const TX_HASH_LENGTH = 64
const PUBLIC_ADDRESS_LENGTH = 98
 
// Functions
func send_http_data(url string,data string) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
        if err != nil {
	  return "error1"
	}
  
	req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
        if err != nil {
	  return "error1"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
fmt.Printf("for %s sending %s received %s\n", url,data,body)

        return string(body)
}
 

func get_http_data(url string) string {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
        if err != nil {
	  return "error1"
	}
  
req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
        if err != nil {
	  return "error1"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        return string(body)
}



func helloWorld(c *fiber.Ctx) error {
	 // Variables
    var id string
    tx_hash := c.Query("tx_hash")
    amount := c.Query("amount")

    // error check
    if (len(tx_hash) != TX_HASH_LENGTH) {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    if _, err := strconv.Atoi(amount); err != nil {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    
fmt.Printf("str1: %s\n", "data")

 
    // return the id
    return c.SendString(URL + id + "}")
}
 
func main() {
  
// setup fiber
app := fiber.New(fiber.Config{
Prefork: true,
DisableStartupMessage: true,
})

// setup routes
app.Post("/test/",helloWorld)
app.Static("/", "/var/www/html/") 
app.Get("/*", func(c *fiber.Ctx) error {
  return c.SendString("Invalid URL")
})
 
  app.Listen(":9000")
}
