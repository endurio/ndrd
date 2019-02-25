package chainutil_test

import (
	"fmt"
	"math"

)

func ExampleAmount() {

	a := chainutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = chainutil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)

	a = chainutil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Atom: 0 BTC
	// 100,000,000 Satoshis: 1 BTC
	// 100,000 Satoshis: 0.001 BTC
}

func ExampleNewAmount() {
	amountOne, err := chainutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := chainutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := chainutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := chainutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 BTC
	// 0.01234567 BTC
	// 0 BTC
	// invalid bitcoin amount
}

func ExampleAmount_unitConversions() {
	amount := chainutil.Amount(44433322211100)

	fmt.Println("Atom to kBTC:", amount.Format(chainutil.AmountKiloBTC))
	fmt.Println("Atom to BTC:", amount)
	fmt.Println("Atom to MilliBTC:", amount.Format(chainutil.AmountMilliBTC))
	fmt.Println("Atom to MicroBTC:", amount.Format(chainutil.AmountMicroBTC))
	fmt.Println("Atom to Atom:", amount.Format(chainutil.AmountSatoshi))

	// Output:
	// Atom to kBTC: 444.333222111 kBTC
	// Atom to BTC: 444333.222111 BTC
	// Atom to MilliBTC: 444333222.111 mBTC
	// Atom to MicroBTC: 444333222111 μBTC
	// Atom to Atom: 44433322211100 Atom
}
