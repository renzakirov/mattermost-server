// DOGEZER RZ:
package sqlstore

import (
	"fmt"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

type SqlMentionStore struct {
	SqlStore
	metrics einterfaces.MetricsInterface
}

const (
	MENTIONS_TIME_CACHE_SIZE = 25000
	MENTIONS_TIME_CACHE_SEC  = 900 // 15 minutes
)

func NewSqlMentionStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.MentionStore {
	s := &SqlMentionStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Mention{}, "Mentions").SetKeys(false, "UserId", "PostId")
		table.ColMap("UserId").SetMaxSize(26)
		table.ColMap("ChannelId").SetMaxSize(26)
		table.ColMap("RootId").SetMaxSize(26)
		table.ColMap("PostId").SetMaxSize(26)
		table.ColMap("AuthorId").SetMaxSize(26)
	}

	return s
}

func (s SqlMentionStore) CreateIndexesIfNotExists() {

	s.CreateIndexIfNotExists("idx_mentions_userid_createat", "Mentions", "UserId, CreateAt")

}

func (s SqlMentionStore) Save(post *model.Post, id string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var mention model.Mention
		mention.UserId = id
		mention.ChannelId = post.ChannelId
		mention.RootId = post.RootId
		mention.PostId = post.Id
		mention.CreateAt = post.CreateAt
		mention.AuthorId = post.UserId

		if err := s.GetMaster().Insert(&mention); err != nil {
			// TODO переделать Err

			fmt.Println(err)
			fmt.Println(err)

			result.Err = model.NewAppError("SqlMentionStore.Save", "store.sql_mention.save", nil, "EROROROROR", 10)
		} else {
			result.Data = mention
		}
	})
}

func (s SqlMentionStore) View(mention *model.Mention) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		result.Err = model.NewAppError("SqlMentionStore.View", "store.sql_post_unread.view", nil, "EROROROROR", 10)
		result.Data = mention
		return
	})
}
