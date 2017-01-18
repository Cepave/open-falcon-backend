package digest

import (
	"time"
)

// This type is used to digest time value
type DigestableTime time.Time

// Generates digest by the unix value of time
//
// if time.Time.IsZero() is true, use NoByteFunc as the digest value
func (t DigestableTime) GetDigest() []byte {
	timeValue := time.Time(t)

	if timeValue.IsZero() {
		return NoByteFunc()
	}

	return GetBytesGetter(timeValue.Unix(), Md5SumFunc)()
}
