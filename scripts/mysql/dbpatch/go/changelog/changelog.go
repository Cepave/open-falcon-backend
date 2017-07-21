// As the loading process for change log
package changelog

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// As the patch configuration
type PatchConfig struct {
	// The id of patch
	Id string `yaml:"id"`
	// The file name of patch
	Filename string `yaml:"filename"`
	// The comment of patch
	Comment string `yaml:"comment"`
}

// Loads configurations of patches from a string
func LoadChangeLog(changeLogOfYaml []byte) (configOfPatches []PatchConfig, err error) {
	/**
	 * Unmarshal the content YAML to PatchConfig
	 */
	configOfPatches = make([]PatchConfig, 0, 8)
	if err = yaml.Unmarshal(changeLogOfYaml, &configOfPatches); err != nil {
		configOfPatches = nil
		return
	}
	// :~)

	return
}

// Loads configurations of patches from a reader
func LoadChangeLogFromReader(readerOfchangeLog io.Reader) (configOfPatches []PatchConfig, err error) {
	var yamlContent []byte

	if yamlContent, err = ioutil.ReadAll(bufio.NewReader(readerOfchangeLog)); err != nil {
		return
	}

	return LoadChangeLog(yamlContent)
}

// Loads configurations of patches from a file name(path)
// The file would be auto-closed by this method
func LoadChangeLogFromFile(changeLogFile string) (configOfPatches []PatchConfig, err error) {
	var fileOfChangeLog *os.File

	/**
	 * Reads change log of YAML
	 */
	if fileOfChangeLog, err = os.Open(changeLogFile); err != nil {
		return nil, err
	}

	defer fileOfChangeLog.Close()
	// :~)

	return LoadChangeLogFromReader(fileOfChangeLog)
}

const ESCAPED_DELIMITER = "!DBPATCH!ESCAPED_DELIMITER!"

// Opens the file of patch and reads the script in it(splitted by delimiter)
func (patchConfig *PatchConfig) loadScripts(folderBase string, delimiter string) (scripts []string, err error) {
	/**
	 * Open file
	 */
	var patchFile = patchConfig.Filename
	if folderBase != "" {
		patchFile = fmt.Sprintf("%s/%s", folderBase, patchFile)
	}

	var targetFile *os.File
	if targetFile, err = os.Open(patchFile); err != nil {
		return
	}

	defer targetFile.Close()
	// :~)

	/**
	 * Reads content of patch file and spilt the content by delimiter
	 * ** Replace the escaped delimiter with special string **
	 */
	var contentOfScript []byte
	if contentOfScript, err = ioutil.ReadAll(targetFile); err != nil {
		return
	}

	stringOfScript := string(contentOfScript)
	stringOfScript = strings.Replace(stringOfScript, delimiter+delimiter, ESCAPED_DELIMITER, -1)

	var rowScripts = strings.Split(
		stringOfScript, delimiter,
	)
	// :~)

	/**
	 * Skip the empty content of script
	 * e.x. CREATE TABLE xxx;;
	 *
	 * ** Replace back the escaped delimiters**
	 */
	scripts = make([]string, 0, len(rowScripts)-1)
	for _, rowScript := range rowScripts {
		rowScript = strings.TrimSpace(rowScript)
		rowScript = strings.Replace(rowScript, ESCAPED_DELIMITER, delimiter, -1)

		if rowScript == "" {
			continue
		}

		scripts = append(scripts, rowScript)
	}
	// :~)

	return
}
