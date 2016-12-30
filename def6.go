package main

import (
	"sort"

	"github.com/scgolang/sc"
)

// loadTHX6 loads the synthdef for the fourth version of the deep note.
func (app *App) loadTHX6() {
	const name = "THX6"

	var (
		outerEnv = app.outerEnv5()
		ampEnv   = app.ampEnv6()
		amp      = outerEnv.Add(sc.C(2)).Mul(ampEnv)
	)
	sig := sc.Limiter{
		In: sc.BLowPass{
			In:   sc.Mix(sc.AR, app.voicesTHX6()),
			Freq: outerEnv.MulAdd(sc.C(18000), sc.C(2000)),
			RQ:   sc.C(0.5),
		}.Rate(sc.AR).Mul(amp),
	}.Rate(sc.AR)

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Multi(sig, sig),
		}.Rate(sc.AR)
	})
}

// voicesTHX6 creates the array of voices for the fifth version of the thx deep note.
func (app *App) voicesTHX6() []sc.Input {
	var (
		fundamentals = app.fundamentals()
		finalPitches = app.finalPitches()
		voices       = make([]sc.Input, app.Config.NumVoices)
		sweepEnv     = app.sweepEnv6()
		invSweep     = sweepEnv.MulAdd(sc.C(-1), sc.C(1))
	)

	// Reverse the order of the fundamentals.
	sort.Sort(sort.Reverse(sort.Float64Slice(fundamentals)))

	for i := range voices {
		var (
			initialDrift = app.voiceDrift6(i)
			initialFreq  = sc.C(fundamentals[i]).Add(initialDrift)
			finalDrift   = app.finalDrift6(i)
			finalFreq    = sc.C(finalPitches[i]).Add(finalDrift)
			freq         = initialFreq.Mul(invSweep).Add(finalFreq.Mul(sweepEnv))
			amp          = sc.C((1 - (1 / float64(i+1))) + 1.5)
		)
		voices[i] = sc.Pan2{
			In: sc.BLowPass{
				In:   sc.Saw{Freq: freq}.Rate(sc.AR),
				Freq: freq.Mul(sc.C(6)),
				RQ:   sc.C(0.6),
			}.Rate(sc.AR),
			Pos:   sc.C(rrand(-0.5, 0.5)),
			Level: amp,
		}.Rate(sc.AR)
	}
	return voices
}

// sweepEnv6 returns an envelope used to modulate the frequency
// of the ith voice in version 6 of the thx deep note.
func (app *App) sweepEnv6() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(rrand(0.1, 0.2)),
				sc.C(1),
			},
			Times: []sc.Input{
				sc.C(rrand(5.5, 6)),
				sc.C(rrand(8.5, 9)),
			},
			Curve: []float64{
				rrand(2, 3),
				rrand(4, 5),
			},
		},
	}.Rate(sc.KR)
}

// ampEnv6 returns the amp envelope for the final version
// of the thx deep note.
func (app *App) ampEnv6() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(1),
				sc.C(1),
				sc.C(0),
			},
			Times: []sc.Input{
				sc.C(3),
				sc.C(21),
				sc.C(3),
			},
			Curve: []int{2, 0, -4},
		},
		Done: sc.FreeEnclosing,
	}.Rate(sc.AR)
}

// voiceDrift6 creates the frequency drift component for the ith voice
// in the final version of the thx deep note.
func (app *App) voiceDrift6(i int) sc.Input {
	driftAmt := sc.C(6 * (app.Config.NumVoices - (i + 1)))
	return sc.LFNoise{
		Interpolation: sc.NoiseQuadratic,
		Freq:          sc.C(0.5),
	}.Rate(sc.AR).Mul(driftAmt)
}

// finalDrift6 returns a noise source used to wobble the pitch
// at the end of the frequency envelope sweep for the final
// version of the thx deep note.
func (app *App) finalDrift6(i int) sc.Input {
	return sc.LFNoise{
		Interpolation: sc.NoiseQuadratic,
		Freq:          sc.C(0.1),
	}.Rate(sc.KR).Mul(sc.C(float32(i) / 3))
}
