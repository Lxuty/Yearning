// Copyright 2019 HenryYee.
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"Yearning-go/src/engine"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"encoding/json"
	"fmt"
	"github.com/gookit/gcli/v2/interact"
	"time"
)

func DataInit(o *engine.AuditRole, other *model.Other, ldap *model.Ldap, message *model.Message, a *model.PermissionList) {
	c, _ := json.Marshal(o)
	oh, _ := json.Marshal(other)
	l, _ := json.Marshal(ldap)
	m, _ := json.Marshal(message)
	ak, _ := json.Marshal(a)
	group, _ := json.Marshal([]string{"admin"})
	model.DB().Debug().Create(&model.CoreAccount{
		Username:   "admin",
		RealName:   "超级管理员",
		Password:   lib.DjangoEncrypt("Yearning_admin", string(lib.GetRandom())),
		Department: "DBA",
		Email:      "",
	})
	model.DB().Debug().Create(&model.CoreGlobalConfiguration{
		Authorization: "global",
		Other:         oh,
		AuditRole:     c,
		Message:       m,
		Ldap:          l,
	})
	model.DB().Debug().Create(&model.CoreGrained{
		Username: "admin",
		Group:    group,
	})
	model.DB().Debug().Create(&model.CoreRoleGroup{
		Name:        "admin",
		Permissions: ak,
	})
}

func Migrate() {
	if !model.DB().HasTable("core_accounts") {
		if !interact.Confirm("是否已将数据库字符集设置为UTF8/UTF8MB4?") {
			return
		}
		model.DB().CreateTable(&model.CoreAccount{})
		model.DB().CreateTable(&model.CoreDataSource{})
		model.DB().CreateTable(&model.CoreGlobalConfiguration{})
		model.DB().CreateTable(&model.CoreGrained{})
		model.DB().CreateTable(&model.CoreSqlOrder{})
		model.DB().CreateTable(&model.CoreSqlRecord{})
		model.DB().CreateTable(&model.CoreRollback{})
		model.DB().CreateTable(&model.CoreQueryRecord{})
		model.DB().CreateTable(&model.CoreQueryOrder{})
		model.DB().CreateTable(&model.CoreAutoTask{})
		model.DB().CreateTable(&model.CoreRoleGroup{})
		model.DB().CreateTable(&model.CoreWorkflowTpl{})
		model.DB().AutoMigrate(&model.CoreWorkflowDetail{})
		o := engine.AuditRole{
			DMLInsertColumns:               false,
			DMLMaxInsertRows:               10,
			DMLWhere:                       false,
			DMLOrder:                       false,
			DMLSelect:                      false,
			DDLCheckTableComment:           false,
			DDLCheckColumnNullable:         false,
			DDLCheckColumnDefault:          false,
			DDLEnableAcrossDBRename:        false,
			DDLEnableAutoincrementInit:     false,
			DDLEnableAutoIncrement:         false,
			DDLEnableAutoincrementUnsigned: false,
			DDLEnableDropTable:             false,
			DDLEnableDropDatabase:          false,
			DDLEnableNullIndexName:         false,
			DDLIndexNameSpec:               false,
			DDLMaxKeyParts:                 5,
			DDLMaxKey:                      5,
			DDLMaxCharLength:               10,
			DDLAllowColumnType:             false,
			DDLPrimaryKeyMust:              false,
			MaxTableNameLen:                10,
			MaxAffectRows:                  1000,
			SupportCharset:                 "",
			SupportCollation:               "",
			CheckIdentifier:                false,
			MustHaveColumns:                "",
			DDLMultiToCommit:               false,
			AllowCreatePartition:           false,
			AllowCreateView:                false,
			AllowSpecialType:               false,
		}

		other := model.Other{
			Limit:       1000,
			IDC:         []string{"Aliyun", "AWS"},
			Query:       false,
			Register:    false,
			Export:      false,
			ExQueryTime: 60,
		}

		ldap := model.Ldap{
			Url:      "",
			User:     "",
			Password: "",
			Type:     "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			Sc:       "",
		}

		message := model.Message{
			WebHook:  "",
			Host:     "",
			Port:     25,
			User:     "",
			Password: "",
			ToUser:   "",
			Mail:     false,
			Ding:     false,
			Ssl:      false,
		}

		a := model.PermissionList{
			DDLSource:   []string{},
			DMLSource:   []string{},
			QuerySource: []string{},
		}
		time.Sleep(2)
		DataInit(&o, &other, &ldap, &message, &a)
		fmt.Println("初始化成功!\n用户名: admin\n密码:Yearning_admin\n请通过./Yearning run 运行,默认地址:http://<host>:8000")
	} else {
		fmt.Println("已初始化过,请不要再次执行")
	}
}

func UpdateData() {
	fmt.Println("检查更新.......")
	model.DB().AutoMigrate(&model.CoreAccount{})
	model.DB().AutoMigrate(&model.CoreDataSource{})
	model.DB().AutoMigrate(&model.CoreGlobalConfiguration{})
	model.DB().AutoMigrate(&model.CoreGrained{})
	model.DB().AutoMigrate(&model.CoreSqlOrder{})
	model.DB().AutoMigrate(&model.CoreSqlRecord{})
	model.DB().AutoMigrate(&model.CoreRollback{})
	model.DB().AutoMigrate(&model.CoreQueryRecord{})
	model.DB().AutoMigrate(&model.CoreQueryOrder{})
	model.DB().AutoMigrate(&model.CoreAutoTask{})
	model.DB().AutoMigrate(&model.CoreRoleGroup{})
	model.DB().AutoMigrate(&model.CoreWorkflowTpl{})
	model.DB().AutoMigrate(&model.CoreWorkflowDetail{})
	model.DB().AutoMigrate(&model.CoreOrderComment{})
	model.DB().LogMode(false).Exec("alter table core_auto_tasks change COLUMN base data_base varchar(50) not null")
	model.DB().LogMode(false).Model(&model.CoreSqlOrder{}).DropColumn("uuid")
	model.DB().LogMode(false).Model(&model.CoreWorkflowDetail{}).DropColumn("rejected")
	model.DB().LogMode(false).Model(&model.CoreAutoTask{}).DropColumn("base")
	fmt.Println("数据已更新!")
}

func DelCol() {
	model.DB().LogMode(false).Model(&model.CoreQueryOrder{}).DropColumn("source")
}

func MargeRuleGroup() {
	fmt.Println("破坏性变更修复…………")
	model.DB().LogMode(false).Model(&model.CoreSqlOrder{}).DropColumn("rejected")
	model.DB().LogMode(false).Model(&model.CoreGrained{}).DropColumn("permissions")
	model.DB().LogMode(false).Model(&model.CoreGrained{}).DropColumn("rule")
	ldap := model.Ldap{
		Url:      "",
		User:     "",
		Password: "",
		Type:     "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
		Sc:       "",
	}
	b, _ := json.Marshal(ldap)
	model.DB().LogMode(false).Model(model.CoreGlobalConfiguration{}).Update(&model.CoreGlobalConfiguration{Ldap: b})
	fmt.Println("修复成功!")
}
