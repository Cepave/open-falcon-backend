package utils

// Converts 64 bits integer(unsigned) to 32 bits one
func UintTo32(source []uint64) []uint32 {
	if source == nil {
		return nil
	}

	result := make([]uint32, len(source))

	for i, v := range source {
		result[i] = uint32(v)
	}

	return result
}

// Converts 64 bits integer(unsigned) to 16 bits one
func UintTo16(source []uint64) []uint16 {
	if source == nil {
		return nil
	}

	result := make([]uint16, len(source))

	for i, v := range source {
		result[i] = uint16(v)
	}

	return result
}

// Converts 64 bits integer(unsigned) to 8 bits one
func UintTo8(source []uint64) []uint8 {
	if source == nil {
		return nil
	}

	result := make([]uint8, len(source))

	for i, v := range source {
		result[i] = uint8(v)
	}

	return result
}

// Converts 64 bits integer to 32 bits one
func IntTo32(source []int64) []int32 {
	if source == nil {
		return nil
	}

	result := make([]int32, len(source))

	for i, v := range source {
		result[i] = int32(v)
	}

	return result
}

// Converts 64 bits integer to 16 bits one
func IntTo16(source []int64) []int16 {
	if source == nil {
		return nil
	}

	result := make([]int16, len(source))

	for i, v := range source {
		result[i] = int16(v)
	}

	return result
}

// Converts 64 bits integer to 8 bits one
func IntTo8(source []int64) []int8 {
	if source == nil {
		return nil
	}

	result := make([]int8, len(source))

	for i, v := range source {
		result[i] = int8(v)
	}

	return result
}
