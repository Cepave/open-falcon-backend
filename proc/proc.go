package proc

import (
	"sync/atomic"
)

var smsCount, mailCount, qqCount uint32

func GetSmsCount() uint32 {
	return atomic.LoadUint32(&smsCount)
}

func GetMailCount() uint32 {
	return atomic.LoadUint32(&mailCount)
}

func GetQQCount() uint32 {
	return atomic.LoadUint32(&qqCount)
}

func IncreSmsCount() {
	atomic.AddUint32(&smsCount, 1)
}

func IncreMailCount() {
	atomic.AddUint32(&mailCount, 1)
}

func IncreQQCount() {
	atomic.AddUint32(&qqCount, 1)
}
