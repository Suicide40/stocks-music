package main

import (
	"fmt"
	"stocks-music/sampler"
	"stocks-music/streamer"
	"strconv"
	"time"
)

type NoteStreamer struct {
	streamer Streamer
}

func NewNoteStreamer(streamer Streamer) *NoteStreamer {
	return &NoteStreamer{streamer: streamer}
}

func (s NoteStreamer) Close() {
	s.streamer.Close()
	fmt.Println("streamer closed")
}

func (s NoteStreamer) Start(output chan sampler.Note) {
	ticks := make(chan streamer.Tick, 100)
	go s.streamer.Start("BTCUSD", ticks)

	for {
		pause := time.After(time.Second)
		select {
		case tick := <-ticks:
			if s.canSkip(tick) {
				break
			}
			output <- s.convertToNote(tick)
		}
		<-pause
	}
}

func (s NoteStreamer) canSkip(tick streamer.Tick) bool {
	return len(tick.Events) > 5
}

func (s NoteStreamer) convertToNote(tick streamer.Tick) sampler.Note {
	if len(tick.Events) == 0 {
		return sampler.A
	}

	price, _ := strconv.ParseFloat(tick.Events[0].Price, 64)

	mapper := map[int]sampler.Note{
		1:  sampler.A,
		2:  sampler.ASharp,
		3:  sampler.B,
		4:  sampler.C,
		5:  sampler.CSharp,
		6:  sampler.D,
		7:  sampler.DSharp,
		8:  sampler.E,
		9:  sampler.F,
		10: sampler.FSharp,
		11: sampler.G,
		12: sampler.GSharp,
	}

	ix := int(price) % 13

	return mapper[ix]
}
