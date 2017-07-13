package digest

// Gets the digest of MD5
type Digestor interface {
	GetDigest() []byte
}

// Function object for Digestor
type DigestorFunc func() []byte

func (f DigestorFunc) GetDigest() []byte {
	return f()
}

type SumFunc func([]byte) []byte

// Sums all of the digestors with customized sum function
func SumAll(sumFunc SumFunc, rest ...Digestor) []byte {
	allBytes := make([]byte, 0)

	for _, digestor := range rest {
		allBytes = append(allBytes, digestor.GetDigest()...)
	}

	return sumFunc(allBytes)
}
