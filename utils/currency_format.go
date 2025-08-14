package utils

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"os"
	"bufio"
	"fmt"

	"golang.org/x/term"
)

// FormatRupiah mengubah float64 jadi "Rp 1.234.567"
// Dibulatkan ke integer rupiah (tanpa koma)
func FormatRupiah(amount float64) string {
	n := int64(amount + 0.5) // round to nearest
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	s := itoa(n) // tanpa alloc fmt

	// sisipkan titik tiap 3 digit dari belakang
	var b strings.Builder
	b.Grow(len(s) + len(s)/3 + 4)
	b.WriteString(sign)
	b.WriteString("Rp ")
	for i, r := range s {
		// posisi dari depan: kalau sisa digit setelah posisi ini kelipatan 3, tambahkan titik
		if (len(s)-i)%3 == 0 && i != 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// itoa versi kecil (int64 -> desimal) agar tidak tarik fmt
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var a [20]byte
	i := len(a)
	for n > 0 {
		i--
		a[i] = byte('0' + n%10)
		n /= 10
	}
	return string(a[i:])
}

func ParseRupiahLikeToFloat(s string) (float64, error) {
	// ambil hanya digit
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	str := b.String()
	if str == "" {
		return 0, errors.New("tidak ada digit pada input")
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return float64(n), nil
}

// ParseIntFlexible: untuk qty "1.000" -> 1000, "002" -> 2
func ParseIntFlexible(s string) (int, error) {
	v, err := ParseRupiahLikeToFloat(strings.TrimSpace(s))
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// === Input uang dengan format live: Rp 1.234.567 saat user mengetik ===

func ReadRupiahInteractive(prompt string) (float64, error) {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, err
	}
	defer term.Restore(fd, oldState)

	r := bufio.NewReader(os.Stdin)
	digits := make([]rune, 0, 32)

	redraw := func() {
		// Clear full line + carriage return, lalu cetak ulang
		fmt.Print("\r\033[2K") // clear entire line
		fmt.Print(Colorize(prompt, ClrCyan) + "Rp " + formatDigitsWithDots(string(digits)))
	}

	redraw()

	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return 0, err
		}

		switch ch {
		case '\r', '\n': // Enter
			fmt.Print("\n")
			return ParseRupiahLikeToFloat(string(digits))

		case 3: // Ctrl+C
			fmt.Print("\n")
			return 0, fmt.Errorf("dibatalkan")

		case 127, 8: // Backspace (DEL atau BS)
			if len(digits) > 0 {
				digits = digits[:len(digits)-1]
			}
			redraw()

		case 27: // ESC — kemungkinan tombol khusus (arrow/delete)
			// Baca beberapa rune berikut, abaikan sequence seperti ESC [ 3 ~ (Delete)
			_ = r.UnreadRune() // taruh kembali
			seq, _ := r.Peek(3)
			// konsumsi minimal ESC
			_, _, _ = r.ReadRune() // ESC
			// konsumsi sisanya kalau ada
			for i := 0; i < len(seq); i++ {
				_, _, _ = r.ReadRune()
			}
			// tidak melakukan apa-apa, hanya mengabaikan
			redraw()

		default:
			if unicode.IsDigit(ch) {
				digits = append(digits, ch)
				redraw()
			}
			// karakter lain diabaikan
		}
	}
}

// sisipkan titik tiap 3 digit (untuk tampilan saat mengetik)
func formatDigitsWithDots(s string) string {
	n := len(s)
	if n == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(n + n/3)
	for i, r := range s {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}
