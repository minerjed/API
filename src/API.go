package main
 
import (
"github.com/gofiber/fiber/v2"
)

/*
global structures in src/structures.go
global constants in src/constants.go
network functions in src/network_functions.go
blockchain API in src/API_blockchain.go
*/

func main() {
  
// setup fiber
app := fiber.New(fiber.Config{
Prefork: true,
DisableStartupMessage: true,
})

// setup blockchain routes
app.Get("/v1/xcash/blockchain/unauthorized/stats/",v1_xcash_blockchain_unauthorized_stats)


// setup global routes
app.Get("/*", func(c *fiber.Ctx) error {
  return c.SendString("Invalid API Request")
})
 
  app.Listen(":9000")
}
