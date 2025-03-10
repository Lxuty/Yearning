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

package group

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/lib"
	"Yearning-go/src/model"
	"encoding/json"
	"fmt"
	"github.com/cookieY/yee"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
)

func SuperGetRuseSource(c yee.Context) (err error) {
	var source []model.CoreDataSource
	var query []model.CoreDataSource
	model.DB().Select("source,source_id").Scopes(commom.AccordingToGroupSourceIsQuery(0, 2)).Find(&source)
	model.DB().Select("source,source_id").Scopes(commom.AccordingToGroupSourceIsQuery(1, 2)).Find(&query)
	return c.JSON(http.StatusOK, commom.SuccessPayload(map[string]interface{}{"source": source, "query": query}))
}

func SuperGroup(c yee.Context) (err error) {
	var page int
	var roles []model.CoreRoleGroup

	f := new(commom.PageChange)
	if err = c.Bind(f); err != nil {
		return err
	}
	start, end := lib.Paging(f.Current, f.PageSize)
	model.DB().Model(model.CoreRoleGroup{}).Scopes(commom.AccordingToOrderName(f.Expr.Text)).Count(&page).Offset(start).Limit(end).Find(&roles)
	return c.JSON(http.StatusOK, commom.SuccessPayload(
		commom.CommonList{
			Page: page,
			Data: roles,
		},
	))
}

func SuperGroupUpdate(c yee.Context) (err error) {
	user := new(lib.Token).JwtParse(c)
	if user.Username == "admin" {
		u := new(policy)
		if err = c.Bind(u); err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
		}
		g, err := json.Marshal(u.PermissionList)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		}

		if u.ID == 0 {
			model.DB().Create(&model.CoreRoleGroup{
				Name:        u.Name,
				Permissions: g,
				GroupId:     uuid.New().String(),
			})
		} else {
			model.DB().Model(model.CoreRoleGroup{}).Scopes(commom.AccordingToIDEqual(u.ID)).Update(&model.CoreRoleGroup{Permissions: g, Name: u.Name})
		}
		return c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(GROUP_CREATE_SUCCESS, u.Name)))
	}
	return c.JSON(http.StatusOK, commom.ERR_REQ_FAKE)
}

func SuperClearUserRule(c yee.Context) (err error) {
	args := c.QueryParam("group_id")
	scape, _ := url.QueryUnescape(args)
	var j []model.CoreGrained
	model.DB().Scopes(commom.AccordingToGroupNameIsLike(scape)).Find(&j)
	for _, i := range j {
		b, err := lib.ArrayRemove(i.Group, scape)
		if err != nil {
			return c.JSON(http.StatusOK, commom.ERR_COMMON_MESSAGE(err))
		}
		go model.DB().Model(model.CoreGrained{}).Scopes(commom.AccordingToUsernameEqual(i.Username)).Update(&model.CoreGrained{Group: b})
	}
	model.DB().Model(model.CoreRoleGroup{}).Where("group_id = ?", scape).Delete(&model.CoreRoleGroup{})
	return c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(fmt.Sprintf(GROUP_DELETE_SUCCESS, scape)))
}

func SuperUserRuleMarge(c yee.Context) (err error) {
	u := new(marge)
	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusOK, commom.ERR_REQ_BIND)
	}
	m3 := lib.MultiUserRuleMarge(strings.Split(u.Group, ","))
	return c.JSON(http.StatusOK, commom.SuccessPayload(m3))
}
