package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (json JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(json).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	Nickname string   `json:"nickname"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
