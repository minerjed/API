package main

import (
"fmt"
"bytes"
"io/ioutil"
"net/http"
"time"
)

func send_http_data(url string,data string) string {
  // create the http request
  req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
  if err != nil {
    return "error1"
  }

  // set the request headers  
  req.Header.Set("Content-Type", "application/json")
 
  // set the http client settings
  client := &http.Client{}
  client.Timeout = time.Second * 2

  // send the request
  resp, err := client.Do(req)
  if err != nil {
    return "error1"
  }

  // close the connection
  defer resp.Body.Close()
  
  // get the response body
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Printf("for %s sending %s received %s\n", url,data,body)

  return string(body)
}
