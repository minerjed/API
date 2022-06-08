package main

import (
"strings"
"strconv"
"os"
"path/filepath"
"encoding/json"
"time"
"github.com/gofiber/fiber/v2"
)

func blockchain_size() (int64, error) {
	var size int64
	err := filepath.Walk(BLOCKCHAIN_DIRECTORY, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func v1_xcash_blockchain_unauthorized_stats(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 BlockchainStats
  var data_read_2 BlockchainBlock
  var output v1XcashBlockchainUnauthorizedStats;
  var count int
  generated_supply := FIRST_BLOCK_MINING_REWARD + XCASH_PREMINE_TOTAL_SUPPLY
  generated_supply_copy := FIRST_BLOCK_MINING_REWARD + XCASH_PREMINE_TOTAL_SUPPLY
  var reward float64
  
  // get info
  data_send = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_info"}`)
  if !strings.Contains(data_send, "\"result\"") {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }

  // get block
  data_send = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block","params":{"height":` + strconv.FormatInt(int64(data_read_1.Result.Height-1), 10) + `}}`)
  if !strings.Contains(data_send, "\"result\"") {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_2); err != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }

  // get the generated supply
  for count = 2; count < data_read_1.Result.Height; count++ {
    if count < XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT {
      generated_supply = generated_supply + (XCASH_TOTAL_SUPPLY - generated_supply) / XCASH_EMMISION_FACTOR
    } else {
      reward = ((XCASH_TOTAL_SUPPLY - generated_supply) / XCASH_DPOPS_EMMISION_FACTOR)
      generated_supply += reward
      if (reward * XCASH_WALLET_DECIMAL_PLACES_AMOUNT) <= EMISSION_BLOCK_REWARD {
        break
      }
    }
  }
  circulating_supply := int64(((generated_supply - (XCASH_PREMINE_TOTAL_SUPPLY - XCASH_PREMINE_CIRCULATING_SUPPLY)) * XCASH_WALLET_DECIMAL_PLACES_AMOUNT))
  generated_supply *= XCASH_WALLET_DECIMAL_PLACES_AMOUNT

  // get the emission data
  for count = 2; count < 2000000; count++ {
    if count < XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT {
      generated_supply_copy = generated_supply_copy + (XCASH_TOTAL_SUPPLY - generated_supply_copy) / XCASH_EMMISION_FACTOR
    } else {
      reward = ((XCASH_TOTAL_SUPPLY - generated_supply_copy) / XCASH_DPOPS_EMMISION_FACTOR)
      generated_supply_copy += reward
      if (reward * XCASH_WALLET_DECIMAL_PLACES_AMOUNT) <= EMISSION_BLOCK_REWARD {
        break
      }
    }
  }
  emission_height := count
  timestamp, _ := strconv.Atoi(strconv.FormatInt(time.Now().UTC().Unix(), 10))
  emission_time := (timestamp) + ((count - data_read_1.Result.Height) * (XCASH_DPOPS_BLOCK_TIME * 60))

  // get the emission data
  generated_supply_copy = FIRST_BLOCK_MINING_REWARD + XCASH_PREMINE_TOTAL_SUPPLY

  for count = 2; count < 2000000; count++ {
    if count < XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT {
      generated_supply_copy = generated_supply_copy + (XCASH_TOTAL_SUPPLY - generated_supply_copy) / XCASH_EMMISION_FACTOR
    } else {
      reward = ((XCASH_TOTAL_SUPPLY - generated_supply_copy) / XCASH_DPOPS_EMMISION_FACTOR)
      if ((reward * XCASH_WALLET_DECIMAL_PLACES_AMOUNT) <= EMISSION_BLOCK_REWARD) {
        reward = (EMISSION_BLOCK_REWARD / XCASH_WALLET_DECIMAL_PLACES_AMOUNT)
      }
      generated_supply_copy += reward
      if (generated_supply_copy) >= XCASH_TOTAL_SUPPLY {
        break
      }
    }
  }
  inflation_height := count
  inflation_timestamp, _ := strconv.Atoi(strconv.FormatInt(time.Now().UTC().Unix(), 10))
  inflation_time := (inflation_timestamp) + ((inflation_height - data_read_1.Result.Height) * (XCASH_DPOPS_BLOCK_TIME * 60))  

  // get the blockchain size
  blockchain_data_size,err := blockchain_size()
  if err != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }

  // fill in the data
  output.Height = data_read_1.Result.Height
  output.Hash = data_read_1.Result.TopBlockHash
  output.Reward = data_read_2.Result.BlockHeader.Reward
  output.Size = blockchain_data_size
  output.Version = CURRENT_BLOCKCHAIN_VERSION
  output.VersionBlockHeight = CURRENT_BLOCKCHAIN_VERSION_HEIGHT
  output.NextVersionBlockHeight = NEXT_BLOCKCHAIN_VERSION_HEIGHT
  output.TotalTx = data_read_1.Result.TxCount
  output.CirculatingSupply = circulating_supply
  output.GeneratedSupply = int64(generated_supply)
  output.TotalSupply = XCASH_TOTAL_SUPPLY
  output.EmissionReward = EMISSION_BLOCK_REWARD
  output.EmissionHeight = emission_height
  output.EmissionTime = emission_time
  output.InflationHeight = inflation_height
  output.InflationTime = inflation_time
    
  return c.JSON(output)
}
