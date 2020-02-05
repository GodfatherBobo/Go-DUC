# Go-DUC
Dynamic DNS Update Client written in Golang.

# Build project
 1) Download main.go
 2) In terminal run go build -ldflags "-H windowsgui"

# Setup
1) Run the exe in setup via command line. use: -setup username password hostname
inputing your specific information in for username password and hostname

2) Install as a service use: -service install

Open Services and start the Dynamic DNS Update Client service. Confirmed to work on windows but it should work on linux and mac as well.



# Supported Dynamic DNS Providers:
No-IP.com
