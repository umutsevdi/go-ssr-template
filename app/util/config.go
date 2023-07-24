package util

/******************************************************************************

 * File: util/config.go
 *
 * Author: Umut Sevdi
 * Created: 07/04/23
 * Description:Configuration parsing utilities.

*****************************************************************************/
import (
	"encoding/json"
	"log"
	"os"
)

var (
	D_URI             string = ""
	D_PORT            uint64 = 8080
	D_API_PORT        uint64 = 8081
	D_CACHE           bool   = false
	D_CACHE_VALID_TTL bool   = true
	C                 Config = Config{}
)

type Cache struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Ttl     uint64 `json:"ttl,omitempty"`
}

type Config struct {
	// URI of the website, will be used while generating sitemaps.xml
	URI *string `json:"uri,omitempty"`
	// Port to launch the application
	Port *uint64 `json:"port,omitempty"`
	// Periodic page indexing
	PIndexing *Cache `json:"page,omitempty"`
	// Will be deprecated
	StaticCache *Cache `json:"static,omitempty"`
	// Caching of filled templates for different users.
	MemoryCache *Cache `json:"memory,omitempty"` //
	// Path to pages static and components directories
	ContentPath *string `json:"content,omitempty"`
}

// Parses the configuration file at $WEBWATCH_CONFIG or config.json file
// whichever is available.
func init() {
	path, found := os.LookupEnv("WEBWATCH_CONFIG")
	if !found || path == "" {
		path = "config.json"
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("ERROR: config.json is not found at \"", path,
			"\". Either create a config file called \"config.json\", or define a valid",
			"$WEBWATCH_CONFIG.")
	}

	json.Unmarshal(file, &C)
	fillEmptyFields()
	log.Println("Server has been started with following configurations:",
		"\n- target:      ", *C.URI, ":", *C.Port,
		"\n- periodicIndexing:   {enabled: ", *C.PIndexing.Enabled, ", ttl: ", C.PIndexing.Ttl, "}",
		"\n- staticCache: {enabled: ", *C.StaticCache.Enabled, ", ttl: ", C.StaticCache.Ttl, "}",
		"\n- memoryCache: {enabled: ", *C.MemoryCache.Enabled, ", ttl: ", C.MemoryCache.Ttl, "}",
	)
}

// Sanitizes the invalid inputs from the configuration file.
func fillEmptyFields() {
	if C.URI == nil {
		log.Println("WARN: \"uri\" is not defined at configuration. Continuing with localhost")
		C.URI = &D_URI
	}
	if C.Port == nil {
		log.Println("WARN : \"port\" is not defined. Continuing with 8080.")
		C.Port = &D_PORT
	}
	if C.StaticCache == nil {
		C.StaticCache = &Cache{}
		C.StaticCache.Enabled = &D_CACHE
	} else if C.StaticCache.Ttl > 0 {
		C.StaticCache.Enabled = &D_CACHE_VALID_TTL
	}
	if C.PIndexing == nil {
		C.PIndexing = &Cache{}
		C.PIndexing.Enabled = &D_CACHE
	} else if C.PIndexing.Ttl > 0 {
		C.PIndexing.Enabled = &D_CACHE_VALID_TTL
	}
	if C.MemoryCache == nil {
		C.MemoryCache = &Cache{}
		C.MemoryCache.Enabled = &D_CACHE
	} else if C.MemoryCache.Ttl > 0 {
		C.MemoryCache.Enabled = &D_CACHE_VALID_TTL
	}
	if C.ContentPath == nil {
		log.Println("WARN : \"port\" is not defined. Defaulting to 8080.")
		*C.ContentPath = "content"
	}
}
