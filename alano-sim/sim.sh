#!/bin/bash

user_name=("user11" "user12" "user13" "user14" "user15" "user21" "user22" "user23" "user24" "user25" "user11111" "user11112" "user11113" "user11114" "user11115" "user21111" "user21112" "user21113" "user21114" "user21115" "user31111" "user31112" "user31113" "user31114" "user31115")

user_total=${#user_name[*]}
echo $user_total

number_pool=(1 2 3 4 5 6)
number_total=${#number_pool[*]}
echo $number_total

for((i=1;i<=10000;i++));
do
tt=$(date +%Y%m%d%H%M%S)
hour=${tt:11:1}
echo $hour

user_index=($RANDOM%$user_total)
user_index1=($RANDOM%$user_total)
number_index=($RANDOM%$number_total)

parameter="'[\"${user_name[$user_index]}\",\"luckystareos\",\"3.0000 SYS\",\"0${number_pool[$number_index]}${user_name[$user_index1]}\"]'"
echo $parameter

eval cleos push action eosio.token transfer $parameter -p ${user_name[$user_index]}@active

if [ $hour == "0" ]
then
	echo "reveal"
	reveal_param="'[1]'"
	eval cleos push action luckystareos reveal $reveal_param -p luckystareos@active
else
	echo "betting"
fi

sleep 1m
done
