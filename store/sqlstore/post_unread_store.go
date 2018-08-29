// DOGEZER
package sqlstore

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

type SqlPostUnreadStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

const (
	LAST_POST_UNREAD_TIME_CACHE_SIZE = 25000
	LAST_POST_UNREAD_TIME_CACHE_SEC  = 900 // 15 minutes

	//LAST_POSTS_CACHE_SIZE = 1000
	//LAST_POSTS_CACHE_SEC  = 900 // 15 minutes
)

func NewSqlPostUnreadStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.PostUnreadStore {
	s := &SqlPostUnreadStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.PostUnread{}, "PostUnreads").SetKeys(false, "PostId", "ChannelId", "UserId")
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("ChannelId").SetMaxSize(26)
		table.ColMap("PostId").SetMaxSize(26)
	}

	return s
}

func (s SqlPostUnreadStore) CreateIndexesIfNotExists() {

	s.CreateIndexIfNotExists("idx_post_unreads_user_channel", "PostUnreads", "UserId, ChannelId")
	/*
		s.CreateIndexIfNotExists("idx_posts_create_at", "Posts", "CreateAt")
		s.CreateIndexIfNotExists("idx_posts_delete_at", "Posts", "DeleteAt")
		s.CreateIndexIfNotExists("idx_posts_channel_id", "Posts", "ChannelId")
		s.CreateIndexIfNotExists("idx_posts_root_id", "Posts", "RootId")
		s.CreateIndexIfNotExists("idx_posts_user_id", "Posts", "UserId")
		s.CreateIndexIfNotExists("idx_posts_is_pinned", "Posts", "IsPinned")

		s.CreateFullTextIndexIfNotExists("idx_posts_message_txt", "Posts", "Message")
		s.CreateFullTextIndexIfNotExists("idx_posts_hashtags_txt", "Posts", "Hashtags")
	*/
}

func (s SqlPostUnreadStore) View(unread *model.PostUnread) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		unread.UpdateAt = model.GetMillis()

		fmt.Println("-----------store.Post_unread_store -> View -> unread = ", unread)

		r, err := s.GetMaster().Exec(`
			UPDATE PostUnreads 
			SET 
				LastPostAt = :LastPostAt, 
				UpdateAt = :UpdateAt, 
				LastPostId = :LastPostId 
			WHERE 
				ChannelId = :ChannelId AND 
				UserId = :UserId AND 
				PostId = :PostId`,
			unread)
		aff, _ := r.RowsAffected()
		fmt.Println("-----------store.Post_unread_store -> View -> r.RowsAffected = ", aff)

		if err == nil && aff == 0 {
			unread.CreateAt = unread.UpdateAt
			err = s.GetMaster().Insert(unread)
			if err != nil {
				result.Err = model.NewAppError("SqlPostUnreadStore.View", "store.sql_post_unread.view", nil, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				return
			}
		}
		result.Data = unread
		if err != nil {
			result.Err = model.NewAppError("SqlPostUnreadStore.View", "store.sql_post_unread.view", nil, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
	})
}
