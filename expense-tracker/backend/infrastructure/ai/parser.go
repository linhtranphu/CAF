package ai

import (
	"regexp"
	"strconv"
	"strings"
)

type MessageParser struct{}

func NewMessageParser() *MessageParser {
	return &MessageParser{}
}

func (p *MessageParser) Parse(message string) (string, int64, error) {
	message = strings.ToLower(strings.TrimSpace(message))
	
	amountRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(triệu|nghìn|đồng|k|tr)?`)
	matches := amountRegex.FindStringSubmatch(message)
	
	var amount int64 = 0
	if len(matches) >= 2 {
		if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
			amount = int64(val)
			if len(matches) >= 3 {
				switch matches[2] {
				case "triệu", "tr":
					amount *= 1000000
				case "nghìn", "k":
					amount *= 1000
				}
			}
		}
	}

	items := message
	if matches != nil {
		items = strings.TrimSpace(strings.Replace(message, matches[0], "", 1))
	}
	if items == "" {
		items = "Chi phí khác"
	}

	return items, amount, nil
}