package util

/******************************************************************************

 * File: util/index.go
 *
 * Author: Umut Sevdi
 * Created: 07/04/23
 * Description: File indexing and caching utilities. Indexes paths of content. If enabled,
 it can periodically update as well.

*****************************************************************************/

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strings"
	"time"
)

type FData struct {
	Path         string
	Content      []byte
	LastModified time.Time
}

type ItemContainer struct {
	// Contains the data of template elements
	components map[string]FData
	// Contains addresses of pages
	pages map[string]FData
	// Contains addresses of cont.static items such as images
	static map[string]FData
	lock   bool
}

var cont ItemContainer = ItemContainer{}

func MapEach(table string, f func(string, FData)) {
	for cont.lock {

	}
	var m *map[string]FData
	switch table {
	case "components":
		m = &cont.components
	case "pages":
		m = &cont.pages
	case "static":
		m = &cont.static
	}
	if m == nil {
		return
	}
	for k, v := range *m {
		f(k, v)
	}
}

func MapGet(table, key string) *FData {
	for cont.lock {

	}
	var d FData
	var ok bool = false
	switch table {
	case "components":
		d, ok = cont.components[key]
	case "pages":
		d, ok = cont.pages[key]
	case "static":
		d, ok = cont.static[key]
	}
	if ok {
		return &d
	}
	return nil
}

// Indexes all available contents.
//
// - If the Page Indexing is enabled at the configuration, runs indexing in the
// background periodically.
func StartIndexing() {
	if *C.PIndexing.Enabled {

		var ticker *time.Ticker
		if C.PIndexing.Ttl > 0 {
			ticker = time.NewTicker(time.Duration(C.PIndexing.Ttl) * time.Minute)
			log.Println("Periodic Caching enabled for every", C.PIndexing.Ttl, "minutes")
		} else {
			ticker = time.NewTicker(time.Duration(1) * time.Second)
			log.Println("Periodic Caching enabled for every second.")
			log.Println("WARN: Do not use this on live environments.")
		}
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					indexContent()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}
	indexContent()
}

// Parses received URL and extracts it's extension
//
//	@param URL path to file
//	@return string corresponding extension type in the format of .type
func Ext(url string) string {
	p := strings.Split(url, "/")
	fname := strings.Split(p[len(p)-1], ".")
	return "." + fname[len(fname)-1]
}

// Indexing function that updates cont.components, cont.pages and cont.static files. Indexes only contain
// the respective file path.
func indexContent() {
	log.Println("INDEX BEGIN   : ", *C.ContentPath)
	cont.lock = true
	if len(cont.components) == 0 {
		cont.components = make(map[string]FData, 100)
	}
	if len(cont.pages) == 0 {
		cont.pages = make(map[string]FData, 100)
	}
	if len(cont.static) == 0 {
		cont.static = make(map[string]FData, 100)
	}

	indexHtml("/components", "/", &cont.components)
	indexHtml("/pages", "/", &cont.pages)
	processTemplate()
	indexStatic("/", &cont.static)
	// TODO cleanup function that deletes removed files
	cont.lock = false
	log.Println("INDEX COMPLETE: ", *C.ContentPath)
}

// Recursively traverses through the content/pages directory
// and indexes pages.
//
// - File extensions are removed while indexing
//
//	@param path - path to traverse
//	@param uris - map to insert items
//
// "content/pages"
// "content/components"
func indexHtml(basePath string, path string, uris *map[string]FData) {
	data, err := os.ReadDir(*C.ContentPath + basePath + path)
	if err != nil {
		return
	}
	for _, v := range data {
		if v.Type().IsDir() {
			indexHtml(basePath, path+v.Name()+"/", uris)
		} else {
			var key string = path + v.Name()
			var path string = *C.ContentPath + basePath + path + v.Name()
			key = strings.Split(key, ".")[0]
			if key == "/index" {
				key = "/"
			}
			fetchFData(key, path, uris)
		}
	}
}

// Recursively traverses through the content/static directory
// and indexes static files at /static URL path.
//
//	@param path - path to traverse
//	@param uris - map to insert items
func indexStatic(path string, uris *map[string]FData) {
	data, err := os.ReadDir(*C.ContentPath + "/static" + path)
	if err != nil {
		return
	}
	for _, v := range data {
		if v.Type().IsDir() {
			indexStatic(path+v.Name()+"/", uris)
		} else {
			var key string = "/static" + path + v.Name()
			var path string = *C.ContentPath + key
			fetchFData(key, path, uris)
		}
	}
}

// Caches given path for the key
// - If the file is not registered, inserts it with last modified date
// - If the file is already cached, checks whether it has been updated or not
// updates only it's changed
func fetchFData(key, path string, files *map[string]FData) {
	if fs, ok := (*files)[key]; ok {
		metadata, err := os.Stat(fs.Path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist")
			return
		}
		// Update only if the file is changed
		if metadata.ModTime().After(fs.LastModified) {
			fs.LastModified = metadata.ModTime()
			fs.Content, err = os.ReadFile(fs.Path)
			if err != nil {
				log.Println("WARN: File", path, "does not exist for metadata")
				return
			}
		}
	} else {
		// if doesn't exist create
		d, err := os.ReadFile(path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist")
			return
		}
		metadata, err := os.Stat(path)
		if err != nil {
			log.Println("WARN: File", path, "does not exist for metadata")
			return
		}
		// Update only if the file is changed
		(*files)[key] = FData{
			Path:         path,
			Content:      d,
			LastModified: metadata.ModTime(),
		}
	}
}

// Replaces template variables on pages with actual components
func processTemplate() {
	comp := make([]string, len(cont.components))
	i := 0

	for _, v := range cont.components {
		if v.Content != nil && len(v.Content) > 0 {
			comp[i] = string(v.Content)
			i++
		}
	}
	for k, v := range cont.pages {
		t, err := template.New(k).Parse(string(v.Content))
		if err != nil {
			log.Println("1", err.Error())
			err = nil
		}
		for _, v := range comp {
			_, err = t.Parse(v)
			if err != nil {
				log.Println("2", err.Error())
				err = nil
			}
		}
		c, _ := cont.pages[k]

		var w bytes.Buffer
		err = t.Execute(&w, nil)
		if err != nil {
			log.Println("3", err.Error())
			err = nil
		}
		cont.pages[k] = FData{
			c.Path,
			[]byte(w.String()),
			c.LastModified,
		}
	}
}
