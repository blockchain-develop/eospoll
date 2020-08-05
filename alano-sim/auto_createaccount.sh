#!/bin/bash

user_name=("user11" "user12" "user13" "user14" "user15" "user21" "user22" "user23" "user24" "user25" "user31" "user32" "user33" "user34" "user35" "user41" "user42" "user43" "user44" "user45" "user51" "user52" "user53" "user54" "user55")

user_total=${#user_name[*]}
echo $user_total

for((i=0;i<$user_total;i++));
do

eval cleos create account eosio ${user_name[$i]} EOS84of4wKVARqDGwmGjrSuF26seeBWxwaaseTHk2UtwaBR8mRdSL EOS84of4wKVARqDGwmGjrSuF26seeBWxwaaseTHk2UtwaBR8mRdSL -p eosio@active


parameter="'[\"${user_name[$i]}\",\"10000.0000 EOS\",\"memo\"]'"
eval cleos push action eosio.token issue $parameter -p eosio@active

sleep 5s
done
