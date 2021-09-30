//go:build go1.16
// +build go1.16

package bling

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"strconv"
)

type Context struct {
	UniqueID  string
	RequestId int64
}

var (
	headerUniqueID  = "bling-uid"
	headerRequestID = "bling-rid"
)

//Bling - Entrypoint for bling
func Bling() gin.HandlerFunc {
	fmt.Println("[BLING] bling loaded!")

	return func(c *gin.Context) {
		blingCtx, existing := detectContext(c)
		if existing {
			fmt.Printf("[BLING] found existing bling context (uid: %s, rid: %d)\n", blingCtx.UniqueID, blingCtx.RequestId)
			blingCtx.RequestId += 1
		} else {
			fmt.Printf("[BLING] created new bling context (uid: %s)\n", blingCtx.UniqueID)
		}

		updateHeaders(c, blingCtx)

		c.Next()
	}
}

func newContext() *Context {
	uid, _ := uuid.NewV4()
	return &Context{UniqueID: uid.String(), RequestId: 1}
}

func detectContext(c *gin.Context) (Context, bool) {
	uid := c.GetHeader(headerUniqueID)
	if uid == "" {
		return *newContext(), false
	}

	requestId := c.GetHeader(headerRequestID)
	if requestId == "" {
		requestId = "0"
	}

	rid, err := strconv.ParseInt(requestId, 10, 64)
	if err != nil {
		rid = 0
	}

	b := Context{
		UniqueID:  uid,
		RequestId: rid,
	}

	return b, true
}

func updateHeaders(c *gin.Context, b Context) {
	c.Writer.Header().Set(headerUniqueID, b.UniqueID)
	c.Writer.Header().Set(headerRequestID, strconv.Itoa(int(b.RequestId)))
}
