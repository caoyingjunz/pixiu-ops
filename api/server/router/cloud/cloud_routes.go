/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"context"
	"fmt"

	"github.com/caoyingjunz/gopixiu/api/server/httputils"
	"github.com/caoyingjunz/gopixiu/api/types"
	"github.com/caoyingjunz/gopixiu/pkg/pixiu"
	"github.com/caoyingjunz/gopixiu/pkg/util"
	
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func readConfig(c *gin.Context) ([]byte, error) {
	config, err := c.FormFile("kubeconfig")
	if err != nil {
		return nil, err
	}
	file, err := config.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func (s *cloudRouter) createCloud(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err   error
		cloud types.Cloud
	)
	cloud.Name = c.Param("name")
	if len(cloud.Name) == 0 {
		httputils.SetFailed(c, r, fmt.Errorf("invaild empty cloud name"))
		return
	}
	cloud.KubeConfig, err = readConfig(c)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.Cloud().Create(context.TODO(), &cloud); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) updateCloud(c *gin.Context) {
	r := httputils.NewResponse()
	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) deleteCloud(c *gin.Context) {
	r := httputils.NewResponse()
	cid, err := util.ParseInt64(c.Param("cid"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.Cloud().Delete(context.TODO(), cid); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) getCloud(c *gin.Context) {
	r := httputils.NewResponse()
	cid, err := util.ParseInt64(c.Param("cid"))
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.Cloud().Get(context.TODO(), cid)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) listClouds(c *gin.Context) {
	r := httputils.NewResponse()
	var err error
	if r.Result, err = pixiu.CoreV1.Cloud().List(context.TODO()); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) createNamespace(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.CloudOptions
		namespace   corev1.Namespace
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = c.ShouldBindJSON(&namespace); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.Cloud().CreateNamespace(context.TODO(), listOptions.CloudName, namespace); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) updateNamespace(c *gin.Context) {
}
func (s *cloudRouter) deleteNamespace(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.CloudOptions
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	namespace := struct {
		Name string `json:"name,omitempty"`
	}{}
	if err = c.ShouldBindJSON(&namespace); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = pixiu.CoreV1.Cloud().DeleteNamespace(context.TODO(), listOptions.CloudName, namespace.Name); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}
func (s *cloudRouter) getNamespace(c *gin.Context) {}

func (s *cloudRouter) listNamespaces(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err          error
		cloudOptions types.CloudOptions
	)
	if err = c.ShouldBindUri(&cloudOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.Cloud().ListNamespaces(context.TODO(), cloudOptions)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) createDeployment(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err        error
		getOptions types.GetOrCreateOptions
		deployment v1.Deployment
	)
	if err = c.ShouldBindUri(&getOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	if err = c.ShouldBindJSON(&deployment); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	deployment.Name = getOptions.ObjectName
	deployment.Namespace = getOptions.Namespace
	if err = pixiu.CoreV1.Cloud().CreateDeployment(context.TODO(), getOptions.CloudName, &deployment); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) deleteDeployment(c *gin.Context) {
	r := httputils.NewResponse()
	var deleteOptions types.GetOrDeleteOptions
	if err := c.ShouldBindUri(&deleteOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	err := pixiu.CoreV1.Cloud().DeleteDeployment(context.TODO(), deleteOptions)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

// listDeployments API: clouds/<cloud_name>/namespaces/<ns>/deployments
func (s *cloudRouter) listDeployments(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.ListOptions
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.Cloud().ListDeployments(context.TODO(), listOptions)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}

func (s *cloudRouter) listJobs(c *gin.Context) {
	r := httputils.NewResponse()
	var (
		err         error
		listOptions types.ListOptions
	)
	if err = c.ShouldBindUri(&listOptions); err != nil {
		httputils.SetFailed(c, r, err)
		return
	}
	r.Result, err = pixiu.CoreV1.Cloud().ListJobs(context.TODO(), listOptions)
	if err != nil {
		httputils.SetFailed(c, r, err)
		return
	}

	httputils.SetSuccess(c, r)
}
