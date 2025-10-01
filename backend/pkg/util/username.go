package util

func GenerateUsername() string {
	return "user_" + RandomString(6, "abcdefghijklmnopqrstuvwxyz0123456789")
}
