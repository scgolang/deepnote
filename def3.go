package main

import (
	"math"

	"github.com/scgolang/sc"
)

// loadTHX3 loads the synthdef for the third version of the deep note.
func (app *App) loadTHX3() {
	const name = "THX3"

	app.synthdefs[name] = sc.NewSynthdef(name, func(params sc.Params) sc.Ugen {
		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Mix(sc.AR, app.voicesTHX3()),
		}.Rate(sc.AR)
	})
}

// voicesTHX3 creates the array of voices for the third version of the thx deep note.
// Set each voice to a sawtooth wave with a random frequency.
// The the big difference from def2 is that we sweep the frequencies so that
// they all converge to the note between D and Eb (MIDI note 14.5).
func (app *App) voicesTHX3() []sc.Input {
	var (
		fundamentals = app.fundamentals()
		finalPitches = app.finalPitches()
		voices       = make([]sc.Input, app.Config.NumVoices)
		sweepEnv     = app.sweepEnv3()
		invSweep     = sweepEnv.MulAdd(sc.C(-1), sc.C(1))
	)
	for i := range voices {
		var (
			drift       = app.voiceDrift(i)
			finalFreq   = sc.C(finalPitches[i])
			initialFreq = sc.C(fundamentals[i]).Add(drift)
			freq        = initialFreq.Mul(invSweep).Add(finalFreq.Mul(sweepEnv))
			amp         = sc.C(1 / float64(app.Config.NumVoices))
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

// finalPitches returns a slice of final pitches for the
// sweeping versions of the thx deep note.
// The final pitches consist of MIDI note 14.5 spanning 6 octaves.
func (app *App) finalPitches() []float32 {
	pitches := make([]float32, app.Config.NumVoices)
	for i := range pitches {
		octave := math.Trunc(float64(i) / float64(app.Config.NumVoices/6))
		pitches[i] = sc.Midicps((12 * float32(octave)) + 14.5)
	}
	return pitches
}

// sweepEnv3 returns the envelope used to sweep the pitches.
func (app *App) sweepEnv3() sc.Input {
	return sc.EnvGen{
		Env: sc.Env{
			Levels: []sc.Input{
				sc.C(0),
				sc.C(1),
			},
			Times: []sc.Input{
				sc.C(13),
			},
		},
	}.Rate(sc.KR)
}
