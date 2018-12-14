package helpers

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/alexmspina/worldmap/server/models"
)

// CreateRegexp creates a map of string keys and their regular expression counterpart values
func CreateRegexp(r map[string]*regexp.Regexp, p []string) {
	for _, s := range p {
		regex, err := regexp.Compile(s)
		models.PanicErrors(err)
		r[s] = regex
	}
}

// GetFilesFromDirectory use the Walk function to create a list of files found in the given directory
func GetFilesFromDirectory(f *[]string, d string) {
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		*f = append(*f, path)
		return nil
	})
	models.PanicErrors(err)
}

// AppendBytes adds a byte slice to another by byte
func AppendBytes(mainslice *[]byte, addingslice []byte) {
	for _, i := range addingslice {
		*mainslice = append(*mainslice, i)
	}
}
