#How to patch

## 先將 mysql 導入 production environment 的資料庫。
``` mysql -h 10.20.30.40 -uroot -p < ./2016-XXXX.sql ```

## 編繹 db-patch 這隻程式
```
cd mysql/dbpatch/go
go build -o db-patch
cp db-patch ../..
```

## 開始 patch
```
cd mysql/dbpatch/change-log
../../db-patch -driverName="mysql" -dataSourceName="root:password@tcp(10.20.30.40:3306)/falcon_portal" '-changeLog=change-log-portal.yaml'
``` 
