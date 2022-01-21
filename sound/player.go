package sound

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
	"os"
)

type Player struct {
	buffers map[string]*beep.Buffer
}

func NewPlayer(sampleRate int, bufferSize int) (*Player, error) {
	err := speaker.Init(beep.SampleRate(sampleRate), bufferSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init speaker")
	}

	return &Player{
		buffers: make(map[string]*beep.Buffer),
	}, nil
}

func (p *Player) Load(files ...string) error {
	for _, v := range files {
		f, err := os.Open(v)
		if err != nil {
			return errors.Wrap(err, "failed to load sound file")
		}
		streamer, format, err := mp3.Decode(f)
		if err != nil {
			return errors.Wrap(err, "failed to decode mp3 file")
		}
		if err != nil {
			return err
		}
		p.buffers[v] = beep.NewBuffer(format)
		p.buffers[v].Append(streamer)
	}
	return nil
}

func (p *Player) Play(file string) {
	speaker.Play(p.buffers[file].Streamer(0, p.buffers[file].Len()))
}
