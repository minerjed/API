# API
API using the api.xcash.foundation

Please read the [documentation](https://docs.xcash.foundation/api/get-started) to use the API

The API covers:  
Blockchain  
DPOPS  
Namespace  
Xpayment  
Xpayment Twitter


# How to build from source

install go  
download the latest go version from https://go.dev/doc/install
 
untar it  
`tar -xf go* && rm go*.tar.gz && mv go /usr/local/`
 
edit the path  
`sudo nano ~/.profile`
 
add this line to the end of the file  
`export PATH=$PATH:/usr/local/go/bin`
 
save the file  
`source ~/.profile`
 
verify the install  
`go version`

Install mongo  
```
wget -qO - https://www.mongodb.org/static/pgp/server-5.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.4 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.4.list
sudo apt update
sudo apt install mongodb-org
```

Install  
`git clone https://github.com/X-CASH-official/turbotx-backend.git && cd API`

copy the systemd file  
`cp -a API.service /lib/systemd/system/ && sudo systemctl daemon-reload`

Build the program  
`make clean ; make release`

Run the program  
`systemctl start API`
