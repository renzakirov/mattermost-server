package app

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

func (a *App) ViewPost(unread *model.PostUnread) (*model.PostUnread, *model.AppError) {
	var pchan store.StoreChannel
	pchan = a.Srv.Store.Post().GetSingle(unread.PostId)
	if pchan != nil {
		if presult := <-pchan; presult.Err != nil {
			return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
		} else {
			post := presult.Data.(*model.Post)
			if post == nil {
				return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
			}
			unread.ChannelId = post.ChannelId
			pchan = a.Srv.Store.PostUnread().View(unread)
			if presult := <-pchan; presult.Err != nil {
				return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusInternalServerError)
			} else {
				a.sendViewThreadEvent(unread)
				return presult.Data.(*model.PostUnread), nil
			}

		}
	}
	//a.Srv.Store.User().Get(post.UserId)
	return nil, model.NewAppError("ViewPost", "api.post.view_post.app_error", nil, "", http.StatusBadRequest)
}
