package utils

func Length(s []byte) (i int) {
	var c byte
	for i, c = range s {
		if c == 0 {
			break
		}
	}
	return i
}
