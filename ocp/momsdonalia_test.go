package ocp

import (
	"io"
	"os"
	"testing"
)

func TestTakeOrder(t *testing.T) {
	m := MomsDonaliaKiosk{}

	tests := []struct {
		order    Order
		expected string
	}{
		{
			order: RegularOrder{
				PricePerUnit: 1.0,
				Quantity:     2,
			},
			expected: "MomsDonalia: 주문 합계는 $2.00 입니다.\nMomsDonalia: 잠시만 기다려주세요. 주문하신 메뉴를 준비하겠습니다.\n",
		},
		{
			order: DiscountOrder{
				PricePerUnit: 1.0,
				Quantity:     2,
				Discount:     0.1,
			},
			expected: "MomsDonalia: 주문 합계는 $1.80 입니다.\nMomsDonalia: 잠시만 기다려주세요. 주문하신 메뉴를 준비하겠습니다.\n",
		},
		{
			order: GiftCardOrder{
				PricePerUnit: 1.0,
				Quantity:     2,
			},
			expected: "MomsDonalia: 주문 합계는 $0.00 입니다.\nMomsDonalia: 잠시만 기다려주세요. 주문하신 메뉴를 준비하겠습니다.\n",
		},
	}

	for _, test := range tests {
		oldOut := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		m.TakeOrder(test.order)

		w.Close()

		out, _ := io.ReadAll(r)

		r.Close()

		os.Stdout = oldOut

		if string(out) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, string(out))
		}
	}
}
