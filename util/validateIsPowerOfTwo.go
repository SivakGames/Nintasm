package util

//Is the value a power of 2 *and* not 0?
func ValidateIsPowerOfTwo(n int) bool {
	return (n != 0) && (n&(n-1) == 0)
}
