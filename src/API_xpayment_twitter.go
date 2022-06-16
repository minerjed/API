package main

import (
"context"
"time"
"github.com/gofiber/fiber/v2"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/bson"
)

func v1_xpayment_twitter_unauthorized_stats(c *fiber.Ctx) error {

  // Variables
  var output v1XpaymentTwitterUnauthorizedStats
  var count int64
  var mongo_sort *mongo.Cursor
  var mongo_results []bson.M
  var error error
  
  var ttg_public int = 0
  var ttg_private int = 0
  var tt24_public int = 0
  var tt24_private int = 0
  var tt1_public int = 0
  var tt1_private int = 0

  var tvg_public int64 = 0
  var tvg_private int64 = 0
  var tv24_public int64 = 0
  var tv24_private int64 = 0
  var tv1_public int64 = 0
  var tv1_private int64 = 0
  
  var settings int
  var amount int64
  var time_settings int
  time_current := time.Now().UTC().Unix()
  

  // setup database
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  // get how many users
  count,error = mongoClient.Database(XPAYMENT_TWITTER_DATABASE).Collection("userWallets").CountDocuments(ctx, bson.D{{}})
  if error != nil {
    error := ErrorResults{"Could not get the xpayment twitter statistics"}
    return c.JSON(error)
  }
  output.TotalUsers = int(count)
  
  // get how many deposits and withdraws
  count,error = mongoClient.Database(XPAYMENT_TWITTER_DATABASE).Collection("doneDeposits").CountDocuments(ctx, bson.D{{}})
  if error != nil {
    error := ErrorResults{"Could not get the xpayment twitter statistics"}
    return c.JSON(error)
  }
  output.TotalDeposits = int(count)
  
  // get how many deposits and withdraws
  count,error = mongoClient.Database(XPAYMENT_TWITTER_DATABASE).Collection("doneWithdrawals").CountDocuments(ctx, bson.D{{}})
  if error != nil {
    error := ErrorResults{"Could not get the xpayment twitter statistics"}
    return c.JSON(error)
  }
  output.TotalWithdraws = int(count)
  
  // get the payment details 
  mongo_sort, error = mongoClient.Database(XPAYMENT_TWITTER_DATABASE).Collection("twitterHistory").Find(ctx, bson.D{{}})
  if error != nil {
    error := ErrorResults{"Could not get the xpayment twitter statistics"}
    return c.JSON(error)
  }

  if error = mongo_sort.All(ctx, &mongo_results); error != nil {
    error := ErrorResults{"Could not get the xcash dpops statistics"}
    return c.JSON(error)
  }

  for _, item := range mongo_results {
      
      // convert the types
      switch v := item["type"].(type) {
	case int32:
		settings = int(v)
	case int64:
		settings = int(v)
    case float32:
		settings = int(v)
	case float64:
		settings = int(v)
	}
	
	 switch v1 := item["amount"].(type) {
	case int32:
		amount = int64(v1)
	case int64:
		amount = int64(v1)
    case float32:
		amount = int64(v1)
	case float64:
		amount = int64(v1)
	}
	
	 switch v2 := item["time"].(type) {
	case int32:
		time_settings = int(v2)
	case int64:
		time_settings = int(v2)
    case float32:
		time_settings = int(v2)
	case float64:
		time_settings = int(v2)
	}
	
	
        if settings == 1 {
             ttg_private++
	         tvg_private += amount

	
	if (int64(time_settings) + time_day) > time_current {
	    tt24_private++
	    tv24_private += amount
	}
	
	if (int64(time_settings) + time_hour) > time_current {
	    tt1_private++
	    tv1_private += amount
	}
        }
        
        
        if settings == 0 {
             ttg_public++
	         tvg_public += amount

	
	if (int64(time_settings) + time_day) > time_current {
	    tt24_public++
	    tv24_public += amount
	}
	
	if (int64(time_settings) + time_hour) > time_current {
	    tt1_public++
	    tv1_public += amount
	}
        }
	}
	
	


  // fill in the data
  output.AvgTipAmount = int((float64(tvg_public + tvg_private) / float64(ttg_public + ttg_private)))
  output.TotalTipsPublic = ttg_public
  output.TotalTipsPrivate = ttg_private
  output.TotalVolumeSentPublic = tvg_public
  output.TotalVolumeSentPrivate = tvg_private
  output.TotalTipsLastDayPublic = tt24_public
  output.TotalTipsLastDayPrivate = tt24_private
  output.TotalVolumeSentLastDayPublic = tv24_public
  output.TotalVolumeSentLastDayPrivate = tv24_private
  output.TotalTipsLastHourPublic = tt1_public
  output.TotalTipsLastHourPrivate = tt1_private
  output.TotalVolumeSentLastHourPublic = tv1_public
  output.TotalVolumeSentLastHourPrivate = tv1_private

  return c.JSON(output)
}
