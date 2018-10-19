package sqlstore

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

// DOGEZER RZ:
func (s *SqlPostStore) GetNLastPosts(userId string, limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		fmt.Println("---------- in store GetNLastPosts -> userId = ", userId)
		pl := model.NewPostList()
		var data []*model.Post

		userBehalf := "{\"on_behalf\":\"" + userId + "\"}"
		params := map[string]interface{}{"UserId": userId, "N": limit, "UserBehalf": userBehalf}

		_, err := s.GetReplica().Select(&data,
			`
				select
					p.*
				from channelmembers as c
				join posts as p
				on c.channelid = p.channelid and
					c.userid = :UserId
				where
					p.userid != :UserId and
					p.type not like '%system%' and 
					p.type not like 'system%' and 
					(p.type != 'custom_dogezer_behalf' or props != :UserBehalf)
				ORDER BY p.createat DESC
				limit :N
			`,
			params)
		if err != nil {
			fmt.Println("---------- ERROR in store GetNLastPosts -> err = ", err)
		} else {
			for _, post := range data {
				pl.AddPost(post)
				pl.AddOrder(post.Id)
			}
		}
		result.Data = pl
	})
}

// DOGEZER RZ:
func (s *SqlChannelStore) GetAllLastPostsAt(userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		fmt.Println("---------- in store GetAllLastPostsAt -> userId = ", userId)
		// var unreadsChannels model.ChannelsUnreads
		var data model.LastsPosts
		params := map[string]interface{}{"UserId": userId}
		_, err := s.GetReplica().Select(&data.Channels,
			`
				select 
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
		userBehalf := "{\"on_behalf\":\"" + userId + "\"}"
		params := map[string]interface{}{"UserId": userId, "UserBehalf": userBehalf}
		_, err := s.GetReplica().Select(&data.Channels,
			`
				select
					p.channelid as ChannelId,
					coalesce(min(c.lastviewedat), 0) as LastViewedAt,
					coalesce(max(c.mentioncount), 0) as MentionCount,
					count(*) as MsgCount
				from channelmembers as c
				right join posts as p
				on c.channelid = p.channelid and
					c.userid = :UserId
				where
					(p.createat > c.lastviewedat or c.lastviewedat is null) and
					p.rootid = '' and
					p.userid != :UserId and
					p.type not like '%system%' and 
					p.type not like 'system%' and 
					(p.type != 'custom_dogezer_behalf' or props != :UserBehalf)
				group by p.channelid
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
					(p.createat > u.lastpostat or u.lastpostat is null) and
					p.rootid != '' and
					p.userid != :UserId and
					p.type not like 'system' and 
					(p.type != 'custom_dogezer_behalf' or props != :UserBehalf)
				group by p.rootid
			`, params)

			if err == nil {
				_, err = s.GetReplica().Select(&data.Mentions,
					`
					select 
						m.rootid as RootId,
						count (*) as Count
					from postunreads as u
					right join mentions as m
					on
						u.userid = :UserId and 
						u.postid = m.rootid
					where
						m.rootid != ''
						and m.userid = :UserId
						and (m.createat > u.lastpostat or u.lastpostat is null)
					group by m.rootid
				`, params)

				if err != nil {
					fmt.Println("-- -> err = ", err)
				}
			}

			fmt.Println("-- -> err = ", err)
			result.Data = &data
		}
	})
}

func (s SqlChannelStore) GetChannelUnreads(channelId, userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var channelUnread model.ChannelUnread
		err := s.GetReplica().SelectOne(&channelUnread,
			`SELECT
				Channels.TeamId TeamId, 
				Channels.Id ChannelId, 
				(Channels.TotalMsgCount - ChannelMembers.MsgCount) MsgCount, 
				ChannelMembers.MentionCount MentionCount, 
				ChannelMembers.NotifyProps NotifyProps
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
			_, err = s.GetReplica().Select(&channelUnread.ByThreads,
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

			// DOGEZER RZ:
			// Нужно кол-во mentions после даты из PostUnreads
			_, err = s.GetReplica().Select(&channelUnread.Mentions,
				`
					select 
						m.rootid as RootId,
						count (*) as Count
					from postunreads as u
					right join mentions as m
					on
						u.userid = :UserId and 
						u.postid = m.rootid
					where
						m.rootid != ''
						and m.channelid = :ChannelId
						and m.userid = :UserId
						and (m.createat > u.lastpostat or u.lastpostat is null)
					group by m.rootid
				`,
				map[string]interface{}{"ChannelId": channelId, "UserId": userId})
			fmt.Println("-- -> err = ", err)

			result.Data = &channelUnread
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

		mCount, err := s.GetReplica().SelectInt(
			`
					select 
						count (*) as Count
					from postunreads as u
					right join mentions as m
					on
						u.postid = :ThreadId and
						u.userid = :UserId
					where
						m.rootid = :ThreadId and
						m.userid = :UserId and
						(m.createat > u.lastpostat or u.lastpostat is null)
				`,
			map[string]interface{}{"ThreadId": threadId, "UserId": userId})
		fmt.Println("-- -> err = ", err)
		threadUnreads.MentionCount = mCount
		result.Data = &threadUnreads
	})
}

