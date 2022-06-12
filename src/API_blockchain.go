package main

import (
"fmt"
"math/rand"
"strings"
"context"
"strconv"
"os"
"path/filepath"
"encoding/json"
"encoding/hex"
"time"
"github.com/gofiber/fiber/v2"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson"
)

func blockchain_size() (int64, error) {
  // Variables
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

func RandStringBytes(n int) string {
  // Constants
  const letterBytes = "0123456789abcdef"

  b := make([]byte, n)
  for i := range b {
      b[i] = letterBytes[rand.Intn(len(letterBytes))]
  }
  return string(b)
}

func get_current_block_height() int {
  // Variables
  var data_read CurrentBlockHeight
  var data_send string
  var error error
  
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block_count"}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    return 0
  }
  if err := json.Unmarshal([]byte(data_send), &data_read); err != nil {
    return 0
  }
  return data_read.Result.Count
}

func get_block_delegate(requestBlockHeight int) string {
  // Variables
  var database_data XcashDpopsReserveBytesCollection

  // get the collection
  block_height_data := strconv.Itoa(int(((requestBlockHeight - XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT) / BLOCKS_PER_DAY_FIVE_MINUTE_BLOCK_TIME))+1)
  collection_number := "reserve_bytes_" + block_height_data
  collection := mongoClient.Database(XCASH_DPOPS_DATABASE).Collection(collection_number)

  // get the reserve bytes
  filter := bson.D{{"block_height", strconv.Itoa(requestBlockHeight)}}
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  err := collection.FindOne(ctx, filter).Decode(&database_data)
  if err == mongo.ErrNoDocuments {
    return ""
  } else if err != nil {
    return ""
  }

  // get the delegate name from the reserve bytes
  delegate_name := database_data.ReserveBytes[strings.Index(database_data.ReserveBytes, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(database_data.ReserveBytes, BLOCKCHAIN_DATA_SEGMENT_STRING)]
  delegate_name_data, err := hex.DecodeString(delegate_name)
  if err != nil {
    return ""
  }

  return string(delegate_name_data)   
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
  var error error
  var database_data XcashAPIStatisticsCollection
  
  // read the tx stats
  collection := mongoClient.Database(XCASH_API_DATABASE).Collection("statistics")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  err := collection.FindOne(ctx, bson.D{{}}).Decode(&database_data)
  if err == mongo.ErrNoDocuments {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  } else if err != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }
  
  // get info
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_info"}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not get the stats"}
    return c.JSON(error)
  }

  // get block
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block","params":{"height":` + strconv.FormatInt(int64(data_read_1.Result.Height-1), 10) + `}}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
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
  output.TotalPublicTx,_ = strconv.Atoi(database_data.Public)
  output.TotalPrivateTx,_ = strconv.Atoi(database_data.Private)
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

func v1_xcash_blockchain_unauthorized_blocks_blockHeight(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 BlockchainStats
  var data_read_2 BlockchainBlock
  var data_read_3 BlockchainBlockJson
  var output v1XcashBlockchainUnauthorizedBlocksBlockHeight;
  var requestBlockHeight string
  var xcash_dpops_status bool
  var xcash_dpops_delegate string
  var error error
  
  // get info
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_info"}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := ErrorResults{"Could not get the block data"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not get the block data"}
    return c.JSON(error)
  }

  // get the resource
  requestBlockHeight = c.Params("blockHeight")
  if requestBlockHeight == "" {
    requestBlockHeight = strconv.FormatInt(int64(data_read_1.Result.Height-1), 10)
  }

  // get block
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block","params":{"height":` + requestBlockHeight + `}}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := ErrorResults{"Could not get the block data"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_2); err != nil {
    error := ErrorResults{"Could not get the block data"}
    return c.JSON(error)
  }

  // get the tx
  s := string(data_read_2.Result.JSON)
  s = strings.Replace(s, "\\n", "", -1)
  s = strings.Replace(s, "\\", "", -1)
  fmt.Println(s)
  if err := json.Unmarshal([]byte(s), &data_read_3); err != nil {
    error := ErrorResults{"Could not get the block data"}
    return c.JSON(error)
  }

  // get the dpops block status
  if (data_read_2.Result.BlockHeader.Height >= XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT) {
    xcash_dpops_status = true
    xcash_dpops_delegate = get_block_delegate(data_read_2.Result.BlockHeader.Height)
  } else {
    xcash_dpops_status = false
    xcash_dpops_delegate = ""
  }

  // fill in the data
  output.Height = data_read_2.Result.BlockHeader.Height
  output.Hash = data_read_2.Result.BlockHeader.Hash
  output.Reward = data_read_2.Result.BlockHeader.Reward
  output.Time = data_read_2.Result.BlockHeader.Timestamp
  output.XcashDPOPS = xcash_dpops_status
  output.DelegateName = xcash_dpops_delegate
  output.Tx = data_read_3.TxHashes
    
  return c.JSON(output)
}

