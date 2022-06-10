package main

import (
"strings"
"strconv"
"sort"
"context"
"encoding/json"
"time"
"github.com/gofiber/fiber/v2"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson"
)

func v1_xcash_dpops_unauthorized_stats(c *fiber.Ctx) error {

  // Variables
  var data_send string
  var database_data_statistics XcashDpopsStatisticsCollection
  var data_read_1 CurrentBlockHeight
  var output v1XcashDpopsUnauthorizedStats;
  var count int64
  var count5 int
  var count4 int
  var online_count int
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var total_voters int
  generated_supply := FIRST_BLOCK_MINING_REWARD + XCASH_PREMINE_TOTAL_SUPPLY
  
  // setup database
  collection_delegates := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("delegates")
  collection_statistics := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("statistics")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the current block Height
  data_send = send_http_data("http://127.0.0.1:18281/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"get_block_count"}`)
  if !strings.Contains(data_send, "\"result\"") {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  if err := json.Unmarshal([]byte(data_send), &data_read_1); err != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  
  // get the xcash dpops statistics
  err := collection_statistics.FindOne(ctx, bson.D{{}}).Decode(&database_data_statistics)
  if err == mongo.ErrNoDocuments {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  } else if err != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  
  // get total delegates
  count,err = collection_delegates.CountDocuments(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  output.TotalRegisteredDelegates = int(count)
  
  // get total delegates votes
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }
  
  count = 0
  online_count = 0
  for _, item := range mongo_results {
        count3,_ := strconv.ParseInt(item["total_vote_count"].(string),10,64)
        count += count3
        
        if item["online_status"].(string) == "true" {
          online_count++
        }
	}
	
	// get the circulating supply 
  for count5 = 2; count5 < data_read_1.Result.Count; count5++ {
    if count5 < XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT {
      generated_supply = generated_supply + (XCASH_TOTAL_SUPPLY - generated_supply) / XCASH_EMMISION_FACTOR
    } else {
      generated_supply += ((XCASH_TOTAL_SUPPLY - generated_supply) / XCASH_DPOPS_EMMISION_FACTOR)
    }
  }
  circulating_supply := int64(((generated_supply - (XCASH_PREMINE_TOTAL_SUPPLY - XCASH_PREMINE_CIRCULATING_SUPPLY)) * XCASH_WALLET_DECIMAL_PLACES_AMOUNT))
  
  // get the total voters
  total_voters = 0
  for count4 = 1; count4 < TOTAL_RESERVE_PROOFS_DATABASES; count4++ {
    count,_ := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("reserve_proofs_" + string(count4)).CountDocuments(ctx, bson.D{{}})
    total_voters += int(count)
  }

  // fill in the data
  output.MostTotalRoundsDelegateName = database_data_statistics.MostTotalRoundsDelegateName
  output.MostTotalRounds,_ = strconv.Atoi(database_data_statistics.MostTotalRounds)
  output.BestBlockVerifierOnlinePercentageDelegateName = database_data_statistics.BestBlockVerifierOnlinePercentageDelegateName
  output.BestBlockVerifierOnlinePercentage,_ = strconv.Atoi(database_data_statistics.BestBlockVerifierOnlinePercentage)
  output.MostBlockProducerTotalRoundsDelegateName = database_data_statistics.MostBlockProducerTotalRoundsDelegateName
  output.MostBlockProducerTotalRounds,_ = strconv.Atoi(database_data_statistics.MostBlockProducerTotalRounds)
  output.TotalVotes = count
  output.TotalVoters = total_voters
  output.AverageVote = int64(int64(count) / int64(total_voters))
  output.VotePercentage = int((float64(count) / float64(circulating_supply)) * 100)
  output.RoundNumber = data_read_1.Result.Count - XCASH_PROOF_OF_STAKE_BLOCK_HEIGHT
  output.TotalOnlineDelegates = online_count
  output.CurrentBlockVerifiersMaximumAmount = BLOCK_VERIFIERS_AMOUNT
  output.CurrentBlockVerifiersValidAmount = BLOCK_VERIFIERS_VALID_AMOUNT
    
  return c.JSON(output)
}

func v1_xcash_dpops_unauthorized_delegates_registered(c *fiber.Ctx) error {

  // Variables
  output:=[]*v1XcashDpopsUnauthorizedDelegatesBasicData{}
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection_delegates := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get total delegates votes
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the delegates registered"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the delegates registered"}
    return c.JSON(error)
  }
  
  for _, item := range mongo_results {
      // fill in the data
      data:=new(v1XcashDpopsUnauthorizedDelegatesBasicData)
      data.Votes,_ = strconv.ParseInt(item["total_vote_count"].(string),10,64)
      
      // get the total voters for the delegates
      total_voters := 0
      for count4 := 1; count4 < TOTAL_RESERVE_PROOFS_DATABASES; count4++ {
        count2,_ := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("reserve_proofs_" + string(count4)).CountDocuments(ctx, bson.D{{"public_address_voted_for",item["public_address"].(string)}})
        total_voters += int(count2)
      }
      data.Voters = total_voters
      
      data.IPAdress = item["IP_address"].(string)
      data.DelegateName = item["delegate_name"].(string)
      if item["shared_delegate_status"].(string) == "solo" {
        data.SharedDelegate = false
      } else {
        data.SharedDelegate = true
      }
      
      if strings.Contains(item["IP_address"].(string), ".xcash.foundation") {
        data.SeedNode = true
      } else {
        data.SeedNode = false
      }
      
      if strings.Contains(item["online_status"].(string), "true") {
        data.Online = true
      } else {
        data.Online = false
      }
      
      data.Fee,_ = strconv.Atoi(item["delegate_fee"].(string))
      data.TotalRounds,_ = strconv.Atoi(item["block_verifier_total_rounds"].(string))
      data.TotalBlockProducerRounds,_ = strconv.Atoi(item["block_producer_total_rounds"].(string))
      data.OnlinePercentage,_ = strconv.Atoi(item["block_verifier_online_percentage"].(string))
      
      output=append(output,data)
	}
	
	// sort the arrray by how xcash dpops sorts the delegates
	sort.Slice(output[:], func(i, j int) bool {
	    var count1 int
	    var count2 int
	    
	    // check if the delegate is a network data node
	    if output[i].IPAdress == "us1.xcash.foundation" {
	      count1 = 5
	    } else if output[i].IPAdress == "europe1.xcash.foundation" {
	      count1 = 4
	    } else if output[i].IPAdress == "europe2.xcash.foundation" {
	      count1 = 3
	    } else if output[i].IPAdress == "europe3.xcash.foundation" {
	      count1 = 2
	    } else if output[i].IPAdress == "oceania1.xcash.foundation" {
	      count1 = 1
	    } else {
	      count1 = 0      
	    }
	    
	    if output[j].IPAdress == "us1.xcash.foundation" {
	      count2 = 5
	    } else if output[j].IPAdress == "europe1.xcash.foundation" {
	      count2 = 4
	    } else if output[j].IPAdress == "europe2.xcash.foundation" {
	      count2 = 3
	    } else if output[j].IPAdress == "europe3.xcash.foundation" {
	      count2 = 2
	    } else if output[j].IPAdress == "oceania1.xcash.foundation" {
	      count2 = 1
	    } else {
	      count2 = 0      
	    }
	    
	    if count1 != count2 {
	      if count2 - count1 < 0 {
	        return true
	      } else {
	        return false
	      }
	    }
	    
	   // check if the delegate is online 
	    if output[i].Online != output[j].Online {
	      if output[i].Online == true {
	          return true
	      } else {
	          return false
	      }
	    }
        return output[i].Votes > output[j].Votes
    })
    
  return c.JSON(output)
}

