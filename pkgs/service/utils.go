package service

func MapToSlice(m map[int]int, s []int) {
	for count, value := range m {
		s[count] = value
	}
}
