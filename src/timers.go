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
    for {
        if time.Now().Minute() % 5 == 0 && time.Now().Second() % 60 == 0 {
            process_block_data()
            time.Sleep(1 * time.Second)
        }
        time.Sleep(1 * time.Second)
    }
}

func timers_build_data() {
    block_height := 800000
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
    _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").InsertOne(ctx, bson.D{{"public", "0"}, {"private", "0"}})
    for {
            fmt.Printf("Processing block: %d\n",block_height)
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
  var private_tx_count int = 0
  var data_send string
  var data_read_1 TxData
  var error error
  var database_data XcashAPIStatisticsCollection
  var data_read_2 CheckTxKey
  var amount int64
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
  
  // get the currrent tx count 
  err := mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").FindOne(ctx, bson.D{{}}).Decode(&database_data)
  if err == mongo.ErrNoDocuments {
    return
  } else if err != nil {
    return
  }
  
  public_tx_count,_ = strconv.Atoi(database_data.Public)
  private_tx_count,_ = strconv.Atoi(database_data.Private)
  
  // parse the reserve bytes 
  delegate = s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(s, BLOCKCHAIN_DATA_SEGMENT_STRING)]
  delegate_name_data,_ := hex.DecodeString(delegate)
  delegate = string(delegate_name_data)
    
  length = len(s) - (len(s) - strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START_DATA)) - 106 - 142

  timestamp = int(varint_decode(s[4:14]))
  reward = varint_decode(s[106 : 106+length])
  
  tx_data := s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_END)+len(BLOCKCHAIN_RESERVED_BYTES_END)+4:]
  if len(tx_data) > TRANSACTION_HASH_LENGTH {
      for len(tx_data) >= TRANSACTION_HASH_LENGTH {
      tx := tx_data[0:TRANSACTION_HASH_LENGTH]
      tx_data = tx_data[TRANSACTION_HASH_LENGTH:]
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
  }
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"public", strconv.Itoa(public_tx_count)}}}})
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"private", strconv.Itoa(private_tx_count)}}}})
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").InsertOne(ctx, bson.D{{"height", strconv.Itoa(block_height)}, {"delegate", delegate},{"reward", strconv.FormatInt(reward, 10)},{"time", strconv.Itoa(timestamp)}})
    
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
  var private_tx_count int = 0
  var data_send string
  var data_read_1 TxData
  var error error
  var database_data XcashAPIStatisticsCollection
  var data_read_2 CheckTxKey
  var amount int64
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the reserve bytes
  s = get_reserve_bytes(block_height)
  if s == "" {
    return
  }
  
  // get the currrent tx count 
  err := mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").FindOne(ctx, bson.D{{}}).Decode(&database_data)
  if err == mongo.ErrNoDocuments {
    return
  } else if err != nil {
    return
  }
  
  public_tx_count,_ = strconv.Atoi(database_data.Public)
  private_tx_count,_ = strconv.Atoi(database_data.Private)
  
  // parse the reserve bytes 
  delegate = s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START)+len(BLOCKCHAIN_RESERVED_BYTES_START):strings.Index(s, BLOCKCHAIN_DATA_SEGMENT_STRING)]
  delegate_name_data,_ := hex.DecodeString(delegate)
  delegate = string(delegate_name_data)
    
  length = len(s) - (len(s) - strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_START_DATA)) - 106 - 142

  timestamp = int(varint_decode(s[4:14]))
  reward = varint_decode(s[106 : 106+length])
  
  tx_data := s[strings.Index(s, BLOCKCHAIN_RESERVED_BYTES_END)+len(BLOCKCHAIN_RESERVED_BYTES_END)+4:]
  if len(tx_data) > TRANSACTION_HASH_LENGTH {
      for len(tx_data) >= TRANSACTION_HASH_LENGTH {
      tx := tx_data[0:TRANSACTION_HASH_LENGTH]
      tx_data = tx_data[TRANSACTION_HASH_LENGTH:]
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
  }
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"public", strconv.Itoa(public_tx_count)}}}})
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("statistics").UpdateOne(ctx, bson.D{{}},bson.D{{"$set", bson.D{{"private", strconv.Itoa(private_tx_count)}}}})
  _,_ = mongoClient.Database(XCASH_API_DATABASE).Collection("blocks").InsertOne(ctx, bson.D{{"height", strconv.Itoa(block_height)}, {"delegate", delegate},{"reward", strconv.FormatInt(reward, 10)},{"time", strconv.Itoa(timestamp)}})
 
  return
}
