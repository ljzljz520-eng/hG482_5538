package followup

import "errors"

var (
	ErrStartupFormat   = errors.New("回访文件格式错误，系统已阻止加载")
	ErrDuplicateRecord = errors.New("回访记录编号已存在")
	ErrRecordMissing   = errors.New("回访记录不存在")
)

type StartupNotice struct {
	Title    string
	Detail   string
	Blocking bool
}

func FormatNotice(err error) StartupNotice {
	if err == nil {
		return StartupNotice{Title: "回访簿已就绪", Detail: "可以开始处理今日回访", Blocking: false}
	}
	return StartupNotice{Title: "无法加载回访簿", Detail: err.Error(), Blocking: true}
}