func v1_xcash_dpops_unauthorized_delegates_online(c *fiber.Ctx) error {

  // Variables
  output:=[]*v1XcashDpopsUnauthorizedDelegatesBasicData{}
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection_delegates := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get total delegates votes
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the delegates online"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the delegates online"}
    return c.JSON(error)
  }
  
  for _, item := range mongo_results {
      if strings.Contains(item["online_status"].(string), "false") {
         continue;
      }
     
      // fill in the data
      data:=new(v1XcashDpopsUnauthorizedDelegatesBasicData)
      data.Votes,_ = strconv.ParseInt(item["total_vote_count"].(string),10,64)
      
      // get the total voters for the delegates
      total_voters := 0
      for count4 := 1; count4 < TOTAL_RESERVE_PROOFS_DATABASES; count4++ {
        count2,_ := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("reserve_proofs_" + string(count4)).CountDocuments(ctx, bson.D{{"public_address_voted_for",item["public_address"].(string)}})
        total_voters += int(count2)
      }
      data.Voters = total_voters
      
      data.IPAdress = item["IP_address"].(string)
      data.DelegateName = item["delegate_name"].(string)
      if item["shared_delegate_status"].(string) == "solo" {
        data.SharedDelegate = false
      } else {
        data.SharedDelegate = true
      }
      
      if strings.Contains(item["IP_address"].(string), ".xcash.foundation") {
        data.SeedNode = true
      } else {
        data.SeedNode = false
      }
      
      data.Online = true
      
      data.Fee,_ = strconv.Atoi(item["delegate_fee"].(string))
      data.TotalRounds,_ = strconv.Atoi(item["block_verifier_total_rounds"].(string))
      data.TotalBlockProducerRounds,_ = strconv.Atoi(item["block_producer_total_rounds"].(string))
      data.OnlinePercentage,_ = strconv.Atoi(item["block_verifier_online_percentage"].(string))
      
      output=append(output,data)
	}
	
	// sort the arrray by how xcash dpops sorts the delegates
	sort.Slice(output[:], func(i, j int) bool {
	    var count1 int
	    var count2 int
	    
	    // check if the delegate is a network data node
	    if output[i].IPAdress == "us1.xcash.foundation" {
	      count1 = 5
	    } else if output[i].IPAdress == "europe1.xcash.foundation" {
	      count1 = 4
	    } else if output[i].IPAdress == "europe2.xcash.foundation" {
	      count1 = 3
	    } else if output[i].IPAdress == "europe3.xcash.foundation" {
	      count1 = 2
	    } else if output[i].IPAdress == "oceania1.xcash.foundation" {
	      count1 = 1
	    } else {
	      count1 = 0      
	    }
	    
	    if output[j].IPAdress == "us1.xcash.foundation" {
	      count2 = 5
	    } else if output[j].IPAdress == "europe1.xcash.foundation" {
	      count2 = 4
	    } else if output[j].IPAdress == "europe2.xcash.foundation" {
	      count2 = 3
	    } else if output[j].IPAdress == "europe3.xcash.foundation" {
	      count2 = 2
	    } else if output[j].IPAdress == "oceania1.xcash.foundation" {
	      count2 = 1
	    } else {
	      count2 = 0      
	    }
	    
	    if count1 != count2 {
	      if count2 - count1 < 0 {
	        return true
	      } else {
	        return false
	      }
	    }
        return output[i].Votes > output[j].Votes
    })
    
  return c.JSON(output)
}

