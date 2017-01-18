## run-patch.sh

This script is used to execute patch with fewer options on OWL system

```bash
run-patch.sh -bin=dbpatch -log-base=. -database=portal "-db-connection=db_user:pp@tcp(192.168.20.50:3306)"
```

## export-changelog.sh

This script is used to export the data of changelog(database patching) of a database.

The exported sql file could be used to combine with global script.

By global script with data of changlog, the database could be patched then even if it is imported with global script.

```bash
export-changelog.sh -db-name=falcon_portal -u db_user -ppp -h 192.168.20.50
```
