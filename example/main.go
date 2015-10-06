// agentsmith project main.go
package main

import (
	"bytes"
	"log"
	"math/rand"
	"text/template"

	"github.com/avesanen/slacksmith"
)

func main() {
	s := smith.NewSmith()
	s.JackIn("xoxo-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

	bs := NewBotsnack()
	bs.AddTemplate("Thank you, {{.User.Profile.FirstName}}, but I don't like {{.Text}}.")
	bs.AddTemplate("Mr. {{.User.Profile.LastName}}, eat that {{.Text}} yourself.")
	bs.AddTemplate("You hear that, Mr. {{.User.Profile.LastName}}? That is the sound of {{.Text}} in garbage disposal.")
	bs.AddTemplate("You brought me {{.Text}}? You should know Better, Mr. {{.User.Profile.LastName}}, I'm a bot.")
	bs.AddTemplate("Mr. Jenkins would love that {{.Text}}. You should give it to him, Mr. {{.User.Profile.LastName}}.")
	s.Handle("botsnack", bs)

	s.Serve()
}

type BotSnack struct {
	Replies []*template.Template
}

func NewBotsnack() *BotSnack {
	b := &BotSnack{}
	return b
}

func (bs *BotSnack) Process(msg smith.MessageEvent, smith *smith.Smith) {
	t := bs.Replies[rand.Intn(len(bs.Replies))]
	buf := new(bytes.Buffer)
	err := t.Execute(buf, msg)
	if err != nil {
		log.Println("Error rendering template:", err.Error())
		return
	}
	reply := buf.String()
	smith.Rtm.Say(reply, msg.Channel)
	return
}

func (bs *BotSnack) AddTemplate(t string) {
	bs.Replies = append(bs.Replies, template.Must(template.New("").Parse(t)))
}