func v1_xcash_dpops_unauthorized_delegates_active(c *fiber.Ctx) error {

  // Variables
  output:=[]*v1XcashDpopsUnauthorizedDelegatesBasicData{}
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection_delegates := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get total delegates votes
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the delegates active"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the delegates active"}
    return c.JSON(error)
  }
  
  for _, item := range mongo_results {
      // fill in the data
      data:=new(v1XcashDpopsUnauthorizedDelegatesBasicData)
      data.Votes,_ = strconv.ParseInt(item["total_vote_count"].(string),10,64)
      
      // get the total voters for the delegates
      total_voters := 0
      for count4 := 1; count4 < TOTAL_RESERVE_PROOFS_DATABASES; count4++ {
        count2,_ := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("reserve_proofs_" + string(count4)).CountDocuments(ctx, bson.D{{"public_address_voted_for",item["public_address"].(string)}})
        total_voters += int(count2)
      }
      data.Voters = total_voters
      
      data.IPAdress = item["IP_address"].(string)
      data.DelegateName = item["delegate_name"].(string)
      if item["shared_delegate_status"].(string) == "solo" {
        data.SharedDelegate = false
      } else {
        data.SharedDelegate = true
      }
      
      if strings.Contains(item["IP_address"].(string), ".xcash.foundation") && item["IP_address"].(string) != "api.xcash.foundation" {
        data.SeedNode = true
      } else {
        data.SeedNode = false
      }
      
      if strings.Contains(item["online_status"].(string), "true") {
        data.Online = true
      } else {
        data.Online = false
      }
      
      data.Fee,_ = strconv.Atoi(item["delegate_fee"].(string))
      data.TotalRounds,_ = strconv.Atoi(item["block_verifier_total_rounds"].(string))
      data.TotalBlockProducerRounds,_ = strconv.Atoi(item["block_producer_total_rounds"].(string))
      data.OnlinePercentage,_ = strconv.Atoi(item["block_verifier_online_percentage"].(string))
      
      output=append(output,data)
	}
	
	// sort the arrray by how xcash dpops sorts the delegates
	sort.Slice(output[:], func(i, j int) bool {
	    var count1 int
	    var count2 int
	    
	    // check if the delegate is a network data node
	    if output[i].IPAdress == "us1.xcash.foundation" {
	      count1 = 5
	    } else if output[i].IPAdress == "europe1.xcash.foundation" {
	      count1 = 4
	    } else if output[i].IPAdress == "europe2.xcash.foundation" {
	      count1 = 3
	    } else if output[i].IPAdress == "europe3.xcash.foundation" {
	      count1 = 2
	    } else if output[i].IPAdress == "oceania1.xcash.foundation" {
	      count1 = 1
	    } else {
	      count1 = 0      
	    }
	    
	    if output[j].IPAdress == "us1.xcash.foundation" {
	      count2 = 5
	    } else if output[j].IPAdress == "europe1.xcash.foundation" {
	      count2 = 4
	    } else if output[j].IPAdress == "europe2.xcash.foundation" {
	      count2 = 3
	    } else if output[j].IPAdress == "europe3.xcash.foundation" {
	      count2 = 2
	    } else if output[j].IPAdress == "oceania1.xcash.foundation" {
	      count2 = 1
	    } else {
	      count2 = 0      
	    }
	    
	    if count1 != count2 {
	      if count2 - count1 < 0 {
	        return true
	      } else {
	        return false
	      }
	    }
	    
	   // check if the delegate is online 
	    if output[i].Online != output[j].Online {
	      if output[i].Online == true {
	          return true
	      } else {
	          return false
	      }
	    }
        return output[i].Votes > output[j].Votes
    })
    
    // only return the top 50
    output = output[0:BLOCK_VERIFIERS_AMOUNT]
    
  return c.JSON(output)
}

