package lib

import (
	"encoding/json"
	"fmt"
	"time"
)

type ShelfItem struct {
	Id          uint64     `json:"id"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Comment     string     `json:"comment"`
	GMTCreated  CustomTime `json:"gmt_created"`
	GMTModified CustomTime `json:"gmt_modified"`
	GMTDeleted  CustomTime `json:"gmt_deleted"`
	Deleted     bool       `json:"deleted"`
}
type ShelfEditItem struct {
	Id         uint64 `json:"id"`
	NewTitle   string `json:"new_title"`
	NewLink    string `json:"new_link"`
	NewComment string `json:"new_comment"`
}

type ShelfDeleteItem struct {
	Id uint64 `json:"id"`
}

type StrShelfResponse[T any] struct {
	Code uint16 `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type PostRequestResponse struct {
	Code   uint16 `json:"code"`
	Result any    `json:"result"`
}

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	return json.Marshal(t.UnixMilli())
}
func (ct *CustomTime) UnmarshalJSON(data []byte) error {

	var ms int64
	if err := json.Unmarshal(data, &ms); err == nil {
		*ct = CustomTime(time.Unix(0, ms*int64(time.Millisecond)))
		return nil
	} else {
		return fmt.Errorf("cannot unmarshal %s into CustomTime: %w", string(data), err)
	}
}
