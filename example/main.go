package main

import (
	"github.com/avesanen/slacksmith"
)

var token string

func main() {
	token = ""

	s := smith.NewSmith()
	s.JackIn(token)

	bs := &BotSnack{}
	bs.AddReply("I don't like %s, it's the smell.")
	bs.AddReply("I'm a bot, I don't need to eat %s.")
	bs.AddReply("No thank you, I think %s smells revolting.")
	bs.AddReply("You can eat eat that %s yourself.")
	bs.AddReply("Why would I eat %s? Are you serious?")
	bs.AddReply("Go feed that %s to Mr.Jenkins, he loves stuff like that.")
	bs.AddReply("I think Mr.Jenkins would love %s.")
	bs.AddReply("You can eat it yourself, %s is against my programming.")

	bq := &BotQuote{}
	bq.AddReply("http://i.imgur.com/e57Qjbb.gif")
	bq.AddReply("http://i.imgur.com/7OcXln9.gif")
	bq.AddReply("http://i.imgur.com/1egZEvE.gif")
	bq.AddReply("http://i.imgur.com/C5RMJCO.gif")
	bq.AddReply("http://i.imgur.com/ahIDqwG.gif")

	bq.AddReply("I'd like to share a revelation that I've had during my time here. It came to me when I tried to classify your species and I realized that you're not actually mammals. Every mammal on this planet instinctively develops a natural equilibrium with the surrounding environment but you humans do not. You move to an area and you multiply and multiply until every natural resource is consumed and the only way you can survive is to spread to another area. There is another organism on this planet that follows the same pattern. Do you know what it is? A virus. Human beings are a disease, a cancer of this planet. You're a plague and we are the cure.")

	s.Handle("botsnack", bs)
	s.Handle("quote", bq)

	s.Serve()
}
