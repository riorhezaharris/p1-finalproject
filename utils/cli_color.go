package utils

// ====== ANSI colors ======
const (
	ClrReset   = "\033[0m"
	ClrBold    = "\033[1m"
	ClrDim     = "\033[2m"
	ClrRed     = "\033[31m"
	ClrGreen   = "\033[32m"
	ClrYellow  = "\033[33m"
	ClrBlue    = "\033[34m"
	ClrMagenta = "\033[35m"
	ClrCyan    = "\033[36m"
	ClrGray    = "\033[90m"
	ClrOrange  = "\033[38;5;208m" // orange
	ClrPurple  = "\033[35m"       // ungu
)

func Colorize(s, color string) string { return color + s + ClrReset }

// =====================================
