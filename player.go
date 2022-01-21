package main

import (
	"github.com/pkg/errors"
	"stocks-music/sampler"
)

type NotePlayer struct {
	soundPlayer SoundPlayer
	samples     map[sampler.Note]string
}

func NewNotePlayer(soundPlayer SoundPlayer, smpl map[sampler.Note]string) *NotePlayer {
	return &NotePlayer{soundPlayer: soundPlayer, samples: smpl}
}

func (np *NotePlayer) Init() error {
	for _, src := range np.samples {
		err := np.soundPlayer.Load(src)
		if err != nil {
			return errors.Wrap(err, "failed to load note file")
		}
	}

	return nil
}

func (np *NotePlayer) Play(n sampler.Note) {
	np.soundPlayer.Play(np.samples[n])
}
