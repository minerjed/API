package main
 
import (
"fmt"
"math/rand"
"strings"
"context"
"strconv"
"bytes"
"io/ioutil"
"net/http"
"time"
"encoding/json"
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data))); err != nil {
		return "error1"
	}
  
	req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req); err != nil {
		return "error2"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
fmt.Printf("for %s sending %s received %s\n", url,data,body)

        return string(body)
}
 

func get_http_data(url string) string {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(""))); err != nil {
		return "error"
	}	
  
req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req); err != nil {
		return "error"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        return string(body)
}
func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
 
func main() {
  
// setup fiber
app := fiber.New(fiber.Config{
Prefork: true,
DisableStartupMessage: true,
})
 
app.Post("/processturbotx/", func(c *fiber.Ctx) error {
    // Variables
    var id string
    tx_hash := c.Query("tx_hash")
    tx_key := c.Query("tx_key")
    sender := c.Query("sender")
    receiver := c.Query("receiver")
    amount := c.Query("amount")

    // error check
    if (len(tx_hash) != TX_HASH_LENGTH || len(tx_key) != TX_HASH_LENGTH || len(sender) != PUBLIC_ADDRESS_LENGTH || len(receiver) != PUBLIC_ADDRESS_LENGTH) {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    if _, err := strconv.Atoi(amount); err != nil {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    // get the id
    id = tx_hash[:IDLENGTH]

    // get the timestamp
    timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
 
    // save the data in the database
    data := string(`{"id": "` + id + `", "tx_hash": "` + tx_hash + `", "tx_key": "` + tx_key + `", "timestamp": "` + timestamp + `", "sender": "` + sender + `", "receiver": "` + receiver + `", "amount": "` + c.Query("amount") + `"}`)
fmt.Printf("str1: %s\n", data)

    err := rdb.Set(ctx, id, data, 1*time.Hour).Err()
    if err != nil {
        error := ErrorResults{"error"}
        return c.JSON(error)
    }
 
    // return the id
    return c.SendString(URL + id + "}")
})
 
app.Get("/getturbotx/", func(c *fiber.Ctx) error {
  id := c.Query("id")
val, _ := rdb.Get(ctx, id).Result()
    if val == "" {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }
fmt.Printf("%s\n", val)
   // convert the string to a json object
   var data TurboTxSave
   json.Unmarshal([]byte(val), &data)
 fmt.Printf("str1: %s\n", "checking data")
fmt.Println("Struct is:", data)

   // check if the amount is correct and the sender and receiver are in the output
   datamount, _ := strconv.Atoi(data.Amount)
   amount,delegate_count,block_status,timestamp := turbo_tx_verify(data)

   if amount < datamount || amount <= 0 {
      error := ErrorResults{"error"}
      return c.JSON(error)
  } 

  if timestamp == "0" {
    timestamp = data.Timestamp
  }
 
  result := TurboTxOut{id, data.TX_Hash, timestamp, data.Sender, data.Receiver, strconv.FormatInt(int64(amount), 10),strconv.FormatInt(int64(delegate_count), 10),block_status}
  return c.JSON(result)
})

app.Static("/", "/var/www/html/")
 
app.Get("/*", func(c *fiber.Ctx) error {
  return c.SendString("Invalid URL")
})
 
  app.Listen(":8000")
}
