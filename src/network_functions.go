package main

import (
"bytes"
"io/ioutil"
"net/http"
"time"
"errors"
)

func send_http_data(url string,data string) (string, error) {
  // create the http request
  req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
  if err != nil {
    return "",errors.New("")
  }

  // set the request headers  
  req.Header.Set("Content-Type", "application/json")
 
  // set the http client settings
  client := &http.Client{}
  client.Timeout = time.Second * 2

  // send the request
  resp, err := client.Do(req)
  if err != nil {
    return "",errors.New("")
  }

  // close the connection
  defer resp.Body.Close()
  
  // get the response body
  body, _ := ioutil.ReadAll(resp.Body)

  return string(body),nil
}
