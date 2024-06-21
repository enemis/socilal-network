package handler

import (
	"github.com/gin-gonic/gin"
	"social-network-otus/internal/post"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/session"
	"social-network-otus/internal/utils"
	"strconv"
)

type FeedHandler interface {
	GetFeed(c *gin.Context)
}

type FeedHandlerInstance struct {
	response *response.ResponseFactory
	session  *session.SessionStorage
	feed     *post.FeedService
}

func NewFeedHandler(responseFactory *response.ResponseFactory, session *session.SessionStorage, feedService *post.FeedService) *FeedHandlerInstance {
	return &FeedHandlerInstance{session: session, response: responseFactory, feed: feedService}
}

func (h *FeedHandlerInstance) GetFeed(c *gin.Context) {
	user := h.session.GetAuthenticatedUser()
	limit, err := strconv.Atoi(c.Param("limit"))
	errorsF := response.F{}
	if err != nil {
		errorsF["limit"] = err.Error()
	}

	offset, err := strconv.Atoi(c.Param("offset"))
	if err != nil {
		errorsF["offset"] = err.Error()
	}

	if len(errorsF) > 0 {
		h.response.BadRequest(c, errorsF)
		return
	}

	feeds, apperror := h.feed.Feed(user, uint(limit), uint(offset))
	if apperror != nil {
		h.response.FromAppError(c, apperror, utils.Ptr("id"))
		return
	}

	h.response.Ok(c, response.F{"feeds": feeds})
}
