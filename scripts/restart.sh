#!/bin/sh 

function func(){
	killall -9 $1

	killall -0 $1
	while [ $? -ne 1 ]; do
		sleep 1
		killall -0 $1
	done
}

if [ $# -eq 0 ]
	then
		func GateServer
		func LoginServer
		func GameServer
		func WorldServer
		func DBServer
		func LogServer
	else
		func $1
fi

$GOGAMESERVER_PATH/bin/LogServer &
sleep 1
$GOGAMESERVER_PATH/bin/DBServer &
sleep 1
$GOGAMESERVER_PATH/bin/GateServer &
sleep 1
$GOGAMESERVER_PATH/bin/LoginServer &
sleep 1
$GOGAMESERVER_PATH/bin/WorldServer &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=1 &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=2 &
sleep 1
$GOGAMESERVER_PATH/bin/GameServer -s=3 &
