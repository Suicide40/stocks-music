package main

import (
	"github.com/gorilla/websocket"
	"stocks-music/sampler"
	"stocks-music/sound"
	"stocks-music/streamer"
	"time"
)

type SoundPlayer interface {
	Play(file string)
	Load(files ...string) error
}

type Streamer interface {
	Start(ticker string, output chan<- streamer.Tick) error
	Close()
}

func main() {
	soundPlayer, err := sound.NewPlayer(44100, 32)
	if err != nil {
		panic(err.Error())
	}

	notePlayer := NewNotePlayer(soundPlayer, sampler.NewPiano())
	err = notePlayer.Init()
	if err != nil {
		panic(err.Error())
	}

	tickerStreamer := streamer.NewGemini(websocket.DefaultDialer)
	noteStreamer := NewNoteStreamer(tickerStreamer)

	noteChan := make(chan sampler.Note, 100)
	go noteStreamer.Start(noteChan)

	timeout := time.After(time.Second * 20)
	noteCount := 15
	for {
		select {
		case <-timeout:
			noteStreamer.Close()
			return
		case n := <-noteChan:
			if noteCount == 0 {
				noteStreamer.Close()
				return
			}
			notePlayer.Play(n)
			noteCount--
		}
	}

}
