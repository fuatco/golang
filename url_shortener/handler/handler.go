package handler

import (
	js "encoding/json"
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	yamlV3 "gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yaml []byte) (dst []map[string]string, err error) {
	err = yamlV3.Unmarshal(yaml, &dst)
	return dst, err
}

func buildMap(parsedYaml []map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for _, entry := range parsedYaml {
		key := entry["path"]
		mergedMap[key] = entry["url"]
	}
	return mergedMap
}

func parseJson(json []byte) (dst map[string]string, err error) {
	err = js.Unmarshal(json, &dst)
	return dst, err
}

func JsonHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJson(json)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedJson, fallback), nil
}

func RedisHanlder(conn redis.Conn, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		path, ok := redis.String(conn.Do("HGET", "urls", r.URL.Path))
		if ok == nil {
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
