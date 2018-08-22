// DOGEZER RZ:

package app

import "github.com/mattermost/mattermost-server/model"

// DOGEZER RZ:
func (a *App) GetAllLastPostsAt(userId string) (*model.LastsPosts, *model.AppError) {
	result := <-a.Srv.Store.Channel().GetAllLastPostsAt(userId)
	if result.Err != nil {
		return nil, result.Err
	}
	lastsPosts := result.Data.(*model.LastsPosts)

	// if channelUnread.NotifyProps[model.MARK_UNREAD_NOTIFY_PROP] == model.CHANNEL_MARK_UNREAD_MENTION {
	// 	channelsUnreads.MsgCount = 0
	// }

	return lastsPosts, nil
}

func (a *App) GetThreadUnreads(threadId, userId string) (*model.ThreadUnread, *model.AppError) {
	result := <-a.Srv.Store.Channel().GetThreadUnreads(threadId, userId)
	if result.Err != nil {
		return nil, result.Err
	}
	threadUnreads := result.Data.(*model.ThreadUnread)

	// if channelUnread.NotifyProps[model.MARK_UNREAD_NOTIFY_PROP] == model.CHANNEL_MARK_UNREAD_MENTION {
	// 	channelsUnreads.MsgCount = 0
	// }

	return threadUnreads, nil
}