func (s SqlChannelStore) DChannelView(channelInfo *model.ChannelInfo, userId string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		fmt.Println("DChannelView -> ", channelInfo)

		// достать last_viewed_at из channelmembers
		// если больше - то сделать
		// меньше или равен достать последние значения

		// var lastV []struct {
		// 	LastViewedAt int64
		// }

		var LastViewedAt int64
		props := make(map[string]interface{})

		props["ChannelId"] = channelInfo.ChannelId
		props["UserId"] = userId

		query := `
			SELECT lastviewedat FROM ChannelMembers 
			WHERE userid = :UserId and
			channelid = :ChannelId
		`
		if err := s.GetMaster().SelectOne(&LastViewedAt, query, props); err != nil {
			result.Err = model.NewAppError("2 SqlChannelStore.DChannelView", "store.sql_channel.update_last_viewed_at.app_error", nil, "channel_id="+channelInfo.ChannelId+", user_id="+userId+", "+err.Error(), http.StatusInternalServerError)
			return
		}

		if LastViewedAt > channelInfo.LastViewedAt {
			channelInfo.MsgCount = 5
			channelInfo.MentionCount = 7
			channelInfo.LastViewedAt = LastViewedAt
		} else {
			props["LastAt"] = channelInfo.LastViewedAt
			props["UserBehalf"] = "{\"on_behalf\":\"" + userId + "\"}"
			var msgCount int64
			err := s.GetReplica().SelectOne(&msgCount,
				`
					select
						count(*)
					from posts
					where channelid = :ChannelId and
					rootid = '' and
					createat > :LastAt and
					userid != :UserId and
					type not like 'system%' and 
					type not like '%system%' and 
					(type != 'custom_dogezer_behalf' or props != :UserBehalf)
				`,
				props)

			if err != nil {
				return
			}
			var mentionCount int64
			err = s.GetReplica().SelectOne(&mentionCount,
				`
					select
						count(*)
					from mentions
					where channelid = :ChannelId and
					rootid = '' and
					createat > :LastAt and
					userid = :UserId 
				`,
				props)
			if err != nil {
				return
			}

			props["MsgCount"] = msgCount
			props["MentionCount"] = mentionCount

			r, err := s.GetMaster().Exec(`
			UPDATE channelmemebers 
			SET 
				LastViewedAt = :LastViewedAt, 
				msgcount = :MsgCount, 
				mentioncount = :MentionCount 
				lastupdateat = :LastViewedAt
			WHERE 
				ChannelId = :ChannelId AND 
				UserId = :UserId`,
				props)
		}

		result.Data = &channelInfo
	})
}
