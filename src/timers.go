package main

import (
"fmt"
"strings"
"context"
"strconv"
"encoding/json"
"encoding/hex"
"time"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson"
)

func timers() {
    var block_height int
    
    for {
        if time.Now().Minute() % XCASH_DPOPS_BLOCK_TIME == 3 && time.Now().Second() == 0 {
            // get the previous block Height
            if block_height = get_current_block_height(); block_height == 0 {
              return
            }
            block_height -= 1
            fmt.Printf("Processing block: %d\n",block_height)
            process_block_data(block_height)
            time.Sleep(1 * time.Second)
        }
        time.Sleep(1 * time.Second)
    }
}

func timers_build_data() {
    block_height := 0
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  if block_height == 0 {
    _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").InsertOne(ctx, bson.D{{"public", "0"}, {"private", "0"}})
  }
    for {
            fmt.Printf("Processing block: %d\n",block_height)
            process_block_data(block_height)
            block_height++
            time.Sleep(1 * time.Second)
        }
    }

func process_block_data(block_height int) {

  // Variables
  var s string
  var delegate string
  var public_tx_count int = 0
  var private_tx_count int = 0
  var data_send string
  var data_read_1 TxData
  var error error
  var database_data XcashAPIStatisticsCollection
  var data_read_2 CheckTxKey
  var amount int64
  var data_read_3 BlockchainBlock
  var data_read_4 BlockchainBlockJson
  var count int64
  
  var block_found bool = false
  
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // check to make sure you have not already added this block 
  count,err := mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").CountDocuments(ctx, bson.D{{"height",strconv.Itoa(block_height)}})
  if err != nil || count != 0 {
    block_found = true
  }
  
  // get the currrent tx count 
  err = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").FindOne(ctx, bson.D{{}}).Decode(&database_data)
  if err == mongo.ErrNoDocuments {
    return
  } else if err != nil {
    return
  }
  
  public_tx_count,_ = strconv.Atoi(database_data.Public)
  private_tx_count,_ = strconv.Atoi(database_data.Private)
  
  
  // get block
  data_send,error = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block","params":{"height":` + strconv.Itoa(block_height) + `}}`)
  if !strings.Contains(data_send, "\"result\"") || error != nil {
    return
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_3); err != nil {
    return
  }

  // get the tx
  s = string(data_read_3.Result.JSON)
  s = strings.Replace(s, "\\n", "", -1)
  s = strings.Replace(s, "\\", "", -1)
  if err := json.Unmarshal([]byte(s), &data_read_4); err != nil {
    return
  }
  
  // parse the reserve bytes 
  if block_height >= XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT {
    // get the reserve bytes
    delegate = get_reserve_bytes(block_height)
    if delegate == "" {
      return
    }
    delegate = delegate[strings.Index(delegate, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(delegate, BLOCKCHAIN_DATA_SEGMENT_STRING)]
    delegate_name_data,_ := hex.DecodeString(delegate)
    delegate = string(delegate_name_data)
  } else {
     delegate = "" 
  }
  
  
  for _, tx := range data_read_4.TxHashes {
      // check to make sure you have not already added this tx 
  count,err = mongoClient.Database(XCASH_API_DATABASE).Collection("tx").CountDocuments(ctx, bson.D{{"tx",tx}})
  if err != nil || count != 0 {
    continue
  }
  
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
      data := data_read_1.TxsAsHex[0][strings.Index(data_read_1.TxsAsHex[0], PUBLIC_TX_PREFIX)+len(PUBLIC_TX_PREFIX):]
      key := data[0:PUBLIC_KEY_LENGTH]
      data = data[PUBLIC_KEY_LENGTH+202:]
      sender := data[0:XCASH_WALLET_LENGTH*2]
      data = data[(XCASH_WALLET_LENGTH*2)+8:]
      receiver := data[0:XCASH_WALLET_LENGTH*2]
     
      sender_data,_ := hex.DecodeString(sender)
      receiver_data,_ := hex.DecodeString(receiver)
      
      sender = string(sender_data)
      receiver = string(receiver_data)
      
      // get the amount
      data_send,error = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + tx + `","tx_key":"` + key + `","address":"` + receiver + `"}}`)
    if !strings.Contains(data_send, "\"result\"") || error != nil {
      return
    }
    if err := json.Unmarshal([]byte(data_send), &data_read_2); err != nil {
      return
    }
    
    amount = data_read_2.Result.Received
    
    // save the public tx in the Database
    _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("tx").InsertOne(ctx, bson.D{{"tx", tx}, {"key", key},{"sender", sender},{"receiver", receiver},{"amount", strconv.FormatInt(amount, 10)},{"height", strconv.Itoa(data_read_1.Txs[0].BlockHeight)},{"time", strconv.Itoa(data_read_1.Txs[0].BlockTimestamp)}})
      
      
  } else {
      private_tx_count++
  }
      }
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"public", strconv.Itoa(public_tx_count)}}}})
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"private", strconv.Itoa(private_tx_count)}}}})
  if !block_found {
    _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").InsertOne(ctx, bson.D{{"height", strconv.Itoa(block_height)}, {"delegate", delegate},{"reward", strconv.FormatInt(data_read_3.Result.BlockHeader.Reward, 10)},{"time", strconv.Itoa(data_read_3.Result.BlockHeader.Timestamp)}})
  }
  return
}
