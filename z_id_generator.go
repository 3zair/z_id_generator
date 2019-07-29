package z_id_generator

import (
	"errors"
	"sync"
	"time"
)

var (
	WorkId int64 = 0          // [0, 31]
	Epoch  int64 = 1561910400 // 起始时间戳 "2019-07-01 00:00:00"
)

const (
	workerIdBits uint8 = 5
	seqBits      uint8 = 16
)

type IdGenerator struct {
	epoch           int64
	seq             int64
	seqMask         int64
	lastTime        int64
	workerIdShifted int64
	timeShift       uint8
	m               sync.Mutex
}

func NewIdGenerator(epoch int64, workId int64) (*IdGenerator, error) {
	if epoch > time.Now().Unix() {
		return nil, errors.New("invalid epoch")
	}
	if workId < 0 || workId > 31 {
		return nil, errors.New("invalid workId")
	}

	return &IdGenerator{
		epoch:           epoch,
		seq:             0,
		seqMask:         2<<seqBits - 1,
		lastTime:        0,
		workerIdShifted: workId << seqBits,
		timeShift:       workerIdBits + seqBits,
	}, nil
}

func (ig *IdGenerator) NewId() int64 {
	ig.m.Lock()
	defer ig.m.Unlock()

	now := time.Now().Unix()

	if now > ig.lastTime {
		ig.seq = 0
	} else {
		ig.seq = (ig.seq + 1) & ig.seqMask

		if ig.seq == 0 {
			for now <= ig.lastTime {
				time.Sleep(time.Millisecond * 100)
				now = time.Now().Unix()
			}
		}
	}

	ig.lastTime = now

	return ((now - ig.epoch) << ig.timeShift) | ig.workerIdShifted | ig.seq
}

var Ig, _ = NewIdGenerator(Epoch, WorkId)

func NewId() int64 {
	return Ig.NewId()
}
