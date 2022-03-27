package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Menu struct{}

func (m *Menu) ListMenus(*gin.Context)        {}
func (m *Menu) SaveOrUpdateMenu(*gin.Context) {}
func (m *Menu) DeleteMenu(*gin.Context)       {}
func (m *Menu) ListMenuOptions(*gin.Context)  {}

func (m *Menu) ListUserMenus(ctx *gin.Context) {
	_session, err := Store.Get(ctx.Request, "CurUser")
	if err != nil {
		Response(ctx, errorcode.Fail, nil, false, "系统异常")
		return
	}
	userid := _session.Values["a_userid"]
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
