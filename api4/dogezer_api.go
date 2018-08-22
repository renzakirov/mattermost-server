// DOGEZER RZ:
package api4

import (
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
)

func getAllChannelsUnreads(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	channelsUnreads, err := c.App.GetAllChannelsUnreads(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(channelsUnreads.ToJson()))
}

func getAllLastPostsAt(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	channelsUnreads, err := c.App.GetAllLastPostsAt(c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(channelsUnreads.ToJson()))
}

func getThreadUnreads(c *Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("--- api -> getThreadUnreads begin")
	c.RequireUserId()
	if c.Err != nil {
		return
	}

	if !c.App.SessionHasPermissionToUser(c.Session, c.Params.UserId) {
		c.SetPermissionError(model.PERMISSION_EDIT_OTHER_USERS)
		return
	}

	threadUnreads, err := c.App.GetThreadUnreads(c.Params.PostId, c.Params.UserId)
	if err != nil {
		c.Err = err
		return
	}
	fmt.Println("--- api -> getThreadUnreads end -> threadUnreads = ", threadUnreads)
	w.Write([]byte(threadUnreads.ToJson()))
}
