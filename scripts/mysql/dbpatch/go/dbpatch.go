package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Cepave/open-falcon-backend/scripts/mysql/dbpatch/go/changelog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var config *changelog.ChangeLogConfig = loadConfigFromFlags()

	log.Printf("Patching -> %v", config)

	if err := changelog.ExecutePatches(config); err != nil {

		stderr("Patching has error: %v", err)
		os.Exit(1)
	}
}

func loadConfigFromFlags() *changelog.ChangeLogConfig {
	var resultConfig changelog.ChangeLogConfig

	/**
	 * Set-up the fields of configuration
	 */
	flag.StringVar(&resultConfig.DriverName, "driverName", "", "The name of driver")
	flag.StringVar(&resultConfig.Dsn, "dataSourceName", "", "The name of data source")
	flag.StringVar(&resultConfig.ChangeLog, "changeLog", "change-log.json", "The file of change log")
	flag.StringVar(&resultConfig.PatchFileBase, "patchFileBase", "patch-files", "The directory of patch files")
	flag.StringVar(&resultConfig.Delimiter, "delimiter", ";", "The directory of patch files")

	if help := flag.Bool("help", false, "Show help"); *help {
		printCmdUsage()
		os.Exit(0)
	}

	flag.Parse()
	// :~)

	/**
	 * Show help information if some parameters are invalid
	 */
	if !checkRunPatchConfig(&resultConfig) {
		printCmdUsage()
		os.Exit(1)
	}
	// :~)

	return &resultConfig
}

func checkRunPatchConfig(config *changelog.ChangeLogConfig) bool {
	if config.DriverName == "" {
		stderr("Error: need -driverName=<driver_name>\n\n")
		return false
	}
	if config.Dsn == "" {
		stderr("Error: need -dataSourceName=<data_source_name>\n\n")
		return false
	}
	if config.ChangeLog == "" {
		stderr("Error: need -changeLog=<change-log.json>\n\n")
		return false
	}
	if config.PatchFileBase == "" {
		stderr("Error: need -patchFileBase=<patchFileBase>\n\n")
		return false
	}

	return true
}

func printCmdUsage() {
	fmt.Printf("db-patch -driverName=<driver_name> -dataSourceName=<data_srouce_name> [-changeLog=change-log.json] [-patchFileBase=patch-files] [\"-dilimeter=;\"]\n\n")
	flag.PrintDefaults()
}

/**
* Default configuration of running patches
 */
func defaultConfig() *changelog.ChangeLogConfig {
	return &changelog.ChangeLogConfig{
		ChangeLog:     "change-log.json",
		PatchFileBase: "patch-files",
	}
}

func stderr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
