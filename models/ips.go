package models

import (
	"github.com/jinfagang/picbed/common"
)

type IP struct {
	Model

	Content string
	Status  int //-1,0,1  || 永久 ，设定时间 block ，白名单
}

func CreateIpInfo(ip IP) bool {
	err := db.Create(&IP{
		Content: ip.Content,
		Status:  ip.Status,
	}).Error
	return modelsError(common.ErrorToString(err))
}

func DeleteIp(ip string) bool {
	err := db.Model(&IP{}).
		Where("content=?", ip).
		Delete(&IP{}).
		Error
	return modelsError(common.ErrorToString(err))
}
