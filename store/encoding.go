package store

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n uint64) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0)
	for n > 0 {
		rem := n % 62
		b = append([]byte{base62Chars[rem]}, b...)
		n = n / 62
	}
	return string(b)
}
