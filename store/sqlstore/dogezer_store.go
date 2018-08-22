package sqlstore

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

// DOGEZER RZ:
func (s SqlChannelStore) GetAllLastPostsAt(userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		fmt.Println("---------- in store GetAllLastPostsAt -> userId = ", userId)
		// var unreadsChannels model.ChannelsUnreads
		var data model.LastsPosts
		params := map[string]interface{}{"UserId": userId}
		_, err := s.GetReplica().Select(&data.Channels,
			`select 
				c.id as Id,
				max(c.lastpostat) as LastPostAt
			from channelmembers as m
			join channels as c 
			on c.id = m.channelid
			WHERE 
				userid = :UserId
			group by c.id
			`,
			params)
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetAllLastPostsAt", "store.sql_channel.get_unread.app_error", nil, " GetAllLastPostsAt "+err.Error(), http.StatusInternalServerError)
			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			_, err = s.GetReplica().Select(&data.Threads,
				`
				select 
					p.rootid as Id,
					max(p.createat) as LastPostAt
				from 
					channelmembers as m
					right join posts as p
					on p.channelid = m.channelid
				where  
					m.userid = :UserId and 
					p.rootid != ''
				group by p.rootid
			`, params)
			fmt.Println("-- -> err = ", err)
			result.Data = &data
		}
	})
}

// DOGEZER RZ:
func (s SqlChannelStore) GetAllChannelsUnreads(userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		// fmt.Println("---------- in store getAllChannelsUnreads -> userId = ", userId)
		// var unreadsChannels model.ChannelsUnreads
		var data model.ChannelsUnreads
		params := map[string]interface{}{"UserId": userId}
		_, err := s.GetReplica().Select(&data.Channels,
			`SELECT
				c.teamid TeamId,
				CM.ChannelId ChannelId, 
				(c.TotalMsgCount - cm.MsgCount) MsgCount, 
				cm.MentionCount MentionCount, 
				cm.LastViewedAt LastViewedAt
			FROM Channels as C join ChannelMembers as CM 
			on CM.channelid = c.id and cm.userid= :UserId
			WHERE
				UserId = :UserId
				AND DeleteAt = 0
			`,
			params)
		// GetChannelUnread
		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetChannelUnread", "store.sql_channel.get_unread.app_error", nil, "all user channels "+err.Error(), http.StatusInternalServerError)
			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			_, err = s.GetReplica().Select(&data.Threads,
				`
				select count(*),
					p.rootid as RootId,
					coalesce(min(u.lastpostat), 0) as LastViewedAt,
					min(p.createat) as FirstUnreadAt,
					max(p.createat) as LastPostAt,
					count(*) as MsgCount
				from postunreads as u
				right join posts as p
				on u.postid = p.rootid and
					u.userid = :UserId
				where
					p.rootid != '' and
					p.userid != :UserId and
					(p.createat > u.lastpostat or u.lastpostat is null)
				group by p.rootid
			`, params)
			fmt.Println("-- -> err = ", err)
			result.Data = &data
		}
	})
}

func (s SqlChannelStore) GetChannelUnreads(channelId, userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var unreadChannel model.ChannelUnread
		err := s.GetReplica().SelectOne(&unreadChannel,
			`SELECT
				Channels.TeamId TeamId, Channels.Id ChannelId, (Channels.TotalMsgCount - ChannelMembers.MsgCount) MsgCount, ChannelMembers.MentionCount MentionCount, ChannelMembers.NotifyProps NotifyProps
			FROM
				Channels, ChannelMembers
			WHERE
				Id = ChannelId
                AND Id = :ChannelId
                AND UserId = :UserId
                AND DeleteAt = 0`,
			map[string]interface{}{"ChannelId": channelId, "UserId": userId})

		if err != nil {
			result.Err = model.NewAppError("SqlChannelStore.GetChannelUnread", "store.sql_channel.get_unread.app_error", nil, "channelId="+channelId+" "+err.Error(), http.StatusInternalServerError)
			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			// DOGEZER RZ:
			// Нужно кол-во постов после даты из PostUnreads
			_, err = s.GetReplica().Select(&unreadChannel.ByThreads,
				`
				select count(*),
					p.rootid as RootId,
					coalesce(min(u.lastpostat), 0) as LastViewedAt,
					min(p.createat) as FirstUnreadAt,
					max(p.createat) as LastPostAt,
					count(*) as MsgCount
				from postunreads as u 
				right join posts as p 
				on u.postid = p.rootid and 
					u.userid = :UserId
				where 
					p.rootid != '' and 
					p.channelid = :ChannelId and 
					p.userid != :UserId and 
					(p.createat > u.lastpostat or u.lastpostat is null)
				group by p.rootid
				`,
				map[string]interface{}{"ChannelId": channelId, "UserId": userId})
			fmt.Println("-- -> err = ", err)
			result.Data = &unreadChannel
		}
	})
}

func (s SqlChannelStore) GetThreadUnreads(threadId, userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var threadUnreads model.ThreadUnread
		err := s.GetReplica().SelectOne(&threadUnreads,
			`
				select 
					p.rootid as RootId,
					coalesce(min(u.lastpostat), 0) as LastViewedAt,
					min(p.createat) as FirstUnreadAt,
					max(p.createat) as LastPostAt,
					count(*) as MsgCount
				from 
				postunreads as u
				right join posts as p
				on 
					u.postid = :ThreadId and
					u.userid = :UserId
				where 
					p.rootid = :ThreadId and
					p.userid != :UserId and
					(p.createat > u.lastpostat or u.lastpostat is null)
				group by p.rootid
			`,
			map[string]interface{}{"ThreadId": threadId, "UserId": userId})
		fmt.Println("-- -> err = ", err)
		result.Data = &threadUnreads
	})
}
