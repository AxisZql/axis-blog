package common

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

/*
* @author:AxisZql
* @date:2022-3-14
* @desc:读取配置文件
 */

// Configure 配置文件结构体
type Configure struct {
	App struct {
		Port       int64  `json:"port" remark:"http端口" must:"true"`
		InitModels bool   `json:"init_models" remark:"是否初始化模型" must:"false"`
		HostName   string `json:"host_name" remark:"部署域名" must:"false"`
		Loglevel   string `json:"loglevel" remark:"日志级别" must:"true"`
		Logfile    string `json:"logfile" remark:"日志文件" must:"true"`
		AvatarDir  string `json:"avatar_dir" remark:"头像目录" must:"ture"`
		ArticleDir string `json:"article_dir" remark:"文章图片目录" must:"ture"`
		VoiceDir   string `json:"voice_dir" remark:"语音目录" must:"true"`
		PhotoDir   string `json:"photo_dir" remark:"照片目录" must:"true"`
		ConfigDir  string `json:"config_dir" remark:"配置图片目录" must:"true"`
		TalkDir    string `json:"talk_dir" remark:"说说图片目录" must:"true"`
	}

	Db struct {
		Type     string `json:"type" remark:"数据库类型" must:"true"`
		Host     string `json:"host" remark:"数据库IP" must:"true"`
		Port     string `json:"port" remark:"数据库端口" must:"true"`
		Username string `json:"username" remark:"数据库用户名" must:"true"`
		Password string `json:"password" remark:"数据库密码" must:"true"`
		DbName   string `json:"dbname" remark:"数据库名称" must:"true"`
	}

	Redis struct {
		Host     string `json:"host" remark:"ip" must:"true"`
		Port     string `json:"port" remark:"redis端口" must:"true"`
		Db       int    `json:"db" remark:"数据库编号" must:"false"`
		Password string `json:"password" remark:"密码" must:"true"`
	}

	RabbitMq struct {
		Host     string `json:"host" remark:"MQ ip" must:"false"`
		Port     string `json:"port" remark:"MQ 端口" must:"false"`
		Username string `json:"username" remark:"MQ 用户名" must:"false"`
		Password string `json:"password" remark:"MQ 密码" must:"false"`
	}

	Jwt struct {
		Key string `json:"key" remark:"验证密钥" must:"true"`
	}

	Mail struct {
		Host     string `json:"host" remark:"邮箱服务器" must:"false"`
		Port     int    `json:"port" remark:"邮箱服务端口" must:"false"`
		Username string `json:"username" remark:"邮箱用户名" must:"false"`
		Password string `json:"password" remark:"邮箱授权码" must:"false"`
	}

	Oss struct {
		Host            string `json:"host" remark:"oss域名" must:"false"`
		AccessKeyId     string `json:"access_key_id" remark:"访问密钥id" must:"false"`
		AccessKeySecret string `json:"access_key_secret" remark:"访问密钥密码" must:"false"`
		BucketName      string `json:"bucket_name" remark:"bucket名称" must:"false"`
	}

	QQ struct {
		Appid         string `json:"appid" remark:"QQ appid" must:"false"`
		CheckTokenUrl string `json:"check_token_url" remark:"校验token地址" must:"false"`
		UserInfoUrl   string `json:"user_info_url" remark:"QQ用户信息地址" must:"false"`
	}

	Weibo struct {
		Appid          string `json:"appid" remark:"微博appid" must:"false"`
		AppSecret      string `json:"app_secret" remark:"微博appSecret" must:"false"`
		GrantType      string `json:"grant_type" remark:"微博登陆类型" must:"false"`
		RedirectUrl    string `json:"redirect_url" remark:"微博回调域名" must:"false"`
		AccessTokenUrl string `json:"access_token_url" remark:"微博访问令牌地址" must:"false"`
		UserInfoUrl    string `json:"user_info_url" remark:"微博用户信息地址" must:"false"`
	}

	Ip struct {
		AppKey  string `json:"app_key" remark:"ip api key" must:"false"`
		AppCode string `json:"app_code" remark:"ip api code" must:"false"`
	}

	Es struct {
		Addr string `json:"addr" remark:"elasticsearch server address" must:"false"`
	}
}

func InitViper() {
	viper.SetConfigName("config") //指定配置文件的文件名称,无需扩展名
	viper.AddConfigPath(".")      //设定配置文件的路径
	viper.AutomaticEnv()          //自动从环境变量中 读取匹配的参数

	// 读取-c输入的路径参数，初始化配置文件 example:./main -c config.yaml
	if len(os.Args) >= 3 {
		if os.Args[1] == "-c" {
			cfgFile := os.Args[2]
			viper.SetConfigFile(cfgFile)
		}
	}

	// 加载对应的配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("加载配置文件出错:%v", err)
	}

	file := viper.GetViper().ConfigFileUsed()
	configData, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("加载配置文件出错:%v", err)
	}
	configText := strings.ReplaceAll(string(configData), "&", "|||")
	r := strings.NewReader(configText)
	err = viper.ReadConfig(r)
	if err != nil {
		log.Fatalf("初始化配置文件出错:%v", err)
	}

}

var Conf = &Configure{}

// InitConfigure 初始化配置
func InitConfigure() (err error) {
	InitViper()
	confValue := reflect.ValueOf(Conf).Elem()
	confType := reflect.TypeOf(*Conf)

	for i := 0; i < confType.NumField(); i++ {
		section := confType.Field(i)
		sectionValue := confValue.Field(i)

		//读取节类型信息
		for j := 0; j < section.Type.NumField(); j++ {
			key := section.Type.Field(j)
			keyValue := sectionValue.Field(j)

			sec := strings.ToLower(section.Name) //配置文件节名小写（为了和config.toml对应位置匹配）
			remark := key.Tag.Get("remark")      //读取配置备注
			must := key.Tag.Get("must")
			tag := key.Tag.Get("json")
			if tag == "" {
				err := fmt.Errorf("can not found a tag name `json` in struct of [%s].[%s]", sec, tag)
				logger.Error(err.Error())
			}

			//绑定环境变量，会优先使用环境变量的值
			logger.Info(fmt.Sprintf("绑定环境变量 AXIS_%s_%s ==> %s.%s", strings.ToUpper(sec), strings.ToUpper(tag), sec, tag))
			envKey := fmt.Sprintf("AXIS_%s_%s", strings.ToUpper(sec), strings.ToUpper(tag))
			_ = viper.BindEnv(sec+"."+tag, envKey)

			//根据类型识别配置字段
			switch key.Type.Kind() {
			case reflect.String:
				value := viper.GetString(sec + "." + tag)
				if value == "" && must != "false" {
					err = fmt.Errorf("get a blank value of must item [%s].%s.%s", sec, tag, remark)
					logger.Error(err.Error())
					return err
				}
				keyValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := viper.GetInt64(sec + "." + tag)
				if value == 0 && must != "false" {
					err = fmt.Errorf("get a zero value of must item [%s].%s %s", sec, tag, remark)
					logger.Error(err.Error())
					return err
				}
				keyValue.SetInt(value)

			case reflect.Bool:
				value := viper.GetBool(sec + "." + tag)
				keyValue.SetBool(value)

			case reflect.Slice:
				value := viper.GetStringSlice(sec + "." + tag)
				val := reflect.ValueOf(&value)
				keyValue.Set(val.Elem())

			default:
				logger.Warn(fmt.Sprintf("unsupported config struct key type %T", key.Type.Kind()))
			}

		}

	}
	return err
}
