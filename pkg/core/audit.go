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

package core

import (
	"context"
	pixiumeta "github.com/caoyingjunz/gopixiu/api/meta"
	"github.com/caoyingjunz/gopixiu/cmd/app/config"
	"github.com/caoyingjunz/gopixiu/pkg/db"
	"github.com/caoyingjunz/gopixiu/pkg/db/model"
	"github.com/caoyingjunz/gopixiu/pkg/log"
)

type AuditGetter interface {
	Audit() AuditInterface
}

type AuditInterface interface {
	Create(ctx context.Context, obj *model.Audit) error
	Delete(c context.Context, ids []int64) error
	List(c context.Context, selector *pixiumeta.ListSelector) (res *model.PageAudit, err error)
}

type audit struct {
	ComponentConfig config.Config
	app             *pixiu
	factory         db.ShareDaoFactory
}

func newAudit(c *pixiu) AuditInterface {
	return &audit{
		ComponentConfig: c.cfg,
		app:             c,
		factory:         c.factory,
	}
}

func (a audit) Create(c context.Context, obj *model.Audit) error {
	if _, err := a.factory.Audit().Create(c, obj); err != nil {
		log.Logger.Errorf("failed to save audit %s: %v", obj.UserID, err)
		return err
	}
	return nil
}

func (a audit) Delete(c context.Context, ids []int64) error {
	if err := a.factory.Audit().Delete(c, ids); err != nil {
		log.Logger.Errorf("batch delete %s audit error: %v", ids, err)
		return err
	}
	return nil
}

func (a audit) List(c context.Context, selector *pixiumeta.ListSelector) (res *model.PageAudit, err error) {
	audits, err := a.factory.Audit().List(c, selector.Page, selector.Limit)
	if err != nil {
		log.Logger.Errorf("list audit error: %v", err)
		return nil, err
	}
	return audits, nil

}