func v1_xcash_blockchain_unauthorized_tx_prove(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 CheckTxKey
  var data_read_2 CheckTxProof
  var output v1XcashBlockchainUnauthorizedTxProve;
  var amount int64
  var valid bool
  var post_data v1XcashBlockchainUnauthorizedTxProvePostData
  var error error

  if err := c.BodyParser(&post_data); err != nil {
    error := v1XcashBlockchainUnauthorizedTxProve{false,0}
    return c.JSON(error)
  }

  // error check
  if post_data.Tx == "" || post_data.Address == "" || post_data.Key == "" || len(post_data.Tx) != TRANSACTION_HASH_LENGTH || len(post_data.Address) != XCASH_WALLET_LENGTH || post_data.Address[0:len(XCASH_WALLET_PREFIX)] != XCASH_WALLET_PREFIX || (len(post_data.Key) != TRANSACTION_HASH_LENGTH && post_data.Key[0:len(CHECK_TX_PROOF_PREFIX)] != CHECK_TX_PROOF_PREFIX) {
    error := v1XcashBlockchainUnauthorizedTxProve{false,0}
    return c.JSON(error)
  }
  
  if len(post_data.Key) == TRANSACTION_HASH_LENGTH {
    // get info
    data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + post_data.Tx + `","tx_key":"` + post_data.Key + `","address":"` + post_data.Address + `"}}`)
    if !strings.Contains(data_send, "\"result\"") || error != nil {
      error := v1XcashBlockchainUnauthorizedTxProve{false,0}
      return c.JSON(error)
    }
    if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
      error := v1XcashBlockchainUnauthorizedTxProve{false,0}
      return c.JSON(error)
    }
    
    valid = true
    amount = data_read_1.Result.Received
  } else {
    // get info
    data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_proof","params":{"txid":"` + post_data.Tx + `","address":"` + post_data.Address + `","signature":"` + post_data.Key + `"}}`)
    if !strings.Contains(data_send, "\"result\"") || error != nil {
      error := v1XcashBlockchainUnauthorizedTxProve{false,0}
      return c.JSON(error)
    }
    if err := json.Unmarshal([]byte(data_send), &data_read_2); err != nil {
      error := v1XcashBlockchainUnauthorizedTxProve{false,0}
      return c.JSON(error)
    }   

    valid = data_read_2.Result.Good
    amount = data_read_2.Result.Received
  }

  // fill in the data
  output.Valid = valid
  output.Amount = amount
    
  return c.JSON(output)
}

func v1_xcash_blockchain_unauthorized_address_prove(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 CheckReserveProof
  var output v1XcashBlockchainUnauthorizedAddressProve;
  var amount int64
  var post_data v1XcashBlockchainUnauthorizedAddressProvePostData
  var error error

  if err := c.BodyParser(&post_data); err != nil {
    error := v1XcashBlockchainUnauthorizedAddressProve{0}
    return c.JSON(error)
  }

  // error check
  if post_data.Address == "" || len(post_data.Address) != XCASH_WALLET_LENGTH || post_data.Address[0:len(XCASH_WALLET_PREFIX)] != XCASH_WALLET_PREFIX || post_data.Signature == "" || post_data.Signature[0:len(CHECK_RESERVE_PROOF_PREFIX)] != CHECK_RESERVE_PROOF_PREFIX {
    error := v1XcashBlockchainUnauthorizedAddressProve{0}
    return c.JSON(error)
  }  
  
  // get info
  data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_reserve_proof","params":{"address":"` + post_data.Address + `","signature":"` + post_data.Signature + `"}}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := v1XcashBlockchainUnauthorizedAddressProve{0}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := v1XcashBlockchainUnauthorizedAddressProve{0}
    return c.JSON(error)
  }   

  if data_read_1.Result.Good && data_read_1.Result.Spent == 0 {
    amount = data_read_1.Result.Total
  } else {
    amount = 0
  }

  // fill in the data
  output.Amount = amount
    
  return c.JSON(output)
}

