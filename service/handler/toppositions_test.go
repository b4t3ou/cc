package handler

import (
	"encoding/json"
	"github.com/b4t3ou/cc/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	testSnapshot = service.Message{
		Type:       "9",
		FSYM:       "ERT",
		TSYM:       "EUR",
		P:          0,
		Q:          0,
		SEQ:        0,
		REPORTEDNS: 0,
		DELAYNS:    0,
		BID: []service.Bid{
			{P: 1.1, Q: 1.1},
			{P: 1.2, Q: 1.2},
			{P: 1.3, Q: 1.3},
			{P: 1.4, Q: 1.4},
			{P: 1.5, Q: 1.5},
			{P: 1.6, Q: 1.6},
			{P: 1.7, Q: 1.7},
			{P: 1.8, Q: 1.8},
			{P: 1.9, Q: 1.9},
			{P: 2.0, Q: 2.0},
			{P: 2.1, Q: 2.1},
			{P: 2.2, Q: 2.2},
			{P: 2.3, Q: 2.3},
		},
	}
)

func TestTopPositions_Process(t *testing.T) {
	amount := 10
	h := NewTopPositions(amount, 1)
	snapshot, _ := json.Marshal(testSnapshot)

	err := h.Process(snapshot)
	assert.Nil(t, err)

	assert.Equal(t, amount, len(h.data["ERT-EUR-P"].list))
	assert.Equal(t, amount, len(h.data["ERT-EUR-Q"].list))
	assert.Equal(t, float32(2.3), h.data["ERT-EUR-P"].list[0])
	assert.Equal(t, float32(1.4), h.data["ERT-EUR-P"].list[9])
	assert.Equal(t, float32(2.3), h.data["ERT-EUR-Q"].list[0])
	assert.Equal(t, float32(1.4), h.data["ERT-EUR-Q"].list[9])

	time.Sleep(time.Second * 2)
	assert.Equal(t, float32(1.85), h.data["ERT-EUR-P"].mid)

	data, _ := json.Marshal(service.Message{
		Type: "8",
		FSYM: "ERT",
		TSYM: "EUR",
		P:    10.41,
		Q:    10.41,
	})

	err = h.Process(data)
	assert.Nil(t, err)

	assert.Equal(t, amount, len(h.data["ERT-EUR-P"].list))
	assert.Equal(t, amount, len(h.data["ERT-EUR-Q"].list))
	assert.Equal(t, float32(10.41), h.data["ERT-EUR-P"].list[0])
	assert.Equal(t, float32(1.5), h.data["ERT-EUR-P"].list[9])
	assert.Equal(t, float32(10.41), h.data["ERT-EUR-Q"].list[0])
	assert.Equal(t, float32(1.5), h.data["ERT-EUR-Q"].list[9])

	time.Sleep(time.Second * 2)
	assert.Equal(t, float32(2.751), h.data["ERT-EUR-P"].mid)
}
