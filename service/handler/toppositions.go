package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/b4t3ou/cc/service"
)

type TopPositions struct {
	amount int
	data   map[string][]int
}

func NewTopPositions(amount int) *TopPositions {
	return &TopPositions{
		amount: amount,
		data:   map[string][]int{},
	}
}

func (tp *TopPositions) Process(data []byte) error {
	msg := &service.Message{}

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	switch msg.Type {
	case "8":
		log.Printf("recv update: %+v", msg)
	case "9":
		tp.createFromSnapshot(msg.FSYM, msg.TSYM, msg.BID)
		log.Printf("recv snaphost: %+v", msg)
	}

	return nil
}

func (tp *TopPositions) createFromSnapshot(fsym, tsym string, bids []service.Bid) {
	baseKeyP := fmt.Sprintf("%s-%s-P", fsym, tsym)
	baseKeyQ := fmt.Sprintf("%s-%s-Q", fsym, tsym)

	if _, exists := tp.data[baseKeyP]; !exists {
		tp.data[baseKeyP] = []int{}
	}

	if _, exists := tp.data[baseKeyQ]; !exists {
		tp.data[baseKeyQ] = []int{}
	}

	for _, bid := range bids {
		tp.data[baseKeyP] = append(tp.data[baseKeyP], int(bid.P*100))
		tp.data[baseKeyQ] = append(tp.data[baseKeyQ], int(bid.Q*100))
	}

	sort.Ints(tp.data[baseKeyP])
	sort.Ints(tp.data[baseKeyQ])
}
