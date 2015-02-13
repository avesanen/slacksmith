package smith

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Rtm struct {
	Ok   bool   `json:"ok"`
	Url  string `json:"url"`
	Self struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
	Team     json.RawMessage `json:"team"`
	Users    json.RawMessage `json:"users"`
	Channels json.RawMessage `json:"channels"`
	Groups   json.RawMessage `json:"groups"`
	Ims      json.RawMessage `json:"ims"`
	Bots     json.RawMessage `json:"bots"`
	Error    string          `json:"error,omitempty"`
	Conn     *websocket.Conn `json:"-"`
	OutChan  chan string     `json:"-"`
	InChan   chan MessageMsg `json:"-"`
}

func (rtm *Rtm) Reader() {
	go func() {
		defer close(rtm.InChan)
		for {
			var b = make([]byte, 2048)
			var n int
			var err error
			n, err = rtm.Conn.Read(b)
			if err != nil {
				log.Println(err.Error())
				return
			}
			var msg Msg
			err = json.Unmarshal(b[:n], &msg)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if msg.Type == "pong" {
				continue
			} else if msg.Type == "message" {
				var message MessageMsg
				err = json.Unmarshal(b[:n], &message)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				rtm.InChan <- message
			}
		}
	}()
}

func (rtm *Rtm) Writer() {
	go func() {
		for {
			select {
			case msg, ok := <-rtm.OutChan:
				if !ok {
					return
				}
				_, err := rtm.Conn.Write([]byte(msg))
				if err != nil {
					log.Println(err.Error())
					continue
				}
			}
		}
	}()
	go func() {
		for {
			<-time.After(time.Second * 1)
			rtm.OutChan <- "{\"id\":1234,\"type\":\"ping\"}"
		}
	}()
}

func (rtm *Rtm) Dial() error {
	u, err := url.Parse(rtm.Url)
	if err != nil {
		return err
	}
	u.Host += ":443"

	ws, err := websocket.Dial(u.String(), "", "http://localhost/")
	if err != nil {
		return err
	}

	rtm.Conn = ws
	return nil
}

func (rtm *Rtm) Say(s string, c string) error {
	m := &ChatMsg{
		Id:      1,
		Type:    "message",
		Channel: c,
		Text:    s,
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	rtm.OutChan <- string(b)
	return nil
}

func StartRtm(rtmUrl string) (*Rtm, error) {
	resp, err := http.Get(rtmUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rtm Rtm
	if err := json.Unmarshal(body, &rtm); err != nil {
		return nil, err
	}

	if !rtm.Ok {
		return nil, errors.New(rtm.Error)
	}

	err = rtm.Dial()
	if err != nil {
		return nil, err
	}

	rtm.OutChan = make(chan string)
	rtm.InChan = make(chan MessageMsg)
	rtm.Reader()
	rtm.Writer()
	return &rtm, nil
}
