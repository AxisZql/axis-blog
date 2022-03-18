package service

import (
	"blog-server/common"
	"blog-server/common/errorcode"
	ctrl "blog-server/controllers"
	"github.com/gin-gonic/gin"
)

type Menu struct {
	ctrl.MenuHandle
}

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
	userid := _session.Values["userid"]
	db := common.GetGorm()
	sql := "select * from v_user_menu where ?"
	userMenu := make([]common.VUserMenu, 0)
	rows, err := db.Raw(sql, userid).Rows()
	for rows.Next() {
		var t common.VUserMenu
		_ = db.ScanRows(rows, &t)
		userMenu = append(userMenu, t)
	}
	// 转换菜单格式为前端接受的格式
	data := common.ConvertMenuType(userMenu)
	Response(ctx, errorcode.Success, data, true, "操作成功")
}
