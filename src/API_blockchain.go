package main

import (
"fmt"
"strconv"
"github.com/gofiber/fiber/v2"
)

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
