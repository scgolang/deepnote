package main

import (
	"github.com/scgolang/sc"
)

// loadTHX2 loads the synthdef for the second version of the deep note.
func (app *App) loadTHX2() {
	const name = "THX2"

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Mix(sc.AR, app.voicesTHX2()),
		}.Rate(sc.AR)
	})
}

// voicesTHX2 creates the slice of voices for the second version
// of the thx deep note.
// Set each oscillator to a sawtooth wave with a random frequency.
// The big difference from def1 is that
// we add some noisy frequency modulation,
// and higher frequencies will have more modulation.
func (app *App) voicesTHX2() []sc.Input {
	var (
		fundamentals = app.fundamentals()
		voices       = make([]sc.Input, app.Config.NumVoices)
	)
	for i := range voices {
		var (
			amp   = sc.C(1 / float64(app.Config.NumVoices))
			drift = app.voiceDrift(i)
			freq  = sc.C(fundamentals[i]).Add(drift)
		)
		voices[i] = sc.Pan2{
			In: sc.BLowPass{
				In:   sc.Saw{Freq: freq}.Rate(sc.AR),
				Freq: freq.Mul(sc.C(5)),
				RQ:   sc.C(0.5),
			}.Rate(sc.AR),
			Pos:   sc.C(rrand(-0.5, 0.5)),
			Level: amp,
		}.Rate(sc.AR)
	}
	return voices
}

// voiceDrift creates the frequency drift component for the ith voice.
func (app *App) voiceDrift(i int) sc.Input {
	driftAmt := sc.C(3 * (i + 1))
	return sc.LFNoise{
		Interpolation: sc.NoiseQuadratic,
		Freq:          sc.C(0.5),
	}.Rate(sc.AR).Mul(driftAmt)
}
