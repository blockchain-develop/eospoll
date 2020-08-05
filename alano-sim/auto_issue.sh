#!/bin/bash

user_name=("user1111" "user1112" "user1113" "user1114" "user1115" "user2111" "user2112" "user2113" "user2114" "user2115" "user3111" "user3112" "user3113" "user3114" "user3115" "user4111" "user4112" "user4113" "user4114" "user4115" "user5111" "user5112" "user5113" "user5114" "user5115")

user_total=${#user_name[*]}
echo $user_total

for((i=0;i<$user_total;i++));
do

parameter="'[\"${user_name[$i]}\",\"10000.0000 EOS\",\"memo\"]'"
eval cleos push action eosio.token issue $parameter -p eosio@active

sleep 10s
done
