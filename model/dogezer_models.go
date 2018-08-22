// DOGEZER RZ:
package model

import (
	"encoding/json"
	"io"
)

// ---------------------------------------- LastPostAt ---------------
type LastPostAt struct {
	Id         string `json:"id"`
	LastPostAt int64  `json:"last_post_at"`
}

func (lpa *LastPostAt) ToJson() string {
	b, _ := json.Marshal(lpa)
	return string(b)
}
func LastPostAtFromJson(data io.Reader) *LastPostAt {
	var o *LastPostAt
	json.NewDecoder(data).Decode(&o)
	return o
}

// ---------------------------------------- LastsPosts ---------------
type LastsPosts struct {
	Channels []*LastPostAt `json:"by_channels"`
	Threads  []*LastPostAt `json:"by_threads"`
}

func (lp *LastsPosts) ToJson() string {
	b, err := json.Marshal(lp)
	if err != nil {
		return "[]"
	} else {
		return string(b)
	}
}
func LastsPostsFromJson(data io.Reader) *LastsPosts {
	decoder := json.NewDecoder(data)
	var o LastsPosts
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

// ---------------------------------------- ChannelsUnreads ---------------
type ChannelsUnreads struct {
	Channels []*ChannelUnread `json:"by_channels"`
	Threads  []*ThreadUnread  `json:"by_threads"`
}

func (cu *ChannelsUnreads) ToJson() string {
	b, err := json.Marshal(cu)
	if err != nil {
		return "[]"
	} else {
		return string(b)
	}
}
func ChannelsUnreadsFromJson(data io.Reader) *ChannelsUnreads {
	decoder := json.NewDecoder(data)
	var o ChannelsUnreads
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

// ---------------------------------------- PostUnread ---------------
type PostUnread struct {
	PostId     string `json:"post_id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	UserId     string `json:"user_id"`
	ChannelId  string `json:"channel_id"`
	LastPostAt int64  `json:"last_post_at"`
	LastPostId string `json:"last_post_id"`
}

func (pu *PostUnread) ToJson() string {
	data, err := json.Marshal(pu)
	if err != nil {
		return ""
	}
	return string(data)
}
func PostUnreadFromJson(data io.Reader) *PostUnread {
	decoder := json.NewDecoder(data)
	var o PostUnread
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

// ---------------------------------------- ThreadUnread ---------------
type ThreadUnread struct {
	RootId        string `json:"root_id"`
	LastViewedAt  int64  `json:"last_viewed_at"`
	FirstUnreadAt int64  `json:"first_unread_at"`
	LastPostAt    int64  `json:"last_post_at"`
	MsgCount      int64  `json:"msg_count"`
}

func (tu *ThreadUnread) ToJson() string {
	data, err := json.Marshal(tu)
	if err != nil {
		return ""
	}
	return string(data)
}

func ThreadUnreadFromJson(data io.Reader) *ThreadUnread {
	decoder := json.NewDecoder(data)
	var o ThreadUnread
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}
