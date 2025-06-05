package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal/log"
)

var i18nData map[string]map[string]string
var DefaultLocale = "en"

func InitI18n() {
	path := "./locales"
	i18nData = make(map[string]map[string]string)
	entries, err := os.ReadDir("locales")
	if err != nil {
		fmt.Println("Error reading locales directory:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filename := path + "/" + entry.Name()
			locale := getLocaleFromFilename(entry.Name())

			data, err := os.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
				continue
			}
			flattenedData, err := Flatten(data, ".")
			if err != nil {
				fmt.Println("Error flattening JSON:", err)
				continue
			}
			i18nData[locale] = flattenedData
			fmt.Printf("Loaded locale: %s from %s\n", locale, filename)
		}
	}
	log.Logger.Println("初始化i18n...success")
}

func getLocaleFromFilename(filename string) string {
	name := filename[:len(filename)-len(".json")]
	return name
}

func Tl(locale string, key string, args ...interface{}) string {
	langData, ok := i18nData[locale] // 使用默认 locale
	if !ok {
		return key
	}
	text, ok := langData[key]
	if !ok {
		return key // 如果 key 不存在，返回 key
	}
	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}

func T(ctx *gin.Context, key string, args ...interface{}) string {
	return Tl(ctx.GetString("locale"), key, args...)
}

func Flatten(jsonBytes []byte, sep string) (map[string]string, error) {
	var jsonObj map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return nil, err
	}

	flattened := make(map[string]string)
	flattenJSON(jsonObj, "", sep, flattened)
	return flattened, nil
}

func flattenJSON(jsonObj map[string]interface{}, parentKey string, sep string, flattened map[string]string) {
	for k, v := range jsonObj {
		newKey := k
		if parentKey != "" {
			newKey = parentKey + sep + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			flattenJSON(val, newKey, sep, flattened)
		case []interface{}:
			for i, item := range val {
				switch itemVal := item.(type) {
				case map[string]interface{}:
					flattenJSON(itemVal, newKey+sep+strconv.Itoa(i), sep, flattened)
				default:
					flattened[newKey+sep+strconv.Itoa(i)] = itemVal.(string)
				}
			}
		default:
			flattened[newKey] = val.(string)
		}
	}
}
