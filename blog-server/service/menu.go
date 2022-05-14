package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Menu struct{}

type reqListMenus struct {
	Keywords string `json:"keywords"`
}
type menusListMenus struct {
	Children   []menusListMenus `json:"children"`
	Component  string           `json:"component"`
	CreateTime time.Time        `json:"createTime"`
	Icon       string           `json:"icon"`
	ID         int64            `json:"id"`
	IsDisable  int              `json:"isDisable"`
	IsHidden   int              `json:"isHidden"`
	OrderNum   int              `json:"orderNum"`
	Path       string           `json:"path"`
	Name       string           `json:"name"`
}

func (m *Menu) ListMenus(ctx *gin.Context) {
	var form reqListMenus
	if err := ctx.ShouldBind(&form); err != nil {
		Response(ctx, errorcode.ValidError, nil, false, "参数校验失败")
		return
	}
	db := common.GetGorm()
	var mList []menusListMenus

	r1 := db.Model(&common.TMenu{}).Where(fmt.Sprintf("isNull(parent_id) AND name LIKE %q", "%"+form.Keywords+"%")).Find(&mList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, m := range mList {
		var child []menusListMenus
		r2 := db.Model(&common.TMenu{}).Where("parent_id = ?", m.ID).Find(&child)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		mList[i].Children = child
	}
	Response(ctx, errorcode.Success, mList, true, "操作成功")

}
func (m *Menu) SaveOrUpdateMenu(ctx *gin.Context) {
	Response(ctx, errorcode.Success, nil, true, "暂时不开放菜单修改和创建功能")
}
func (m *Menu) DeleteMenu(ctx *gin.Context) {
	Response(ctx, errorcode.Success, nil, true, "暂时不开放菜单删除功能")
}

type menusOptions struct {
	ID       int64          `json:"id"`
	Name     string         `json:"label"`
	Children []menusOptions `json:"children"`
}

func (m *Menu) ListMenuOptions(ctx *gin.Context) {
	db := common.GetGorm()
	var mList []menusOptions

	r1 := db.Model(&common.TMenu{}).Where("isNull(parent_id)").Find(&mList)
	if r1.Error != nil && r1.Error != gorm.ErrRecordNotFound {
		logger.Error(r1.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for i, m := range mList {
		var child []menusOptions
		r2 := db.Model(&common.TMenu{}).Where("parent_id = ?", m.ID).Find(&child)
		if r2.Error != nil && r2.Error != gorm.ErrRecordNotFound {
			logger.Error(r2.Error.Error())
			Response(ctx, errorcode.Fail, nil, false, "系统异常")
			return
		}
		mList[i].Children = child
	}
	Response(ctx, errorcode.Success, mList, true, "操作成功")
}

func (m *Menu) ListUserMenus(ctx *gin.Context) {
	userid, exist := ctx.Get("a_userid")
	if !exist {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	db := common.GetGorm()
	userAuth := common.TUserAuth{}
	r := db.Where("id = ?", userid).First(&userAuth)
	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		logger.Error(r.Error.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	if r.Error != nil {
		Response(ctx, errorcode.UsernameNotExist, nil, false, "用户名不存在")
		return
	}
	sql := "select * from v_user_menu where user_id = ?"
	userMenu := make([]common.VUserMenu, 0)
	rows, err := db.Raw(sql, userAuth.UserInfoId).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err.Error())
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	for rows.Next() {
		var t common.VUserMenu
		_ = db.ScanRows(rows, &t)
		userMenu = append(userMenu, t)
	}
	// 转换菜单格式为前端接受的格式
	data := common.ConvertMenuType(userMenu)
	Response(ctx, errorcode.Success, data, true, "操作成功")
}
