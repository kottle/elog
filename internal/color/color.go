package color

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func Blue(msg string) string {
	return blue + msg + reset
}

func Red(msg string) string {
	return red + msg + reset
}
func Green(msg string) string {
	return green + msg + reset
}
func Yellow(msg string) string {
	return yellow + msg + reset
}
func Purple(msg string) string {
	return purple + msg + reset
}
func Cyan(msg string) string {
	return cyan + msg + reset
}
func Gray(msg string) string {
	return gray + msg + reset
}
func White(msg string) string {
	return white + msg + reset
}
