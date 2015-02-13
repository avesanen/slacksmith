package main

import (
	"math/rand"
)

type BotQuote struct {
	Replies []string
}

func (bq *BotQuote) Process(s string) string {
	if bq.Replies == nil {
		bq.Replies = []string{"Me, me, me. http://i.imgur.com/e57Qjbb.gif"}
	}
	return bq.Replies[rand.Intn(len(bq.Replies))]
}

func (bq *BotQuote) AddReply(s string) {
	bq.Replies = append(bq.Replies, s)
}
