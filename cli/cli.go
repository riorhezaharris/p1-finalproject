package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"p1finalproject/handler"
	"p1finalproject/utils"
)

type CLI struct {
	Handler *handler.Handler
	scanner *bufio.Scanner
	userID  int
	email   string
}

func NewCLI(h *handler.Handler) *CLI {
	return &CLI{
		Handler: h,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (c *CLI) Init() {
	// Banner
	shirt := []string{
		"   ___ ___   ",
		" /| |/|\\| |\\ ",
		"/_| ´ |.` |_\\",
		"  |   |.  |  ",
		"  |   |.  |  ",
		"  |___|.__|  ",
	}
	logo := []string{
		"      _       _   _               ",
		"     | |     | | | |              ",
		"  ___| | ___ | |_| |__   ___  ___ ",
		" / __| |/ _ \\| __| '_ \\ / _ \\/ __|",
		"| (__| | (_) | |_| | | |  __/\\__ \\",
		" \\___|_|\\___/ \\__|_| |_|\\___||___/",
	}
	for i := 0; i < len(shirt); i++ {
		fmt.Println(utils.ClrYellow + shirt[i] + utils.ClrMagenta + logo[i] + utils.ClrReset)
	}
	fmt.Println(utils.Colorize("==== Welcome To Cloth Store ====", utils.ClrBold+utils.ClrCyan))
	c.login()
	for {
		c.mainMenu()
	}
}

// --- Login / Skip Auth ---
func (c *CLI) login() {
	if os.Getenv("SKIP_AUTH") == "1" {
		email := "demo@store.local"
		u, err := c.Handler.GetOrCreateUserByEmail(email)
		if err != nil {
			fmt.Println(utils.Colorize("Gagal auto-login:", utils.ClrRed), err)
			os.Exit(1)
		}
		c.userID = u.ID
		c.email = u.Email
		fmt.Println(utils.Colorize("Mode SKIP_AUTH aktif → auto login sebagai ", utils.ClrYellow) + c.email)
		return
	}

	fmt.Print("Masukkan email untuk mulai (contoh: user@example.com): ")
	c.scanner.Scan()
	email := strings.TrimSpace(c.scanner.Text())
	u, err := c.Handler.GetOrCreateUserByEmail(email)
	if err != nil {
		fmt.Println(utils.Colorize("Gagal login:", utils.ClrRed), err)
		os.Exit(1)
	}
	c.userID = u.ID
	c.email = u.Email
	fmt.Println("Halo,", utils.Colorize(c.email, utils.ClrGreen))
}

func (c *CLI) mainMenu() {
	fmt.Println("")
	fmt.Println(utils.Colorize("1) Berbelanja", utils.ClrBold), " ", utils.Colorize("2) Transaction history", utils.ClrBold), " ", utils.Colorize("0) Keluar", utils.ClrBold))
	fmt.Print(utils.Colorize("Pilih: ", utils.ClrCyan))
	c.scanner.Scan()
	switch strings.TrimSpace(c.scanner.Text()) {
	case "1":
		c.shopFlow()
	case "2":
		c.history()
	case "0":
		fmt.Println(utils.Colorize("Terima kasih sudah berbelanja!", utils.ClrGreen))
		os.Exit(0)
	default:
		fmt.Println(utils.Colorize("Pilihan tidak valid.", utils.ClrRed))
	}
}

// --------- Flow belanja ----------
func (c *CLI) shopFlow() {
	for {
		// 2. Show product list
		c.showProducts()

		// 3. Input product code
		fmt.Print(utils.Colorize("Mau beli apa? Input Product Code (atau 0 untuk cek keranjang): ", utils.ClrCyan))
		c.scanner.Scan()
		codeStr := strings.TrimSpace(c.scanner.Text())
		if codeStr == "0" {
			c.cartMenu()
			return
		}
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			fmt.Println(utils.Colorize("Kode produk tidak valid.", utils.ClrRed))
			continue
		}

		// 4. Info size dari product
		p, err := c.Handler.GetProductByCode(code)
		if err != nil {
			fmt.Println(utils.Colorize("Produk tidak ditemukan: ", utils.ClrRed), err)
			continue
		}
		fmt.Printf("%s %s | Size: %s | Harga: %s\n",
			utils.Colorize("Produk:", utils.ClrYellow), utils.Colorize(p.Name, utils.ClrBold), p.SizeName, utils.FormatRupiah(p.Price))

		// 5. Jumlah
		fmt.Print(utils.Colorize("Mau jumlah berapa: ", utils.ClrCyan))
		c.scanner.Scan()
		qtyStr := strings.TrimSpace(c.scanner.Text())
		qty, parseErr := utils.ParseIntFlexible(qtyStr)
		if parseErr != nil {
			fmt.Println(utils.Colorize("Quantity tidak valid.", utils.ClrRed))
			continue
		}
		if qty <= 0 {
			fmt.Println(utils.Colorize("Jumlah harus > 0", utils.ClrRed))
			continue
		}

		if err := c.Handler.AddToCart(c.userID, p.ID, qty); err != nil {
			fmt.Println(utils.Colorize("Gagal menambah ke keranjang:", utils.ClrRed), err)
			continue
		}

		fmt.Println(utils.Colorize("Barang dimasukkan ke keranjang.", utils.ClrGreen))
	}
}

func (c *CLI) showProducts() {
	fmt.Println("\n" + utils.Colorize("=== Product List ===", utils.ClrBold+utils.ClrBlue))
	items, err := c.Handler.ListProducts()
	if err != nil {
		fmt.Println(utils.Colorize("Tidak bisa mengambil produk:", utils.ClrRed), err)
		return
	}

	headers := []string{"ID", "Name", "Size", "Price"}
	var rows [][]string
	for _, p := range items {
		rows = append(rows, []string{
			strconv.Itoa(p.ID),
			p.Name,
			p.SizeName,
			utils.FormatRupiah(p.Price),
		})
	}
	printTable(headers, rows)
}

func (c *CLI) cartMenu() {
	for {
		items, total, err := c.Handler.GetCart(c.userID)
		if err != nil {
			fmt.Println(utils.Colorize("Gagal mengambil keranjang:", utils.ClrRed), err)
			return
		}
		fmt.Println("\n" + utils.Colorize("=== Keranjang Kamu ===", utils.ClrBold+utils.ClrBlue))
		if len(items) == 0 {
			fmt.Println(utils.Colorize("(kosong)", utils.ClrGray))
		} else {
			headers := []string{"Code", "Name", "Size", "Price", "Qty", "Subtotal"}
			var rows [][]string
			for _, it := range items {
				rows = append(rows, []string{
					strconv.Itoa(it.ProductID),
					it.Name,
					it.SizeName,
					utils.FormatRupiah(it.Price),
					strconv.Itoa(it.Quantity),
					utils.FormatRupiah(it.SubTotal),
				})
			}
			printTable(headers, rows)
			fmt.Println(utils.Colorize("Total: "+utils.FormatRupiah(total), utils.ClrBold+utils.ClrGreen))
		}

		fmt.Print(utils.Colorize("Input product code untuk edit/delete, ketik 'checkout' untuk membuat order, atau 'back' untuk belanja lagi: ", utils.ClrCyan))
		c.scanner.Scan()
		cmd := strings.TrimSpace(strings.ToLower(c.scanner.Text()))

		if cmd == "checkout" {
			c.checkout(total)
			return
		}
		if cmd == "back" {
			return
		}
		pid, err := strconv.Atoi(cmd)
		if err != nil {
			fmt.Println(utils.Colorize("Perintah tidak dikenali.", utils.ClrRed))
			continue
		}
		fmt.Print(utils.Colorize("Masukkan quantity baru (0 untuk hapus): ", utils.ClrCyan))
		c.scanner.Scan()
		qtyStr := strings.TrimSpace(c.scanner.Text())
		qty, parseErr := utils.ParseIntFlexible(qtyStr)
		if parseErr != nil {
			fmt.Println(utils.Colorize("Quantity tidak valid.", utils.ClrRed))
			continue
		}
		if err := c.Handler.UpdateCartItem(c.userID, pid, qty); err != nil {
			fmt.Println(utils.Colorize("Gagal mengubah keranjang:", utils.ClrRed), err)
		} else {
			if qty == 0 {
				fmt.Println(utils.Colorize("Item dihapus dari keranjang.", utils.ClrYellow))
			} else {
				fmt.Println(utils.Colorize("Quantity diperbarui.", utils.ClrGreen))
			}
		}
	}
}

func (c *CLI) checkout(currentTotal float64) {
	items, total, err := c.Handler.GetCart(c.userID)
	if err != nil {
		fmt.Println(utils.Colorize("Gagal mengambil keranjang:", utils.ClrRed), err)
		return
	}
	if len(items) == 0 {
		fmt.Println(utils.Colorize("Keranjang kosong.", utils.ClrYellow))
		return
	}
	fmt.Println(utils.Colorize("Total yang harus dibayar: "+utils.FormatRupiah(total), utils.ClrBold+utils.ClrGreen))
	// fmt.Print(utils.Colorize("Masukkan nominal pembayaran: ", utils.ClrCyan))
	// c.scanner.Scan()
	paid, err := utils.ReadRupiahInteractive("Masukkan nominal pembayaran: ")
	if err != nil {
		fmt.Println(utils.Colorize("Input pembayaran dibatalkan / tidak valid.", utils.ClrRed), err)
		return
	}

	res, err := c.Handler.Checkout(c.userID, paid)
	if err != nil {
		fmt.Println(utils.Colorize("Checkout gagal:", utils.ClrRed), err)
		return
	}
	fmt.Println("\n" + utils.Colorize("=== Checkout Result ===", utils.ClrBold+utils.ClrBlue))
	headers := []string{"Order ID", "Total", "Dibayar", "Kembalian", "Order Status", "Payment"}
	rows := [][]string{{
		strconv.Itoa(res.OrderID),
		utils.FormatRupiah(res.Total),
		utils.FormatRupiah(res.Paid),
		func() string {
			if res.Paid >= res.Total {
				return utils.FormatRupiah(res.Change)
			}
			return "-"
		}(),
		res.OrderStatus,
		res.PayStatus,
	}}
	printTable(headers, rows)

	c.history()
}

func (c *CLI) history() {
	fmt.Println("\n" + utils.Colorize("=== Transaction History ===", utils.ClrBold+utils.ClrBlue))
	list, err := c.Handler.GetOrderHistory(c.userID)
	if err != nil {
		fmt.Println(utils.Colorize("Gagal mengambil history:", utils.ClrRed), err)
		return
	}
	if len(list) == 0 {
		fmt.Println(utils.Colorize("(belum ada transaksi)", utils.ClrGray))
		return
	}

	headers := []string{"Order#", "Created At", "Total", "Status", "Paid"}
	var rows [][]string
	for _, o := range list {
		rows = append(rows, []string{
			"#" + strconv.Itoa(o.ID),
			o.CreatedAt,
			utils.FormatRupiah(o.TotalPrice),
			o.Status,
			utils.FormatRupiah(o.Payment),
		})
	}
	printTable(headers, rows)
}

// ======================================================
// Tabel rapi: auto ukur lebar kolom agar header & rows lurus
// ======================================================

func strWidth(s string) int { return utf8.RuneCountInString(s) }

func repeat(ch string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(ch, n)
}

// printTable: header bold, garis pemisah, angka/uang rata kanan
func printTable(headers []string, rows [][]string) {
	col := len(headers)
	widths := make([]int, col)
	for i, h := range headers {
		if w := strWidth(h); w > widths[i] {
			widths[i] = w
		}
	}
	for _, r := range rows {
		for i := 0; i < col && i < len(r); i++ {
			if w := strWidth(r[i]); w > widths[i] {
				widths[i] = w
			}
		}
	}

	padding := 2

	// header
	for i, h := range headers {
		fmt.Print(utils.Colorize(padRight(h, widths[i]), utils.ClrBold))
		if i < col-1 {
			fmt.Print(repeat(" ", padding))
		}
	}
	fmt.Println()

	// separator
	sep := ""
	for i := 0; i < col; i++ {
		if i > 0 {
			sep += repeat(" ", padding)
		}
		sep += repeat("─", widths[i])
	}
	fmt.Println(utils.Colorize(sep, utils.ClrGray))

	// rows
	for _, r := range rows {
		for i := 0; i < col; i++ {
			cell := ""
			if i < len(r) {
				cell = r[i]
			}
			if isNumericOrMoney(cell) {
				fmt.Print(padLeft(cell, widths[i]))
			} else {
				fmt.Print(padRight(cell, widths[i]))
			}
			if i < col-1 {
				fmt.Print(repeat(" ", padding))
			}
		}
		fmt.Println()
	}
}

func padRight(s string, w int) string {
	diff := w - strWidth(s)
	if diff > 0 {
		return s + repeat(" ", diff)
	}
	return s
}

func padLeft(s string, w int) string {
	diff := w - strWidth(s)
	if diff > 0 {
		return repeat(" ", diff) + s
	}
	return s
}

// Anggap "numeric" juga untuk format uang "Rp 1.234.567"
func isNumericOrMoney(s string) bool {
	if s == "" {
		return false
	}
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "Rp ") {
		s = strings.TrimPrefix(s, "Rp ")
		s = strings.ReplaceAll(s, ".", "")
	}
	for _, r := range s {
		if (r < '0' || r > '9') && r != '.' && r != ',' && r != '-' && r != '+' {
			return false
		}
	}
	return true
}
