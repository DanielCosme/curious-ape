package colors

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	purple = "\033[35m"
	cyan   = "\033[36m"
	gray   = "\033[37m"
	white  = "\033[97m"
)

func Red(s string) string {
	return red + s + reset
}

func Green(s string) string {
	return green + s + reset
}

func Yellow(s string) string {
	return yellow + s + reset
}

func Blue(s string) string {
	return blue + s + reset
}

func Purple(s string) string {
	return purple + s + reset
}

func Cyan(s string) string {
	return cyan + s + reset
}

func Gray(s string) string {
	return gray + s + reset
}

func White(s string) string {
	return white + s + reset
}
