/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (phe "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plan

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"

	"github.com/caoyingjunz/pixiu/api/server/errors"
	"github.com/caoyingjunz/pixiu/pkg/db/model"
	"github.com/caoyingjunz/pixiu/pkg/types"
)

func (p *plan) preCreateConfig(ctx context.Context, planId int64, req *types.CreatePlanConfigRequest) error {
	_, err := p.factory.Plan().GetConfigByPlan(ctx, planId)
	if err == nil {
		return fmt.Errorf("plan(%d) 配置已存在", planId)
	}

	return nil
}

func (p *plan) CreateConfig(ctx context.Context, pid int64, req *types.CreatePlanConfigRequest) error {
	// 创建前检查
	if err := p.preCreateConfig(ctx, pid, req); err != nil {
		return err
	}

	planConfig, err := p.makePlanConfig(ctx, req)
	if err != nil {
		return err
	}
	// 创建配置
	if _, err = p.factory.Plan().CreatConfig(ctx, planConfig); err != nil {
		klog.Errorf("failed to create plan(%s) config(%d): %v", req.Name, pid, err)
		return err
	}

	return nil
}

// UpdateConfig
// TODO
func (p *plan) UpdateConfig(ctx context.Context, pid int64, cfgId int64, req *types.UpdatePlanConfigRequest) error {
	return nil
}

func (p *plan) DeleteConfig(ctx context.Context, pid int64, cfgId int64) error {
	if _, err := p.factory.Plan().DeleteConfig(ctx, cfgId); err != nil {
		klog.Errorf("failed to delete plan(%d) config(%d): %v", pid, cfgId, err)
		return errors.ErrServerInternal
	}

	return nil
}

func (p *plan) GetConfig(ctx context.Context, pid int64, cfgId int64) (*types.PlanConfig, error) {
	object, err := p.factory.Plan().GetConfigByPlan(ctx, pid)
	if err != nil {
		klog.Errorf("failed to get plan(%d) config(%d): %v", pid, cfgId, err)
		return nil, errors.ErrServerInternal
	}

	return p.modelConfig2Type(object)
}

func (p *plan) makePlanConfig(ctx context.Context, req *types.CreatePlanConfigRequest) (*model.Config, error) {
	kubeConfig, err := req.Kubernetes.Marshal()
	if err != nil {
		return nil, err
	}
	networkConfig, err := req.Network.Marshal()
	if err != nil {
		return nil, err
	}
	runtimeConfig, err := req.Runtime.Marshal()
	if err != nil {
		return nil, err
	}

	return &model.Config{
		PlanId:     req.PlanId,
		Kubernetes: kubeConfig,
		Network:    networkConfig,
		Runtime:    runtimeConfig,
	}, nil
}

func (p *plan) modelConfig2Type(o *model.Config) (*types.PlanConfig, error) {
	ks := &types.KubernetesSpec{}
	if err := ks.Unmarshal(o.Kubernetes); err != nil {
		return nil, err
	}
	ns := &types.NetworkSpec{}
	if err := ns.Unmarshal(o.Network); err != nil {
		return nil, err
	}
	rs := &types.RuntimeSpec{}
	if err := rs.Unmarshal(o.Runtime); err != nil {
		return nil, err
	}

	return &types.PlanConfig{
		PixiuMeta: types.PixiuMeta{
			Id:              o.Id,
			ResourceVersion: o.ResourceVersion,
		},
		TimeMeta: types.TimeMeta{
			GmtCreate:   o.GmtCreate,
			GmtModified: o.GmtModified,
		},
		PlanId:      o.PlanId,
		Name:        o.Name,
		Region:      o.Region,
		Description: o.Description,
		Kubernetes:  *ks,
		Network:     *ns,
		Runtime:     *rs,
	}, nil
}
