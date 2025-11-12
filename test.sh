#!/bin/bash

# build 

mage all

# run server

./minicoin-server --fail-in 6 2> log.txt &

# add initial transactions
./client -id 1 --action 0 --value 20 
./client -id 2 --action 0 --value 30 
./client -id 3 --action 0 --value 10 

#consult balance
./client -id 1 --action 2
./client -id 2 --action 2
./client -id 3 --action 2

# make changes to balance
./client -id 1 --action 0 --value 30 
./client -id 2 --action 0 --value -31 
./client -id 3 --action 0 --value -5 

#consult balance
./client -id 1 --action 2 
./client -id 2 --action 2  
./client -id 3 --action 2 

# check balance of unknown client

./client -id 4 --action 2 

# draw from unknown client

./client -id 4 --action 0 --value -30 

# check blockchain integrity

./client -id 1 --action 1 

# make transaction to meet fail-in criteria

./client -id 1 --action 0 --value 30 

# check blockchain integrity

./client -id 1 --action 1 

# make transaction after blockchain corruption 

./client -id 1 --action 0 --value 30 

kill `pgrep minicoin-server`
