package assets

import (
	"log"
	"net/http"
	"sync"

	"github.com/rakyll/statik/fs"
)

var instance http.FileSystem
var once sync.Once

// GetFS return assets file system instance
func GetFS() http.FileSystem {
	once.Do(func() {
		var err error
		instance, err = fs.New()
		if err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
