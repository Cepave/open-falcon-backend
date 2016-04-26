# (This document is deprecated)

# nqm

The module of network quality measurement

	
# 測試環境設定方法

## 清空 graph 模組的資料

* 先確認 nqm 的 scripts 沒有在跑，有的話要 kill 掉
* 在 Vagrant 的 single_instance 環境中：

		$ docker stop graph
		$ cd /home/openfalcon/graph
		$ sudo rm -rf ./*

## 重置 graph 這個資料庫的資料表

	$ mysql -h 10.20.30.40 -u root -ppassword -D graph
	mysql> truncate endpoint;
	mysql> truncate endpoint_counter;
	mysql> truncate tag_ endpoint;
	mysql> exit

## 重啟 graph 的 container

	$ docker start graph
	
## 重置 falcon_portal 這個資料庫

	$ mysql -h 127.0.0.1 -u root -ppassword
	mysql> drop database falcon_portal;
	mysql> create database falcon_portal;
	mysql> exit

## 執行 db-patch 確保資料庫的 schema 不是舊的

	$ cd workspace/src/github.com/Cepave/scripts
	$ go get gopkg.in/yaml.v2
	$ go build -o db-patch ./dbpatch/go
	$ ./db-patch -driverName=mysql -dataSourceName="root:password@tcp(10.20.30.40:3306)/falcon_portal" -changeLog=dbpatch/changchange-log-portal.yaml -patchFileBase=dbpatch/change-log/schema-portal

## 目前 graph 的 counters 有中文的話會出問題，因此把 falcon_portal 資料庫的一些資料都改成英文拼音

	$ cd workspace/src/github.com/Cepave/nqm
	$ mysql -h 10.20.30.40 -u root -ppassword -D falcon_portal < ./test/sql/dbSchema.sql

## 先跑兩個 nqm 模組的 instance 各一次，讓 falcon_portal 資料庫的 nqm_agent 資料表多兩個 rows

	$ ./nqm
	$ ./nqm -connectionId nqm-agent-2@10.20.30.40

## 加上測資

	$ mysql -h 10.20.30.40 -u root -ppassword -D falcon_portal < ./test/sql/testData.sql

## 執行

	$ cd test/run
	$ ./run1.sh&
	$ ./run2.sh&
