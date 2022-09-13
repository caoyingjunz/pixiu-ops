package cloud

import (
	"context"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"

	"github.com/caoyingjunz/gopixiu/api/server/httputils"
	"github.com/caoyingjunz/gopixiu/api/types"
	"github.com/caoyingjunz/gopixiu/pkg/pixiu"
)

func (s *cloudRouter) listNodes(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.NodeListOptions
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.Cloud().Nodes(listOptions.CloudName).List(context.TODO())
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) createNode(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.NodeListOptions
		node        v1.Node
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.Cloud().Nodes(listOptions.CloudName).Create(context.TODO(), &node); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	httputils.SetSuccess(c, r)
}
