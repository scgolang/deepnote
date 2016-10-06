package main

import (
	"github.com/scgolang/sc"
)

// loadTHX1 loads the synthdef for the first version of the deep note.
func (app *App) loadTHX1() {
	const name = "THX1"

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Mix(sc.AR, app.voicesTHX1()),
		}.Rate(sc.AR)
	})

}

// voicesTHX1 creates a slice of voices for the first version of the
// thx deep note.
func (app *App) voicesTHX1() []sc.Input {
	voices := make([]sc.Input, app.Config.NumVoices)

	// Set each oscillator to a sawtooth wave with a random frequency.
	for i := range voices {
		var (
			freq = rrand(app.Config.FreqMin, app.Config.FreqMax)
			amp  = sc.C(1 / float64(app.Config.NumVoices))
		)
		voices[i] = sc.Pan2{
			In:    sc.Saw{Freq: sc.C(freq)}.Rate(sc.AR),
			Pos:   sc.C(rrand(-0.5, 0.5)),
			Level: amp,
		}.Rate(sc.AR)
	}
	return voices
}
