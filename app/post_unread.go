// DOGEZER
package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

func (a *App) ViewPost(unread *model.PostUnread) (*model.ThreadUnread, *model.AppError) {
	var pchan store.StoreChannel
	pchan = a.Srv.Store.Post().GetSingle(unread.PostId)
	if pchan != nil {
		presult := <-pchan
		if presult.Err != nil {
			return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
		}
		post := presult.Data.(*model.Post)
		if post == nil {
			return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
		}
		unread.ChannelId = post.ChannelId
		pchan = a.Srv.Store.PostUnread().View(unread)
		presult = <-pchan
		if presult.Err != nil {
			return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusInternalServerError)
		}

		pchan = a.Srv.Store.PostUnread().GetUnreadsByUserAndRootId(unread)
		presult = <-pchan
		if presult.Err != nil {
			return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusInternalServerError)
		}
		threadUnread := presult.Data.(*model.ThreadUnread)
		a.sendViewThreadEvent(threadUnread, unread)
		return threadUnread, nil
	}
	return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
}