func v1_xcash_blockchain_unauthorized_address_create_integrated(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 CreateIntegratedAddress
  var output v1XcashBlockchainUnauthorizedAddressCreateIntegrated;
  var post_data v1XcashBlockchainUnauthorizedAddressCreateIntegratedPostData
  var error error

  if err := c.BodyParser(&post_data); err != nil {
    error := ErrorResults{"Could not create the integrated address"}
    return c.JSON(error)
  }

  // error check
  if post_data.Address == "" || len(post_data.Address) != XCASH_WALLET_LENGTH || post_data.Address[0:len(XCASH_WALLET_PREFIX)] != XCASH_WALLET_PREFIX || (post_data.PaymentID != "" && len(post_data.PaymentID) != ENCRYPTED_PAYMENT_ID_LENGTH) {
    error := ErrorResults{"Could not create the integrated address"}
    return c.JSON(error)
  }  
  
  // get info
  if post_data.PaymentID == "" {
    post_data.PaymentID = RandStringBytes(ENCRYPTED_PAYMENT_ID_LENGTH);
  }

  data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"make_integrated_address","params":{"standard_address":"` + post_data.Address + `", "payment_id":"` + post_data.PaymentID + `"}}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := ErrorResults{"Could not create the integrated address"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not create the integrated address"}
    return c.JSON(error)
  } 

  // fill in the data
  output.IntegratedAddress = data_read_1.Result.IntegratedAddress
  output.PaymentID = data_read_1.Result.PaymentID
    
  return c.JSON(output)
}

func v1_xcash_blockchain_unauthorized_tx_txHash(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var data_read_1 TxData
  var data_read_2 CheckTxKey
  var data_read_3 CurrentBlockHeight
  var output v1XcashBlockchainUnauthorizedTxTxHash;
  var tx string
  var sender_data string
  var receiver_data string
  var key string
  var error error

  // get the resource
  if tx = c.Params("txHash"); tx == "" {
    error := ErrorResults{"Could not get the tx details"}
    return c.JSON(error)
  }
  
  // get info
  data_send,error = send_http_data("http://127.0.0.1:18281/get_transactions",`{"txs_hashes":["` + tx + `"]}`)
  if !strings.Contains(data_send, "\"status\": \"OK\"") || error != nil {
    error := ErrorResults{"Could not get the tx details"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not get the tx details"}
    return c.JSON(error)
  }
  
  // get the public tx info
  if strings.Contains(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX) {
    output.Type = "public"
    
    // decode the tx data
    key = data_read_1.TxsAsHex[0][strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX)+len(PUBLIC_TX_PREFIX):strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX)+len(PUBLIC_TX_PREFIX)+TRANSACTION_HASH_LENGTH]
    sender_data = data_read_1.TxsAsHex[0][strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_XCASH_PREFIX)+2 : strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_XCASH_PREFIX)+2+(XCASH_WALLET_LENGTH*2)]
    data_read_1.TxsAsHex[0] = strings.Replace(data_read_1.TxsAsHex[0], sender_data, "", -1)
    receiver_data = data_read_1.TxsAsHex[0][strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_XCASH_PREFIX)+2 : strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_XCASH_PREFIX)+2+(XCASH_WALLET_LENGTH*2)]
    
    data1,_ := hex.DecodeString(receiver_data)
    data2,_ := hex.DecodeString(sender_data)
    
    output.Receiver = string(data1)
    output.Sender = string(data2)
    
    // get the amount
    data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + tx + `","tx_key":"` + key + `","address":"` + output.Receiver + `"}}`)
    if !strings.Contains(data_send, "\"result\"") || error != nil {
      error := ErrorResults{"Could not get the tx details"}
      return c.JSON(error)
    }
    if err := json.Unmarshal([]byte(data_send), &data_read_2); err != nil {
      error := ErrorResults{"Could not get the tx details"}
      return c.JSON(error)
    }
    
    output.Amount = data_read_2.Result.Received
  } else {
    output.Type = "private"
    output.Receiver = ""
    output.Sender = ""
    output.Amount = 0
  }
  
  // get the current block Height
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block_count"}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    error := ErrorResults{"Could not get the tx details"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_3); err != nil {
    error := ErrorResults{"Could not get the tx details"}
    return c.JSON(error)
  }

  // fill in the data
  output.Height = data_read_1.Txs[0].BlockHeight
  output.Confirmations = data_read_3.Result.Count - data_read_1.Txs[0].BlockHeight
  output.Time = data_read_1.Txs[0].BlockTimestamp
    
  return c.JSON(output)
}
