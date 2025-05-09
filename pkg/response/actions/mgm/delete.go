package mgm

/*
 * @Author: lwnmengjing
 * @Date: 2023/1/25 17:13:59
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2023/1/25 17:13:59
 */

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/boilingsitar/mss-boot/pkg/response"
)

// Delete action
type Delete struct {
	Base
	Key string
}

// NewDelete new deleteMgm action
func NewDelete(b Base, key string) *Delete {
	return &Delete{
		Base: b,
		Key:  key,
	}
}

func (e *Delete) Handler() gin.HandlersChain {
	h := func(c *gin.Context) {
		ids := make([]string, 0)
		v := c.Param(e.Key)
		if v == "batch" {
			api := response.Make(c).Bind(&ids)
			if api.Error != nil || len(ids) == 0 {
				api.Err(http.StatusUnprocessableEntity)
				return
			}
			e.delete(c, ids...)
			return
		}
		e.delete(c, v)
	}
	chain := gin.HandlersChain{h}
	return chain
}

// String action name
func (*Delete) String() string {
	return "deleteMgm"
}

// delete batch and single deleteMgm
func (e *Delete) delete(c *gin.Context, ids ...string) {
	api := response.Make(c)
	if len(ids) == 0 {
		api.OK(nil)
		return
	}
	_, err := mgm.Coll(e.Model).DeleteMany(c, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			api.OK(nil)
			return
		}
		api.Err(http.StatusInternalServerError)
		return
	}
	api.OK(nil)
}
