package handler

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/b4t3ou/cc/service"
)

type TopPositions struct {
	listAmount         int
	midRefreshInterval int
	data               map[string]Item
	mu                 sync.Mutex
}

func NewTopPositions(listAmount, midRefreshInterval int) *TopPositions {
	tp := &TopPositions{
		listAmount:         listAmount,
		midRefreshInterval: midRefreshInterval,
		data:               map[string]Item{},
	}

	go tp.refreshMidPrice()

	return tp
}

type Item struct {
	list BidList
	mid  float32
}

type BidList []float32

func (b BidList) Len() int           { return len(b) }
func (b BidList) Less(i, j int) bool { return b[i] > b[j] }
func (b BidList) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (tp *TopPositions) Process(data []byte) error {
	msg := &service.Message{}

	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}

	switch msg.Type {
	case "8":
		tp.update(msg.FSYM, msg.TSYM, "P", msg.P)
		tp.update(msg.FSYM, msg.TSYM, "Q", msg.Q)
	case "9":
		tp.createFromSnapshot(msg)
	}

	return nil
}

func (tp *TopPositions) createFromSnapshot(msg *service.Message) {
	bidsP := BidList{}
	bidsQ := BidList{}

	for _, bid := range msg.BID {
		bidsP = append(bidsP, bid.P)
		bidsQ = append(bidsQ, bid.Q)
	}

	sort.Sort(bidsP)
	sort.Sort(bidsQ)

	tp.data[fmt.Sprintf("%s-%s-P", msg.FSYM, msg.TSYM)] = Item{list: bidsP[:tp.listAmount]}
	tp.data[fmt.Sprintf("%s-%s-Q", msg.FSYM, msg.TSYM)] = Item{list: bidsQ[:tp.listAmount]}
}

func (tp *TopPositions) update(fsym, tsym, typeName string, value float32) {
	key := fmt.Sprintf("%s-%s-%s", fsym, tsym, typeName)
	item := tp.data[key]

	item.list = append(item.list, value)

	sort.Sort(item.list)
	item.list = item.list[:tp.listAmount]

	tp.data[key] = item
}

func (tp *TopPositions) refreshMidPrice() {
	for range time.Tick(1 * time.Second) {
		if time.Now().Second()%tp.midRefreshInterval != 0 {
			continue
		}

		tp.mu.Lock()
		for key, value := range tp.data {
			tp.calculateMid(key, value)
		}
		tp.mu.Unlock()
	}
}

func (tp *TopPositions) calculateMid(key string, item Item) {
	var (
		sum float32
	)

	for _, value := range item.list {
		sum += value
	}

	item.mid = sum / float32(tp.listAmount)
	tp.data[key] = item

	fmt.Printf("%s mid calculated for %s: %f\n", time.Now().Format("2006-01-02T15:04:05"), key, item.mid)
}
