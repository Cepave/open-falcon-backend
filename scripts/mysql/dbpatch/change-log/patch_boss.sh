#!/bin/bash

./dbpatch -driverName="mysql" -dataSourceName="root:password@tcp(10.20.30.40:3306)/boss" '-changeLog'=change-log-boss.yaml '-patchFileBase=schema-boss'

