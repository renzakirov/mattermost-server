package model

import (
	"encoding/json"
	"io"
)

// DOGEZER RZ:

type PostUnread struct {
	PostId     string `json:"post_id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	UserId     string `json:"user_id"`
	ChannelId  string `json:"channel_id"`
	LastPostAt int64  `json:"last_post_at"`
	LastPostId string `json:"last_post_id"`
}

type ChannelPostUnread struct {
	RootId        string `json:"root_id"`
	LastViewedAt  int64  `json:"last_viewed_at"`
	FirstUnreadAt int64  `json:"first_unread_at"`
	LastPostAt    int64  `json:"last_post_at"`
	MsgCount      int64  `json:"msg_count"`
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
