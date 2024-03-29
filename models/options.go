package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jinfagang/picbed/common"
	"github.com/jinfagang/picbed/configs"
)

type Option struct {
	Model

	OptionsKey   string
	OptionsValue string `gorm:"type:text;"`
	Group        string
}

var commonPrefix = "Option_"

func SetOption(key string, value interface{}, group ...string) bool {
	//首先载入缓存
	cacheKey := commonPrefix + group[0] + "_" + key
	s := formatValue(value)
	config.Cache.Put(cacheKey, s, time.Hour)
	//然后载入数据库
	err := db.Create(&Option{
		OptionsKey:   key,
		OptionsValue: s,
		Group:        group[0],
	}).Error
	return modelsError(common.ErrorToString(err))
}

func UpdateOption(key string, value interface{}, group ...string) bool {
	//首先载入缓存key
	cacheKey := commonPrefix + group[0] + "_" + key
	s := formatValue(value)
	//然后载入数据库
	err := db.Model(&Option{}).
		Where(&Option{OptionsKey: key, Group: group[0]}).Update("options_value", s).Error
	//最后载入缓存
	config.Cache.Delete(cacheKey)
	config.Cache.Put(cacheKey, s, time.Hour)
	return modelsError(common.ErrorToString(err))
}

func GetOption(key string, group ...string) string {
	cacheKey := getCacheKey(group[0], key)
	//尝试从缓存中检索
	v := config.Cache.Get(cacheKey)
	if v != nil {
		return formatValue(v)
	}
	var option Option
	//从数据库中查询
	err := dbOptionSelect(key, group[0]).First(&option).Error
	if option.ID <= 0 {
		modelsError(common.ErrorToString(err))
		return ""
	}
	//cache
	config.Cache.Put(cacheKey, formatValue(option.OptionsValue), time.Hour)
	return formatValue(option.OptionsValue)
}

func getCacheKey(group, key string) string {
	return commonPrefix + group + "_" + key
}

func formatValue(value interface{}) string {
	//检测是否可以被转s为 string 类型
	s, ok := value.(string)
	if !ok {
		//进行json 序列化
		b, err := json.Marshal(value)
		if err != nil {
			modelsError(common.ErrorToString(err))
		}
		//转为 string
		s = string(b)
	}
	return s
}

func dbOptionSelect(key, group string) *gorm.DB {
	return db.Model(&Option{}).
		Where(&Option{OptionsKey: key, Group: group})
}

func MigrateOptions() error {
	if db.HasTable(&Option{}) {
		err := db.DropTable(&Option{}).Error
		err = db.CreateTable(&Option{}).Error
		CreateDefaultOptions()
		return err
	} else {
		err := db.CreateTable(&Option{}).Error
		CreateDefaultOptions()
		return err
	}

}

func CreateDefaultOptions() {
	//Sina
	var sina = common.SinaAccount{
		Proxy:               common.ProxyConf{Status: false, Node: "https://images.weserv.nl/?url="},
		DefultPicSize:       "large",
		ResetSinaCookieTime: 3600,
	}
	SetOption("sina", sina, "conf")
	//Flickr
	SetOption("flickr", config.SiteConfig.SiteUploadWay.FlickrAccount, "conf")
	//Local
	SetOption("local", config.SiteConfig.SiteUploadWay.LocalStore, "conf")
	//Imgur
	var imgur = common.ImgurAccount{
		Status:  false,
		ImgurID: "",
		Proxy:   common.ProxyConf{Status: false, Node: "https://images.weserv.nl/?url="},
	}
	SetOption("imgur", imgur, "conf")
	//Github
	var github = common.GithubAccount{
		AccessToken: "",
		Owner:       "",
		Repo:        "",
		Proxy:       common.ProxyConf{Status: false, Node: "https://images.weserv.nl/?url="},
	}

	SetOption("github", github, "conf")

	//Gitee
	var gitee = common.GiteeAccount{
		AccessToken: "",
		Owner:       "",
		Repo:        "",
	}
	SetOption("gitee", gitee, "conf")
	SetOption("proxy", "https://images.weserv.nl/?url=", "conf")
	//Set Root
	rId := GetStoreInfoByApi("local")
	var dispatch = common.Dispatch{
		Status:    true,
		Root:      "local",
		RootID:    rId.ID,
		CacheTime: time.Second * 3600 * 2,
	}
	SetOption("dispatch", dispatch, "conf")
	//Api set
	apiOption := common.ApiOptions{
		Status:     false,
		Auth:       true,
		ApiDefault: "SouGou",
	}
	SetOption("api_option", apiOption, "conf")

	//Limit
	var limit = common.LimitConfig{
		Status:    false, //是否启用
		BlockTime: 3600,  //暂时封禁的时间
		DeadLine:  10,    //最大次数
		AllowNum:  100,   //张
		AllowTime: 3600,  //s
	}
	SetOption("ip_limit", limit, "conf")
	//Site
	var site = common.SiteBase{}
	site.SiteName = "AuXpI API 图床"
	site.SiteUrl = "http://yoursite.com/"
	site.SiteFooter = "新一代图床"
	site.Logo = "/static/app/images/logo.jpg"
	site.SiteUploadMaxSize = 5
	site.SiteUploadMaxNumber = 10
	site.AllowTourists = false
	site.AllowRegister = false
	site.JwtSecret = config.GenerateUniqueString()
	site.JwtDueTime = 3
	site.MailConfig.Status = false
	SetOption("site_base", site, "conf")
}
