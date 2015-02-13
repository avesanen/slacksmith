package smith

import (
	"strings"
)

type Handler interface {
	Process(msg string) string
}

type Smith struct {
	CmdPrefix string
	Handlers  map[string]Handler
	Rtm       *Rtm
}

func NewSmith() *Smith {
	s := &Smith{}
	s.CmdPrefix = "!"
	s.Handlers = make(map[string]Handler)
	return s
}

func (s *Smith) Serve() {
	for {
		select {
		case msg, ok := <-s.Rtm.InChan:
			if !ok {
				return
			}
			reply := s.Parse(msg.Text)
			s.Rtm.Say(reply, msg.Channel)
		}
	}
}

func (s *Smith) JackIn(token string) error {
	url := "https://slack.com/api/rtm.start?token=" + token
	rtm, err := StartRtm(url)
	if err != nil {
		return err
	}
	s.Rtm = rtm
	return nil
}

func (s *Smith) Parse(msg string) string {
	if !strings.HasPrefix(msg, s.CmdPrefix) {
		return ""
	}
	msg = strings.TrimLeft(msg, "!")
	fields := strings.Fields(msg)
	if len(fields) > 0 {
		for i, k := range s.Handlers {
			if i == fields[0] {
				msg = strings.TrimLeft(msg, i)
				msg = strings.TrimLeft(msg, " ")
				return k.Process(msg)
			}
		}
	}
	return ""
}

func (s *Smith) Handle(command string, handler Handler) {
	s.Handlers[command] = handler
}
