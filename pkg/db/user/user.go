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

package user

import (
	"context"
	"time"

	"gorm.io/gorm"

	dberrors "github.com/caoyingjunz/gopixiu/pkg/db/errors"
	"github.com/caoyingjunz/gopixiu/pkg/db/model"
	"github.com/caoyingjunz/gopixiu/pkg/log"
)

type UserInterface interface {
	Create(ctx context.Context, obj *model.User) (*model.User, error)
	Update(ctx context.Context, uid int64, resourceVersion int64, updates map[string]interface{}) error
	Delete(ctx context.Context, uid int64) error
	Get(ctx context.Context, uid int64) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)

	GetByName(ctx context.Context, name string) (*model.User, error)
	GetRoleIDByUser(ctx context.Context, uid int64) (map[string][]int64, error)
	SetUserRoles(ctx context.Context, uid int64, rid []int64) error
	GetButtonsByUserID(ctx context.Context, uid, menuId int64) (*[]model.Menu, error)
	GetLeftMenusByUserID(ctx context.Context, uid int64) (*[]model.TreeMenus, error)
	DeleteRolesByUserID(ctx context.Context, uid int64) error
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) UserInterface {
	return &user{db}
}

func (u *user) Create(ctx context.Context, obj *model.User) (*model.User, error) {
	// 系统维护字段
	now := time.Now()
	obj.GmtCreate = now
	obj.GmtModified = now

	if err := u.db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (u *user) Update(ctx context.Context, uid int64, resourceVersion int64, updates map[string]interface{}) error {
	// 系统维护字段
	updates["gmt_modified"] = time.Now()
	updates["resource_version"] = resourceVersion + 1

	f := u.db.Model(&model.User{}).
		Where("id = ? and resource_version = ?", uid, resourceVersion).
		Updates(updates)
	if f.Error != nil {
		return f.Error
	}

	if f.RowsAffected == 0 {
		return dberrors.ErrRecordNotUpdate
	}

	return nil
}

func (u *user) Delete(ctx context.Context, uid int64) error {
	return u.db.
		Where("id = ?", uid).
		Delete(&model.User{}).
		Error
}

func (u *user) Get(ctx context.Context, uid int64) (*model.User, error) {
	var obj model.User
	if err := u.db.Where("id = ?", uid).First(&obj).Error; err != nil {
		return nil, err
	}

	return &obj, nil
}

func (u *user) List(ctx context.Context) ([]model.User, error) {
	var us []model.User
	if err := u.db.Find(&us).Error; err != nil {
		return nil, err
	}

	return us, nil
}

func (u *user) GetByName(ctx context.Context, name string) (*model.User, error) {
	var obj model.User
	if err := u.db.Where("name = ?", name).First(&obj).Error; err != nil {
		return nil, err
	}
	return &obj, nil
}

//SetUserRoles 分配用户角色
func (u *user) SetUserRoles(ctx context.Context, uid int64, rids []int64) (err error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Errorf(err.(error).Error())
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		log.Logger.Errorf(err.Error())
		tx.Rollback()
		return err
	}
	if err := tx.Where(&model.UserRole{UserID: uid}).Delete(&model.UserRole{}).Error; err != nil {
		log.Logger.Errorf(err.Error())
		tx.Rollback()
		return err
	}
	if len(rids) > 0 {
		for _, rid := range rids {
			rm := new(model.UserRole)
			rm.RoleID = rid
			rm.UserID = uid
			if err := tx.Create(rm).Error; err != nil {
				log.Logger.Errorf(err.Error())
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

//GetRoleIDByUser 查询用户角色
func (u *user) GetRoleIDByUser(ctx context.Context, uid int64) (map[string][]int64, error) {
	roleList := []int64{}
	err := u.db.Model(&model.UserRole{}).Where("user_id = ?", uid).Pluck("role_id", &roleList).Error
	if err != nil {
		log.Logger.Errorf(err.Error())
		return nil, err
	}
	roleInfo := map[string][]int64{
		"role_ids": roleList,
	}
	return roleInfo, nil
}

//GetButtonsByUserID 获取菜单按钮
func (u *user) GetButtonsByUserID(ctx context.Context, uid, menuId int64) (*[]model.Menu, error) {
	var menus []model.Menu
	err := u.db.Table("menu").Select(" menu.id, menu.parent_id,menu.name, menu.url, menu.icon,menu.sequence,"+
		"menu.code,menu.method, menu.menu_type").
		Joins("left join role_menu on menu.id = role_menu.menu_id ").
		Joins("left join user_role on user_role.user_id = role_menu.role_id  and user_role.user_id = ?", uid).
		Where("menu.menu_type = 2 and menu.parent_id = ?", menuId).
		Order("parent_id ASC").
		Order("sequence ASC").
		Scan(&menus).Error
	if err != nil {
		log.Logger.Errorf(err.Error())
		return nil, err
	}

	return &menus, nil
}

// GetLeftMenusByUserID 根据用户ID获取左侧菜单
func (u *user) GetLeftMenusByUserID(ctx context.Context, uid int64) (*[]model.TreeMenus, error) {
	var menus []model.Menu
	err := u.db.Table("menu").Select(" menu.id, menu.parent_id,menu.name, menu.url, menu.icon,menu.sequence,"+
		"menu.code,menu.method, menu.menu_type").
		Joins("left join role_menu on menu.id = role_menu.menu_id ").
		Joins("left join user_role on user_role.user_id = role_menu.role_id  and user_role.user_id = ?", uid).
		Where("menu.menu_type = 1").
		Order("parent_id ASC").
		Order("sequence ASC").
		Scan(&menus).Error
	if err != nil {
		log.Logger.Errorf(err.Error())
		return nil, err
	}
	treeMenusList := getTreeMenus(menus, 0)
	return &treeMenusList, nil
}

func (u *user) DeleteRolesByUserID(ctx context.Context, uid int64) error {
	err := u.db.Where("user_id = ?", uid).Delete(&model.UserRole{}).Error
	if err != nil {
		log.Logger.Errorf(err.Error())
		return err
	}
	return nil
}
