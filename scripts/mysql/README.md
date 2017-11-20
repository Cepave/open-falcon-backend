# Patching

* [liquibase/](liquibase/) - The runtime files of Liquibase
* [liquibase-changelog/](liquibase-changelog/) - The files of change log for Liquibase

## Liquibase

[Liquibase](http://www.liquibase.org/documentation/index.html) is comprehensive utility to perform incremental upgrading on RDB.

## Scripts

* `patch-owl-database.sh` - Use `patch-owl-database.sh --help` to see help messages
```
./patch-owl-databases.sh [--mysql_conn=192.168.20.50:3306] [--command=update] [--options=<args>] [--database=<database>] [--change-log-base=<directory>] [--prefix=<db prefix>] [--suffix=<db suffix>]

This script would update databases(By Liquibase): "uic falcon_portal falcon_links grafana graph boss dashboard imdb"

	--mysql_conn=<host:port> - host and port for database connection of MySql.
		For example: 192.168.20.50:3306(default value)
	--command=<command> - The command of Liquibase
		For example: "update", "updateCount 13"
	--database=<database> - Name of OWL database.
		Value domain: "uic falcon_portal falcon_links grafana graph boss dashboard imdb"
	--options=<args> - Arguments of liquibase.
		For example: "--username=abc --password=cepave"
	--change-log-base=<directory> - Base directory of files for changelog
		Default value: "./liquibase-changelog/"
	--prefix=<prefix> - The prefix to be added to name of database
	--suffix=<suffix> - The suffix to be appended to name of database
	--help - Show this message
```
* `recreate-owl-databases.sh` - Use `patch-owl-database.sh --help` to see help messages
```
./recreate-owl-databases.sh [--mysql=<args>] [--prefix=<db prefix>] [--suffix=<db suffix>] [--action=<action>]

This script would drop and create databases: "imdb falcon_portal uic falcon_links grafana graph boss dashboard"

	--action=<action> - The action to be perfomed
		Value domain: recreate(default), drop, create(if not existing)
	--prefix=<prefix> - The prefix to be added to name of database
	--suffix=<suffix> - The suffix to be appended to name of database
	--mysql=<args> - The arguments to be fed to "mysql" command
	--help - Show this message
```

# Dependencies of OWL Databases

The dependencies of OWL databases is depicted by following graph:
```
+-----+     +---------------+     +------+
| uic +-----> falcon_portal +-----> imdb |
+-----+     +---------------+     +------+
```

The **recreate-owl-databases.sh** and **patch-owl-databases.sh** script complies with the dependencies.

# Old patching(Deprecated)

**Global schema**: Which is put into [db_schema](db_schema/) directory, these files won't be maintinaed anymore.

In [dbpatch/](dbpatch/), there is a simple patching tool(written by GoLang) which is put in [dbpatch/tool](dbpatch/tool/) directory.

* [change-log/](dbpatch/change-log) - Static files of change log of databases
* [go/](dbpatch/go) - Go source code of patching
* [tool/](dbpatch/tool) - Out-of-box utilities for patching

## Tools

* `dbpatch-$OS-$ARCH` files - Corresponding to OS and arch, these executives are used to parse files of change log and to perform operations on database. You could execute the file to see help message.
	```
	dbpatch -driverName=<driver_name> -dataSourceName=<data_srouce_name> [-changeLog=change-log.json] [-patchFileBase=patch-files] ["-dilimeter=;"]

	-changeLog string
		The file of change log (default "change-log.json")
	-dataSourceName string
		The name of data source
	-delimiter string
		The directory of patch files (default ";")
	-driverName string
		The name of driver
	-help
		Show help
	-patchFileBase string
		The directory of patch files (default "patch-files")
	```
* `run-owl-patch.sh` - Executes patch to any of databases in OWL system,
	```
	Options:
	-bin=<bin file> - default value is "dbpatch"
	-bin-os=<os for bin file> - If this value is provided, ignores "-bin" option
		linux, windows, or osx
	-bin-arch=<32> - If this value is provided, ignores "-bin" option
		64 or 32
	-database=<boss|portal|uic|links|graph|grafana|dashboard>
	-log-base=<base directory>. Defualt value is "../change-log/"
	-db-connection=<connection string> - default value is "root:cepave@tcp(192.168.20.50:3306)"
	-db-type=<type> - default value is "mysql"
	```
* `export-changelog.sh` - Dump table of `sysdb_change_log` to maintain **global schema**
	```
	Usage export-changelog.sh -db-name=<database_name> {{ mysqldump options... }}
	Options:

	-db-name=<database_name>
	```
* `patch_boss.sh`(**DEPRECATED**) - Dedicated to patching "boss" database, you should use **run-owl-patch.sh** instead.
