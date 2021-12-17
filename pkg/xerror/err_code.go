package xerror

type Code uint32

const (
	ParamErr      Code = iota + 1000 // 参数错误
	SourceNotFind                    // 资源不存在
	SystemErr                        // 系统错误
	CopierErr                        // 复制错误
	WechatErr                        // 微信错误
	JwtErr                           // Token错误
	MarshalErr                       // 数据格式化错误
	SQLErr                           // 数据库错误
	RedisErr                         // Redis错误
	AuthErr                          // 请登录
	YunPianErr                       // 云片网错误
	EmailErr                         // 邮件发送失败
	Invalid                          // 非法请求
)

func (i Code) GetCode() uint32 {
	return uint32(i)
}
