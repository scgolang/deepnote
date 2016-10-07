package main

import (
	"github.com/scgolang/sc"
)

// loadTHX4 loads the synthdef for the fourth version of the deep note.
func (app *App) loadTHX4() {
	const name = "THX4"

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Mix(sc.AR, app.voicesTHX4()),
		}.Rate(sc.AR)
	})
}

// voicesTHX4 creates the array of voices for the third version of the thx deep note.
// Set each voice to a sawtooth wave with a random frequency.
// The the big difference from def2 is that we sweep the frequencies so that
// they all converge to the note between D and Eb (MIDI note 14.5).
func (app *App) voicesTHX4() []sc.Input {
	var (
		fundamentals = app.fundamentals()
		finalPitches = app.finalPitches()
		voices       = make([]sc.Input, app.Config.NumVoices)
		sweepEnv     = app.sweepEnv4()
		invSweep     = sweepEnv.MulAdd(sc.C(-1), sc.C(1))
	)
	for i := range voices {
		var (
			initialDrift = app.voiceDrift(i)
			initialFreq  = sc.C(fundamentals[i]).Add(initialDrift)
			finalDrift   = app.finalDrift(i)
			finalFreq    = sc.C(finalPitches[i]).Add(finalDrift)
			freq         = initialFreq.Mul(invSweep).Add(finalFreq.Mul(sweepEnv))
			amp          = sc.C(1 / float64(app.Config.NumVoices))
		)
		voices[i] = sc.Pan2{
			In: sc.BLowPass{
				In:   sc.Saw{Freq: freq}.Rate(sc.AR),
				Freq: freq.Mul(sc.C(8)),
				RQ:   sc.C(0.5),
			}.Rate(sc.AR),
			Pos:   sc.C(rrand(-0.5, 0.5)),
			Level: amp,
		}.Rate(sc.AR)
	}
	return voices
}

// finalDrift returns a noise source used to wobble the pitch
// at the end of the frequency envelope sweep.
func (app *App) finalDrift(i int) sc.Input {
	return sc.LFNoise{
		Interpolation: sc.NoiseQuadratic,
		Freq:          sc.C(0.1),
	}.Rate(sc.KR).Mul(sc.C(float32(i) / 4))
}

// sweepEnv4 returns the envelope used to sweep the pitches.
func (app *App) sweepEnv4() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(0.1),
				sc.C(1),
			},
			Times: []sc.Input{
				sc.C(5),
				sc.C(8),
			},
			Curve: []interface{}{2, 5},
		},
	}.Rate(sc.KR)
}
