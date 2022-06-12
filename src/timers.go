package main

import (
"strings"
"fmt"
"context"
"strconv"
"encoding/json"
"encoding/hex"
"time"
"go.mongodb.org/mongo-driver/bson"
)

func timers() {
    for {
        if time.Now().Minute() % 5 == 0 && time.Now().Second() % 60 == 0 {
            process_block_data()
            time.Sleep(1 * time.Second)
        }
        time.Sleep(1 * time.Second)
    }
}

func timers_build_data() {
    block_height := XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT
    for {
            process_block_data_build_data(block_height)
            block_height++
            time.Sleep(10 * time.Second)
        }
    }


func process_block_data() {

  // Variables
  var block_height int
  var s string
  var delegate string
  var length int
  var timestamp int
  var reward int64
  var public_tx_count int = 0
  var data_send string
  var data_read_1 TxData
  var error error
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the previous block Height
  if block_height = get_current_block_height(); block_height == 0 {
     return
  }
  block_height -= 1
  
  // get the reserve bytes
  s = get_reserve_bytes(block_height)
  if s == "" {
    return
  }
  
  // parse the reserve bytes 
  delegate = s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(s, BLOCKCHAIN_DATA_SEGMENT_STRING)]
  delegate_name_data,_ := hex.DecodeString(delegate)
  delegate = string(delegate_name_data)
    
  length = len(s) - (len(s) - strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START_DATA)) - 106 - 142

  timestamp = int(varint_decode(s[4:14]))
  reward = varint_decode(s[106 : 106+length])
  
  tx_data := s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_END)+len(BLOCKCHAIN_RESERVED_BYTES_END)+4:]
  if len(tx_data) < TRANSACTION_HASH_LENGTH {
      public_tx_count = 0
  } else {
      for len(tx_data) >= TRANSACTION_HASH_LENGTH {
      tx := tx_data[0:TRANSACTION_HASH_LENGTH]
      tx_data = tx_data[TRANSACTION_HASH_LENGTH:]
      fmt.Println(tx)
      // get the tx details
  data_send,error = send_http_data("http://127.0.0.1:18281/get_transactions",`{"txs_hashes":["` + tx + `"]}`)
  if !strings.Contains(data_send, "\"status\": \"OK\"") || error != nil {
    return
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    return
  }
  
  // get the public tx info
  if strings.Contains(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX) {
      public_tx_count++
      
      // parse the public tx
  }
      }
  }
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").InsertOne(ctx, bson.D{{"height", strconv.Itoa(block_height)}, {"delegate", delegate},{"reward", strconv.FormatInt(reward, 10)},{"time", strconv.Itoa(timestamp)},{"public_tx", strconv.Itoa(public_tx_count)}})
 
  fmt.Println(block_height)
   fmt.Println(delegate)
   fmt.Println(timestamp)
   fmt.Println(reward)
   fmt.Println(public_tx_count)
    
  return
}

func process_block_data_build_data(block_height int) {

  // Variables
  var s string
  var delegate string
  var length int
  var timestamp int
  var reward int64
  var public_tx_count int = 0
  var data_send string
  var data_read_1 TxData
  var error error
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the reserve bytes
  s = get_reserve_bytes(block_height)
  if s == "" {
    return
  }
  
  // parse the reserve bytes 
  delegate = s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(s, BLOCKCHAIN_DATA_SEGMENT_STRING)]
  delegate_name_data,_ := hex.DecodeString(delegate)
  delegate = string(delegate_name_data)
    
  length = len(s) - (len(s) - strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START_DATA)) - 106 - 142

  timestamp = int(varint_decode(s[4:14]))
  reward = varint_decode(s[106 : 106+length])
  
  tx_data := s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_END)+len(BLOCKCHAIN_RESERVED_BYTES_END)+4:]
  if len(tx_data) < TRANSACTION_HASH_LENGTH {
      public_tx_count = 0
  } else {
      for len(tx_data) >= TRANSACTION_HASH_LENGTH {
      tx := tx_data[0:TRANSACTION_HASH_LENGTH]
      tx_data = tx_data[TRANSACTION_HASH_LENGTH:]
      fmt.Println(tx)
      // get the tx details
  data_send,error = send_http_data("http://127.0.0.1:18281/get_transactions",`{"txs_hashes":["` + tx + `"]}`)
  if !strings.Contains(data_send, "\"status\": \"OK\"") || error != nil {
    return
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    return
  }
  
  // get the public tx info
  if strings.Contains(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX) {
      public_tx_count++
      
      // parse the public tx
  }
      }
  }
  
  // save the block to the database
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").InsertOne(ctx, bson.D{{"height", strconv.Itoa(block_height)}, {"delegate", delegate},{"reward", strconv.FormatInt(reward, 10)},{"time", strconv.Itoa(timestamp)},{"public_tx", strconv.Itoa(public_tx_count)}})
  
  fmt.Println(block_height)
   fmt.Println(delegate)
   fmt.Println(timestamp)
   fmt.Println(reward)
   fmt.Println(public_tx_count)
    
  return
}
