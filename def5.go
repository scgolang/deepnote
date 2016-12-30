package main

import (
	"github.com/scgolang/sc"
)

// loadTHX5 loads the synthdef for the fourth version of the deep note.
func (app *App) loadTHX5() {
	const name = "THX5"

	outerEnv := app.outerEnv5()
	sig := sc.BLowPass{
		In:   sc.Mix(sc.AR, app.voicesTHX5()),
		Freq: outerEnv.MulAdd(sc.C(18000), sc.C(2000)),
		RQ:   sc.C(0.5),
	}.Rate(sc.AR)

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Multi(sig, sig),
		}.Rate(sc.AR)
	})
}

// voicesTHX5 creates the array of voices for the fifth version of the thx deep note.
func (app *App) voicesTHX5() []sc.Input {
	var (
		fundamentals = app.fundamentals()
		finalPitches = app.finalPitches()
		voices       = make([]sc.Input, app.Config.NumVoices)
		sweepEnv     = app.sweepEnv5()
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

// sweepEnv5 returns an envelope used to modulate the frequency
// of the ith voice in version 5 of the thx deep note.
func (app *App) sweepEnv5() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(rrand(0.1, 0.2)),
				sc.C(1),
			},
			Times: []sc.Input{
				sc.C(rrand(5, 6)),
				sc.C(rrand(8, 9)),
			},
			Curve: []float64{
				rrand(2, 3),
				rrand(4, 5),
			},
		},
	}.Rate(sc.KR)
}

// outerEnv5 creates an envelope used to modulate the cutoff
// frequency of the outermost lowpass filter in version 5 of
// the thx deep note.
func (app *App) outerEnv5() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(0.1),
				sc.C(1),
			},
			Times: []sc.Input{
				sc.C(8),
				sc.C(4),
			},
			Curve: []interface{}{2, 4}, // exp, welch
		},
	}.Rate(sc.KR)
}
