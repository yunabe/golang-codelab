package unittest

// IntSum returns the sum of args.
func IntSum(args ...int) int {
	sum := 0
	for _, arg := range args {
		sum += arg
	}
	return sum
}
