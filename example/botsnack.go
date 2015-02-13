package main

import (
	"fmt"
	"math/rand"
)

type BotSnack struct {
	Replies []string
}

func (bs *BotSnack) Process(s string) string {
	if bs.Replies == nil {
		bs.Replies = []string{"Thank you, I love %s!"}
	}
	return fmt.Sprintf(bs.Replies[rand.Intn(len(bs.Replies))], s)
}

func (bs *BotSnack) AddReply(s string) {
	bs.Replies = append(bs.Replies, s)
}
