package services

import (
	"context"
	"fmt"

	"github.com/fpay/foundation-go"

	"github.com/fpay/foundation-go/log"

	"github.com/owenliu1122/notice"
)

// NewSenderService return target channel service interface.
func NewSenderService(logger *log.Logger, msgType string, sendToolCfg notice.SendServiceConfig, pc foundation.JobManager) SenderService {

	switch msgType {
	case NoticeTypeWeChat:
		return NewWeChatSenderService(logger, sendToolCfg, pc)
	case NoticeTypeMail:
		return NewMailSenderService(logger, sendToolCfg, pc)
	case NoticeTypePhone:
		return NewPhoneSenderService(logger, sendToolCfg, pc)
	default:
		panic(fmt.Sprintf("Unknown MsgType: %s", msgType))
	}
}

// SenderService sender message to target channel service interface.
type SenderService interface {
	Handler(ctx context.Context, msg *notice.UserMessage) error
}
