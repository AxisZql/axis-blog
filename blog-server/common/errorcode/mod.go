package errorcode

/*
* @author:AxisZql
* @date:2022-3-16 11:19 AM
* @desc:业务错误码描述
 */

type ErrCode = int64

// Success 操作成功
const Success ErrCode = 20000

// AuthorizedError 没有操作权限
const AuthorizedError ErrCode = 40300

// ExpireLoginTime 登陆过期
const ExpireLoginTime ErrCode = 40100

// SystemError 系统异常
const SystemError ErrCode = 50000

// Fail  失败
const Fail ErrCode = 51000

// ValidError 参数校验失败
const ValidError ErrCode = 52000

// NotFoundResource 找不到资源
const NotFoundResource ErrCode = 404000

// UsernameExistError 用户名已经存在
const UsernameExistError ErrCode = 52001

// UsernameNotExist 用户名不存在
const UsernameNotExist ErrCode = 52002

// QQLoginError QQ登陆错误
const QQLoginError ErrCode = 53001

// WeiboLoginError 微博登陆错误
const WeiboLoginError ErrCode = 53002