func v1_xcash_dpops_unauthorized_delegates(c *fiber.Ctx) error {

  // Variables
  output_data:=[]*v1XcashDpopsUnauthorizedDelegatesBasicData{}
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var delegate string
  var database_data_delegates XcashDpopsDelegatesCollection
  var output v1XcashDpopsUnauthorizedDelegatesData;
  var total_voters int
  
  delegate = c.Params("delegateName")
  if delegate == "" {
    error := ErrorResults{"Could not get the delegates data"}
    return c.JSON(error)
  }
  
  // setup database
  collection_delegates := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the delegate
  err := collection_delegates.FindOne(ctx, bson.D{{"delegate_name",delegate}}).Decode(&database_data_delegates)
  if err == mongo.ErrNoDocuments {
    error := ErrorResults{"Could not get the delegates data"}
    return c.JSON(error)
  } else if err != nil {
    error := ErrorResults{"Could not get the delegates data"}
    return c.JSON(error)
  }
  
  if (database_data_delegates.OnlineStatus == "true") {
     output.Online = true 
  } else {
     output.Online = false
  }
  
  // get total voters
  for count4 := 1; count4 < TOTAL_RESERVE_PROOFS_DATABASES; count4++ {
    count2,_ := mongoClient.Database("XCASH_PROOF_OF_STAKE").Collection("reserve_proofs_" + string(count4)).CountDocuments(ctx, bson.D{{"public_address_voted_for",database_data_delegates.PublicAddress}})
    total_voters += int(count2)
  }
  
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the delegates data"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the delegates data"}
    return c.JSON(error)
  }
  
  for _, item := range mongo_results {
      // fill in the data
      data:=new(v1XcashDpopsUnauthorizedDelegatesBasicData)
      data.DelegateName = item["delegate_name"].(string)
      data.Votes,_ = strconv.ParseInt(item["total_vote_count"].(string),10,64)
      data.IPAdress = item["IP_address"].(string)
      if strings.Contains(item["online_status"].(string), "true") {
        data.Online = true
      } else {
        data.Online = false
      }
      output_data=append(output_data,data)
	}
	
	// sort the arrray by how xcash dpops sorts the delegates
	sort.Slice(output_data[:], func(i, j int) bool {
	    var count1 int
	    var count2 int
	    
	    // check if the delegate is a network data node
	    if output_data[i].IPAdress == "us1.xcash.foundation" {
	      count1 = 5
	    } else if output_data[i].IPAdress == "europe1.xcash.foundation" {
	      count1 = 4
	    } else if output_data[i].IPAdress == "europe2.xcash.foundation" {
	      count1 = 3
	    } else if output_data[i].IPAdress == "europe3.xcash.foundation" {
	      count1 = 2
	    } else if output_data[i].IPAdress == "oceania1.xcash.foundation" {
	      count1 = 1
	    } else {
	      count1 = 0      
	    }
	    
	    if output_data[j].IPAdress == "us1.xcash.foundation" {
	      count2 = 5
	    } else if output_data[j].IPAdress == "europe1.xcash.foundation" {
	      count2 = 4
	    } else if output_data[j].IPAdress == "europe2.xcash.foundation" {
	      count2 = 3
	    } else if output_data[j].IPAdress == "europe3.xcash.foundation" {
	      count2 = 2
	    } else if output_data[j].IPAdress == "oceania1.xcash.foundation" {
	      count2 = 1
	    } else {
	      count2 = 0      
	    }
	    
	    if count1 != count2 {
	      if count2 - count1 < 0 {
	        return true
	      } else {
	        return false
	      }
	    }
	    
	   // check if the delegate is online 
	    if output_data[i].Online != output_data[j].Online {
	      if output_data[i].Online == true {
	          return true
	      } else {
	          return false
	      }
	    }
        return output_data[i].Votes > output_data[j].Votes
    })
    
    // get the Rank
    for rank, item := range output_data {
        if item.DelegateName == database_data_delegates.DelegateName {
           output.Rank = rank+1
           break
        }
    }
    
    if database_data_delegates.SharedDelegateStatus == "solo" {
        output.SharedDelegate = false
      } else {
        output.SharedDelegate = true
      }
      
      if strings.Contains(database_data_delegates.IPAddress, ".xcash.foundation") && database_data_delegates.IPAddress != "api.xcash.foundation" {
        output.SeedNode = true
      } else {
        output.SeedNode = false
      }

  // fill in the data
  output.Votes,_ = strconv.ParseInt(database_data_delegates.TotalVoteCount,10,64)
  output.Voters = int(total_voters)
  output.IPAdress = database_data_delegates.IPAddress
  output.DelegateName = database_data_delegates.DelegateName
  output.PublicAddress = database_data_delegates.PublicAddress
  output.About = database_data_delegates.About
  output.Website = database_data_delegates.Website
  output.Team = database_data_delegates.Team
  output.Specifications = database_data_delegates.ServerSpecs
  output.Fee,_ = strconv.Atoi(database_data_delegates.DelegateFee)
  output.TotalRounds,_ = strconv.Atoi(database_data_delegates.BlockVerifierTotalRounds)
  output.TotalBlockProducerRounds,_ = strconv.Atoi(database_data_delegates.BlockProducerTotalRounds)
  output.OnlinePercentage,_ = strconv.Atoi(database_data_delegates.BlockVerifierOnlinePercentage)
    
  return c.JSON(output)
}
