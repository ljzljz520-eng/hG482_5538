package excel

import "errors"

var (
	ErrMissingSheet    = errors.New("回访工作簿缺少回访记录工作表")
	ErrMissingHeaders  = errors.New("回访工作簿缺少必要表头")
	ErrInvalidWorkbook = errors.New("回访工作簿无法读取")
)

const FollowUpSheet = "回访记录"

var RequiredHeaders = []string{
	"记录编号", "客户编号", "客户姓名", "服务类型", "阿姨", "服务日期",
	"下次回访日期", "满意度", "满意度评价", "改进建议", "状态", "备注",
}
