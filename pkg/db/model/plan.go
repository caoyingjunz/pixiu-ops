/*
Copyright 2024 The Pixiu Authors.

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

package model

import "github.com/caoyingjunz/pixiu/pkg/db/model/pixiu"

func init() {
	register(&Plan{})
	register(&Node{})
}

type Plan struct {
	pixiu.Model

	Name        string `gorm:"index:idx_name,unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

func (plan *Plan) TableName() string {
	return "plans"
}

type KubeRole int

const (
	NodeRole   KubeRole = iota // kubernetes node role
	MasterRole                 // kubernetes master role
)

type Node struct {
	pixiu.Model

	Name   string   `gorm:"index:idx_name,unique" json:"name"`
	PlanId int64    `json:"plan_id"`
	Role   KubeRole `json:"role"` // k8s 节点的角色，master 为 1 和 node 为 0
	Ip     string   `json:"ip"`
	Auth   string   `json:"auth"`
}

func (*Node) TableName() string {
	return "nodes"
}
