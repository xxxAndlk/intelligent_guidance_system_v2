package entity

import "errors"

var (
	ErrInvalidPaymentID = errors.New("invalid payment ID")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrPaymentNotFound  = errors.New("payment not found")
)

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCash
	PaymentMethodWeChat
	PaymentMethodAlipay
	PaymentMethodCard
)

func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodCash:
		return "现金"
	case PaymentMethodWeChat:
		return "微信支付"
	case PaymentMethodAlipay:
		return "支付宝"
	case PaymentMethodCard:
		return "银行卡"
	default:
		return "未知"
	}
}

func PaymentMethodFromCode(code string) PaymentMethod {
	switch code {
	case "CASH":
		return PaymentMethodCash
	case "WECHAT":
		return PaymentMethodWeChat
	case "ALIPAY":
		return PaymentMethodAlipay
	case "CARD":
		return PaymentMethodCard
	default:
		return PaymentMethodUnknown
	}
}

func (m PaymentMethod) Code() string {
	switch m {
	case PaymentMethodCash:
		return "CASH"
	case PaymentMethodWeChat:
		return "WECHAT"
	case PaymentMethodAlipay:
		return "ALIPAY"
	case PaymentMethodCard:
		return "CARD"
	default:
		return "UNKNOWN"
	}
}

type PaymentStatus int

const (
	PaymentStatusUnknown PaymentStatus = iota
	PaymentStatusPending
	PaymentStatusProcessing
	PaymentStatusSuccess
	PaymentStatusFailed
	PaymentStatusRefunded
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusPending:
		return "待支付"
	case PaymentStatusProcessing:
		return "处理中"
	case PaymentStatusSuccess:
		return "支付成功"
	case PaymentStatusFailed:
		return "支付失败"
	case PaymentStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

func PaymentStatusFromCode(code string) PaymentStatus {
	switch code {
	case "PENDING":
		return PaymentStatusPending
	case "PROCESSING":
		return PaymentStatusProcessing
	case "SUCCESS":
		return PaymentStatusSuccess
	case "FAILED":
		return PaymentStatusFailed
	case "REFUNDED":
		return PaymentStatusRefunded
	default:
		return PaymentStatusUnknown
	}
}

func (s PaymentStatus) Code() string {
	switch s {
	case PaymentStatusPending:
		return "PENDING"
	case PaymentStatusProcessing:
		return "PROCESSING"
	case PaymentStatusSuccess:
		return "SUCCESS"
	case PaymentStatusFailed:
		return "FAILED"
	case PaymentStatusRefunded:
		return "REFUNDED"
	default:
		return "UNKNOWN"
	}
}

func (s PaymentStatus) CanRefund() bool {
	return s == PaymentStatusSuccess
}

func (s PaymentStatus) IsTerminal() bool {
	return s == PaymentStatusSuccess || s == PaymentStatusFailed || s == PaymentStatusRefunded
}