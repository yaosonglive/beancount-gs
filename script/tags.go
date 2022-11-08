package script

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

type Tag struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

var ledgerTagsMap map[string]map[string]string

func GetLedgerTags(ledgerId string) map[string]string {
	return ledgerTagsMap[ledgerId]
}

func UpdateLedgerTags(ledgerId string, accountTypesMap map[string]string) {
	ledgerTagsMap[ledgerId] = accountTypesMap
}

func ClearLedgerTags(ledgerId string) {
	delete(ledgerTagsMap, ledgerId)
}

func GetShowTag(ledgerId string, acc string) string {
	tags := ledgerTagsMap[ledgerId]
	tag := acc
	for key, name := range tags {
		if strings.Compare(strings.TrimSpace(acc), key) == 0 {
			tag = name
		}
	}
	return tag
}

func GetTagByShow(ledgerId string, acc string) string {
	tags := ledgerTagsMap[ledgerId]
	tag := strings.TrimSpace(acc)
	hasHan := false
	for _, r := range []rune(tag) {
		if unicode.Is(unicode.Han, r) {
			hasHan = true
		}
	}
	tagKey := tag
	if hasHan {
		for key, name := range tags {
			if strings.Compare(tag, name) == 0 {
				tagKey = key
			}
		}
	}
	return tagKey
}

func LoadLedgerTagsMap(config Config) error {
	path := GetLedgerTagsMapFilePath(config.DataPath)
	if !FileIfExist(path) {
		err := WriteFile(path, "{}")
		if err != nil {
			return err
		}
	}
	fileContent, err := ReadFile(path)
	if err != nil {
		return err
	}
	tags := make(map[string]string)
	err = json.Unmarshal(fileContent, &tags)
	if err != nil {
		LogSystemError("Failed unmarshal config file (" + path + ")")
		return err
	}
	if ledgerTagsMap == nil {
		ledgerTagsMap = make(map[string]map[string]string)
	}
	ledgerTagsMap[config.Id] = tags
	LogSystemInfo(fmt.Sprintf("Success load [%s] account type cache", config.Mail))
	return nil
}
