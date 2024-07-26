package util

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func YesNo(b bool) string {
	if b {
		return Yes
	}
	return No
}

func YESNO(b bool) string {
	if b {
		return YES
	}
	return NO
}
