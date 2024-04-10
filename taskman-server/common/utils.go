package common

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strconv"
	"time"
)

func GetLowVersionUnixMillis(date string) string {
	var t time.Time
	loc, _ := time.LoadLocation("Local")
	t, _ = time.ParseInLocation(models.DateTimeFormat, date, loc)
	millisecondsSinceEpoch := t.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).Nanoseconds() / 1e6 // 计算从1970年1月1日起经过的微秒数，再除以1000得到毫秒数
	return fmt.Sprintf("%d", millisecondsSinceEpoch)
}

func BuildVersionNum(version string) string {
	if version == "" {
		return "v1"
	}
	tmpV, err := strconv.Atoi(version[1:])
	if err != nil {
		return fmt.Sprintf("v%d", time.Now().Unix())
	}
	return fmt.Sprintf("v%d", tmpV+1)
}

func CompareUpdateConfirmTime(updatedTime, confirmTime string) bool {
	result := false
	ut, _ := time.Parse(models.DateTimeFormat, updatedTime)
	ct, _ := time.Parse(models.DateTimeFormat, confirmTime)
	if ut.Unix() > ct.Unix() {
		result = true
	}
	return result
}
