package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"

	"github.com/joaovictorsl/dcache"
	"github.com/joaovictorsl/fooche"
	"github.com/joaovictorsl/fooche/evict"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type dcacheConfig struct {
	Port           uint16        `json:"port"`
	Cache          string        `json:"cache"`
	EvictionPolicy string        `json:"eviction-policy"`
	CleanInterval  cleanInterval `json:"clean-interval"`
	SizeCapConfig  map[int]int   `json:"size-cap-config"`
	MaxValueLength uint
}

type cleanInterval struct {
	Interval int `json:"interval"`
}

type genPolicyFn func(capacity int) evict.EvictionPolicy[string]

func main() {
	config := loadConfig()
	cache := createCache(config)

	dcacheServer := dcache.NewServer(
		config.Port,
		cache,
		config.MaxValueLength,
	)

	if err := dcacheServer.Start(); err != nil {
		panic(err)
	}
}

func loadConfig() dcacheConfig {
	cfile := flag.String("config-file", "example.json", "DCache config file")
	flag.Parse()

	data, err := readConfigJSON(*cfile)
	if err != nil {
		panic(err)
	}

	var config dcacheConfig
	json.Unmarshal(data, &config)

	// Get the maximum value length
	// Sort keys in descending order
	keys := maps.Keys(config.SizeCapConfig)
	config.MaxValueLength = uint(slices.Max(keys))

	return config
}

func readConfigJSON(cfile string) ([]byte, error) {
	f, err := os.Open(cfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getCreatePolicyFn(config dcacheConfig) genPolicyFn {
	return func(capacity int) evict.EvictionPolicy[string] {
		switch config.EvictionPolicy {
		case "lru":
			return evict.NewLRU[string](capacity)
		default:
			panic("Invalid eviction policy")
		}
	}
}

func createCache(config dcacheConfig) fooche.ICache {
	createPolicyFn := getCreatePolicyFn(config)

	switch config.Cache {
	case "clean-interval":
		return fooche.NewCleanIntervalBounded(
			time.Duration(config.CleanInterval.Interval)*time.Second,
			config.SizeCapConfig,
			createPolicyFn,
		)
	case "simple":
		return fooche.NewSimpleBounded(
			config.SizeCapConfig,
			createPolicyFn,
		)
	default:
		panic("Invalid cache type")
	}
}
