package assets

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/openshift-psap/special-resource-operator/pkg/color"
	errs "github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	log logr.Logger
)

func init() {
	log = zap.New(zap.UseDevMode(true)).WithName(color.PrettyPrint("manifests", color.Brown))
}

// Metadata manifests filename and content
type Metadata struct {
	Name    string
	Content []byte
}

// GetFrom reads all manifests 0000- from path and returns them
func GetFrom(assets string) []Metadata {

	manifests := []Metadata{}
	files, err := filePathWalkDir(assets, ".yaml")
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		buffer, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		manifests = append(manifests, Metadata{path.Base(file), buffer})
	}
	return manifests
}

func filePathWalkDir(root string, ext string) ([]string, error) {

	var files []string

	if _, err := os.Stat(root); os.IsNotExist(err) {
		if errs.Wrap(err, "Directory does note exists, giving up: "+root) != nil {
			log.Info("Exiting On", "error", err)
			os.Exit(1)
		}
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			log.Info("WalkDir", "path IsDir", path)
			// Ignore root directory but skipdir any subdirectories
			if path == root {
				return nil
			}
			return filepath.SkipDir
		}
		if filepath.Ext(path) != ext {
			log.Info("WalkDir", "path does not match *.yaml", path)
			return nil
		}

		if valid := filePathPatternValid(path); !valid {
			return nil
		}

		log.Info("WalkDir", "path valid", path)
		files = append(files, path)
		return nil
	})
	return files, err
}

func filePathPatternValid(path string) bool {

	patterns := []string{
		"[0-9][0-9][0-9][0-9]-*.yaml",
		"[0-9][0-9][0-9][0-9]_*.yaml",
	}

	for _, pattern := range patterns {
		if result, _ := filepath.Match(pattern, filepath.Base(path)); !result {
			continue
		}
		return true
	}
	return false
}
