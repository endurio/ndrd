package chainutil_test

import (
	"fmt"
	"math"

	"github.com/endurio/ndrd/types"
)

func ExampleAmount() {

	a := types.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = types.Amount(1e8)
	fmt.Println("100,000,000 Atom:", a)

	a = types.Amount(1e5)
	fmt.Println("100,000 Atom:", a)
	// Output:
	// Zero Atom: 0 Coin
	// 100,000,000 Atom: 1 Coin
	// 100,000 Atom: 0.001 Coin
}

func ExampleNewAmount() {
	amountOne, err := types.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := types.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := types.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := types.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 Coin
	// 0.01234567 Coin
	// 0 Coin
	// invalid bitcoin amount
}

func ExampleAmount_unitConversions() {
	amount := types.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(chainutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(chainutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(chainutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(chainutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kCoin
	// Atom to Coin: 444333.222111 Coin
	// Atom to MilliCoin: 444333222.111 mCoin
	// Atom to MicroCoin: 444333222111 μCoin
	// Atom to Atom: 44433322211100 Atom
}
