package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

type Role struct{}

func (r *Role) ListUserRoles(ctx *gin.Context) {
	db := common.GetGorm()
	var roleList []struct {
		ID       int64  `json:"id"`
		RoleName string `json:"roleName"`
	}
	r1 := db.Model(&common.TRole{}).Find(&roleList)
	if r1.Error != nil {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	Response(ctx, errorcode.Success, roleList, true, "操作成功")
}

type reqListRoles struct {
	Current  int    `form:"current"`
	Size     int    `form:"size"`
	Keywords string `form:"keywords"`
}

func (r *Role) ListRoles(ctx *gin.Context) {
	var form reqListRoles
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "系统异常")
		return
	}
	if form.Current <= 0 || form.Size <= 0 {
		form.Current = 1
		form.Size = 10
	}
	db := common.GetGorm()
	var count int64
	var roleList []struct {
		ID             int64   `json:"id"`
		IsDisable      int     `json:"isDisable"`
		MenuIdList     []int64 `json:"menuIdList"`
		ResourceIdList []int64 `json:"resourceIdList"`
		RoleLabel      string  `json:"roleLabel"`
		RoleName       string  `json:"roleName"`
	}
	if form.Keywords == "" {
		r1 := db.Model(&common.TRole{}).Count(&count)
		r1 = db.Model(&common.TRole{}).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&roleList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
		}
	} else {
		r1 := db.Model(&common.TRole{}).Where(fmt.Sprintf("role_name LIKE %q", "%"+form.Keywords+"%")).Count(&count)
		r1 = db.Model(&common.TRole{}).Where(fmt.Sprintf("role_name LIKE %q", "%"+form.Keywords+"%")).Limit(form.Size).Offset((form.Current - 1) * form.Size).Find(&roleList)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
		}
	}

	for i, val := range roleList {
		var mList []common.TRoleMenu
		var rList []common.TRoleResource
		r2 := db.Where("role_id = ?", val.ID).Find(&mList)
		r2 = db.Where("role_id = ?", val.ID).Find(&rList)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		for _, v := range mList {
			roleList[i].MenuIdList = append(roleList[i].MenuIdList, v.MenuId)
		}
		for _, v := range rList {
			roleList[i].ResourceIdList = append(roleList[i].ResourceIdList, v.ResourceId)
		}
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["recordList"] = roleList
	Response(ctx, errorcode.Success, data, true, "操作成功")
}

type reqSaveOrUpdateRole struct {
	ID             int64   `json:"id"`
	IsDisable      int     `json:"isDisable"`
	MenuIdList     []int64 `json:"menuIdList"`
	ResourceIdList []int64 `json:"resourceIdList"`
	RoleLabel      string  `json:"roleLabel"`
	RoleName       string  `json:"roleName"`
}

func (r *Role) SaveOrUpdateRole(ctx *gin.Context) {
	var form reqSaveOrUpdateRole
	if err := ctx.ShouldBindJSON(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	if form.ID == 0 {
		role := common.TRole{
			IsDisable: form.IsDisable,
			RoleLabel: form.RoleLabel,
			RoleName:  form.RoleName,
		}
		r1 := db.Model(&common.TRole{}).Create(&role)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		var rmList []common.TRoleMenu
		for _, val := range form.MenuIdList {
			tmp := common.TRoleMenu{
				RoleId: role.ID,
				MenuId: val,
			}
			rmList = append(rmList, tmp)
		}
		var rrList []common.TRoleResource
		for _, val := range form.ResourceIdList {
			tmp := common.TRoleResource{
				RoleId:     role.ID,
				ResourceId: val,
			}
			rrList = append(rrList, tmp)
		}
		if rmList != nil {
			r1 = db.Create(&rmList)
			if r1.Error != nil {
				logger.Error(r1.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		if rrList != nil {
			r1 = db.Create(&rrList)
			if r1.Error != nil {
				logger.Error(r1.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}

	} else {
		var role common.TRole
		r1 := db.Where("id = ?", form.ID).First(&role)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		role.IsDisable = form.IsDisable
		role.RoleLabel = form.RoleLabel
		role.RoleName = form.RoleName
		role.UpdateTime = time.Now()
		r1 = db.Save(&role)
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}

		var rmList []common.TRoleMenu
		for _, val := range form.MenuIdList {
			var c common.TRoleMenu
			r2 := db.Where("role_id = ? AND menu_id = ?", role.ID, val).First(&c)
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			//如果存在则没必要创建
			if r2.Error == nil {
				continue
			}
			tmp := common.TRoleMenu{
				RoleId: role.ID,
				MenuId: val,
			}
			rmList = append(rmList, tmp)
		}
		var rrList []common.TRoleResource
		for _, val := range form.ResourceIdList {
			var c common.TRoleResource
			r2 := db.Where("role_id = ? AND resource_id = ?", role.ID, val).First(&c)
			if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
				logger.Error(r2.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
			//如果存在则没必要创建
			if r2.Error == nil {
				continue
			}
			tmp := common.TRoleResource{
				RoleId:     role.ID,
				ResourceId: val,
			}
			rrList = append(rrList, tmp)
		}
		if rmList != nil {
			r1 = db.Create(&rmList)
			if r1.Error != nil {
				logger.Error(r1.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
		if rrList != nil {
			r1 = db.Create(&rrList)
			if r1.Error != nil {
				logger.Error(r1.Error.Error())
				Response(ctx, errorcode.Fail, nil, false, "系统异常")
				return
			}
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
func (r *Role) DeleteRoles(ctx *gin.Context) {
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	str := fmt.Sprintf("%v", string(data))
	var idList []int64
	err := json.Unmarshal([]byte(str), &idList)
	if err != nil {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	db := common.GetGorm()
	for _, val := range idList {
		r1 := db.Where("role_id = ?", val).Delete(&common.TRoleMenu{})
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		r1 = db.Where("role_id = ?", val).Delete(&common.TRoleResource{})
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		r1 = db.Where("role_id = ?", val).Delete(&common.TUserRole{})
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		r1 = db.Where("id = ?", val).Delete(&common.TRole{})
		if r1.Error != nil {
			logger.Error(r1.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
	}
	Response(ctx, errorcode.Success, nil, true, "操作成功")
}
