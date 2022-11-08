package service

import (
	"github.com/beancount-gs/script"
	"github.com/gin-gonic/gin"
	"strings"
)

type Tags struct {
	Value string `bql:"distinct tags" json:"value"`
}

func QueryTags(c *gin.Context) {
	ledgerConfig := script.GetLedgerConfigFromContext(c)
	tags := make([]Tags, 0)
	err := script.BQLQueryList(ledgerConfig, nil, &tags)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	result := make([]string, 0)
	m := make(map[string]bool, 0)
	for _, t := range tags {
		if t.Value != "" {
			itemTags := strings.Split(t.Value, ",")
			for _, rt := range itemTags {
				rt = strings.TrimSpace(rt)
				if rt != "" {
					if _, ok := m[rt]; !ok {
						result = append(result, script.GetShowTag(ledgerConfig.Id, rt))
						m[rt] = true
					}
				}
			}
		}
	}

	OK(c, result)
}
