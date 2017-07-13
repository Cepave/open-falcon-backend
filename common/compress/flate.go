package compress

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// compressor of falte algorithm
type FlateCompressor struct {
	// The level of compression
	Level int
}

func NewDefaultFlateCompressor() *FlateCompressor {
	return &FlateCompressor{
		Level: flate.DefaultCompression,
	}
}

func (c *FlateCompressor) CompressString(v string) (result []byte, err error) {
	var b = new(bytes.Buffer)
	var flateWriter *flate.Writer

	flateWriter, err = flate.NewWriter(b, c.Level)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize flate compressor. Error: %v", err)
	}

	_, err = io.Copy(flateWriter, strings.NewReader(v))
	if err != nil {
		return nil, fmt.Errorf("Cannot write data to compressor(flate). Error: %v", err)
	}

	err = flateWriter.Flush()
	if err != nil {
		return nil, fmt.Errorf("Cannot flush data with compressor(flate). Error: %v", err)
	}

	err = flateWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("Close flate writer has error. Error: %v", err)
	}

	result = b.Bytes()

	return
}
func (c *FlateCompressor) DecompressToString(compressedContent []byte) (string, error) {
	flateReader := flate.NewReader(bytes.NewBuffer(compressedContent))

	defer flateReader.Close()

	result, err := ioutil.ReadAll(flateReader)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c *FlateCompressor) MustCompressString(v string) []byte {
	result, err := c.CompressString(v)
	if err != nil {
		panic(err)
	}

	return result
}
func (c *FlateCompressor) MustDecompressToString(compressedContent []byte) string {
	result, err := c.DecompressToString(compressedContent)
	if err != nil {
		panic(err)
	}

	return result
}
