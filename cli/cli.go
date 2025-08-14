package cli

import (
	"bufio"
	"fmt"
	"os"
	"p1finalproject/handler"
	"strings"
)

type CLI struct {
	Handler handler.Handler
	scanner bufio.Scanner
}

func NewCli(h handler.Handler) *CLI {
	return &CLI{Handler: h, scanner: *bufio.NewScanner(os.Stdin)}
}

func (c *CLI) Init() {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Welcome to Sport Store")

	// // Scan for the input and trim the white space
	// fmt.Println("Enter email:")
	// var input string
	// if c.scanner.Scan() {
	// 	input = strings.TrimSpace(c.scanner.Text())
	// }
	// // Error handler for invalid input
	// if err := c.scanner.Err(); err != nil {
	// 	fmt.Println("Error reading input:", err)
	// }
	// if input == "" {
	// 	fmt.Println("Input is required, please try again.")
	// }
	// result, err := c.Handler.GetUser("bob.williams@example.com", "password2")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Log In User:")
	// userId := result.UserId
	// fmt.Println(userId)

	// cart, subtotal, err := c.Handler.GetCart(userId)
	// if err != nil {
	// 	panic(err)
	// }
	// if len(cart) == 0 {
	// 	fmt.Println("Cart is empty")
	// }
	// for _, CartItem := range cart {
	// 	fmt.Println(CartItem)
	// }
	// fmt.Println(subtotal)

	// c.Handler.AddItem(userId, 20)
	// c.Handler.AddItem(userId, 20)
	// c.Handler.AddItem(userId, 20)
	// c.Handler.RemoveItem(userId, 20)
	// c.Handler.RemoveItem(userId, 20)
	// c.Handler.RemoveItem(userId, 20)
	// c.Handler.AddItem(userId, 10)
	// c.Handler.ResetCart(userId)
	// c.Handler.AddItem(userId, 10)
	// cart, subtotal, err = c.Handler.GetCart(userId)
	// if err != nil {
	// 	panic(err)
	// }
	// if len(cart) == 0 {
	// 	fmt.Println("Cart is empty")
	// }
	// for _, CartItem := range cart {
	// 	fmt.Println(CartItem)
	// }
	// fmt.Println(subtotal)

	// orderId, err := c.Handler.CreateOrder(userId)
	// if err != nil {
	// 	panic(err)
	// }

	// err = c.Handler.CreatePayment(userId, orderId, 100)
	// if err != nil {
	// 	panic(err)
	// }

	// orders, err := c.Handler.GetOrders(userId)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, order := range orders {
	// 	fmt.Println(order)
	// }
}
