package main

import (
"strings"
"strconv"
"regexp"
"sort"
"context"
"time"
"github.com/gofiber/fiber/v2"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson"
)

func v1_xcash_namespace_unauthorized_stats(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedStats
  var total_names_registered int
  var total_amount int64
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data_delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  mongo_sort, err = collection.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the namespace statistics"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the namespace statistics"}
    return c.JSON(error)
  }
  
  total_names_registered = 0
  total_amount = 0
  
  for _, item := range mongo_results {
    total,_ := strconv.Atoi(item["total_registered_renewed_amount"].(string))
    total_names_registered += total
    totaldata,_ := strconv.ParseInt(item["total_amount"].(string), 10, 64)
    total_amount += totaldata
  }
  
  // fill in the data
  output.TotalNamesRegisteredOrRenewed = total_names_registered
  output.TotalVolume = total_amount
  
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_delegates_registered(c *fiber.Ctx) error {

  output:=[]*v1XcashNamespaceUnauthorizedDelegatesRegistered{}
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection_delegates := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data_delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the delegates
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
      data:=new(v1XcashNamespaceUnauthorizedDelegatesRegistered)
      data.DelegateName = item["name"].(string)
      data.PublicAddress = item["public_address"].(string)
      data.Amount,_ = strconv.ParseInt(item["amount"].(string),10,64)
      output=append(output,data)
	}
	
	// sort the arrray by the amount
	sort.Slice(output[:], func(i, j int) bool {
        return output[i].Amount > output[j].Amount
    })
    
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_delegates_delegate_name(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedDelegatesDelegateName
  var database_data XcashDpopsRemoteDataDelegatesCollection
  var delegate string
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data_delegates")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if delegate = c.Params("delegateName"); delegate == "" {
    error := ErrorResults{"Could not get the delegates details"}
    return c.JSON(error)
  }
  
  // get the delegates data
    err := collection.FindOne(ctx, bson.D{{"name", delegate}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      error := ErrorResults{"Could not get the delegates details"}
    return c.JSON(error)
    } else if err != nil {
      error := ErrorResults{"Could not get the delegates details"}
    return c.JSON(error)
    }
  
  // fill in the data
  output.DelegateName = database_data.Name
  output.PublicAddress = database_data.PublicAddress
  output.Amount,_ = strconv.ParseInt(database_data.Amount, 10, 64)
  output.TotalNamesRegisteredOrRenewed,_ = strconv.Atoi(database_data.TotalRegisteredRenewedAmount)
  output.TotalVolume,_ = strconv.ParseInt(database_data.TotalAmount, 10, 64)
  
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_names_name(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedNamesName
  var database_data XcashDpopsRemoteDataCollection
  var name string
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if name = c.Params("name"); name == "" {
    error := ErrorResults{"Could not get the name details"}
    return c.JSON(error)
  }
  
  // get the delegates data
    err := collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      error := ErrorResults{"Could not get the name details"}
    return c.JSON(error)
    } else if err != nil {
      error := ErrorResults{"Could not get the name details"}
    return c.JSON(error)
    }
  
  // fill in the data
  output.Address = database_data.Address
  output.Saddress = database_data.Saddress
  output.Paddress = database_data.Paddress
  output.Expires,_ = strconv.Atoi(database_data.Timestamp)
  output.DelegateName = database_data.ReserveDelegateAddress
  output.DelegateAmount,_ = strconv.ParseInt(database_data.ReserveDelegateAmount, 10, 64)
  
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_names_status_name(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedNamesStatusName
  var database_data XcashDpopsRemoteDataCollection
  var name string
  var valid_name = regexp.MustCompile(VALID_NAME_DATA).MatchString
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if name = c.Params("name"); name == "" {
    error := ErrorResults{"Could not get the name status"}
    return c.JSON(error)
  }
  
  // check if the name is valid
  if !valid_name(name) {
      error := ErrorResults{"Could not get the name status"}
    return c.JSON(error)
  }
  
  // get the delegates data
    err := collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      output.Status = true
    return c.JSON(output)
    } else if err != nil {
      error := ErrorResults{"Could not get the name details"}
      return c.JSON(error)
    }
    
    if database_data.Timestamp == REMOTE_DATA_TIMESTAMP_DEFAULT_AMOUNT && database_data.TxHash == "" {
      output.Status = true
    } else {
      output.Status = false
    }
  
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_names_status_address(c *fiber.Ctx) error {

  var output v1XcashNamespaceUnauthorizedNamesStatusAddress
  var address string
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var err error
  
  // setup database
  collection_delegates := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if address = c.Params("address"); address == "" {
    error := ErrorResults{"Could not get the address status"}
    return c.JSON(error)
  }
  
  // get the remote data
  mongo_sort, err = collection_delegates.Find(ctx, bson.D{{}})
  if err != nil {
    error := ErrorResults{"Could not get the address status"}
    return c.JSON(error)
  }
  
  if err = mongo_sort.All(ctx, &mongo_results); err != nil {
    error := ErrorResults{"Could not get the address status"}
    return c.JSON(error)
  }
  
  for _, item := range mongo_results {
      if item["address"].(string) == address {
        output.Status = "address"
        return c.JSON(output)
      } else if item["saddress"].(string) == address {
        output.Status = "saddress"
        return c.JSON(output)
      } else if item["paddress"].(string) == address {
        output.Status = "paddress"
        return c.JSON(output)
      } else if strings.Contains(item["saddress_list"].(string), address) {
        output.Status = "saddress"
        return c.JSON(output)
      } else if strings.Contains(item["paddress_list"].(string), address) {
        output.Status = "paddress"
        return c.JSON(output)
      }
	}
	
	output.Status = "not registered"
	
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_names_convert_name(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedNamesConvertName
  var database_data XcashDpopsRemoteDataCollection
  var name string
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if name = c.Params("name"); name == "" {
    error := ErrorResults{"Could not convert the name to an address"}
    return c.JSON(error)
  }
  
  // get the delegates data
    err := collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
       error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    return c.JSON(output)
    } else if err != nil {
      error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    }
    
    // fill in the data
  output.Address = database_data.Address
  output.Saddress = database_data.Saddress
  output.Paddress = database_data.Paddress
  
  return c.JSON(output)
}

func v1_xcash_namespace_unauthorized_names_convert_address(c *fiber.Ctx) error {

  // Variables
  var output v1XcashNamespaceUnauthorizedNamesConvertAddress
  var database_data XcashDpopsRemoteDataCollection
  var address string
  
  // setup database
  collection := mongoClient.Database(XCASH_NAMESPACE_DATABASE).Collection("remote_data")
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get the resource
  if address = c.Params("address"); address == "" {
    error := ErrorResults{"Could not convert the name to an address"}
    return c.JSON(error)
  }
  
  // get the delegates data
    err := collection.FindOne(ctx, bson.D{{"address", address}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      goto SADDRESS
    } else if err != nil {
      error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    }
    
    output.Name = database_data.Name
    output.Extension = "address"
    return c.JSON(output)
    
    SADDRESS:
    
    // get the delegates data
    err = collection.FindOne(ctx, bson.D{{"saddress", address}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      goto PADDRESS
    } else if err != nil {
      error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    }
    
    output.Name = database_data.Name
    output.Extension = "saddress"
    return c.JSON(output)
    
    PADDRESS:
    
    // get the delegates data
    err = collection.FindOne(ctx, bson.D{{"paddress", address}}).Decode(&database_data)
    if err == mongo.ErrNoDocuments {
      error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    } else if err != nil {
      error := ErrorResults{"Could not convert the name to an address"}
      return c.JSON(error)
    }
    
    output.Name = database_data.Name
    output.Extension = "paddress"
    return c.JSON(output)
}
