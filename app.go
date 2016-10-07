package main

import (
	"os"
	"sort"

	"github.com/scgolang/sc"
)

// App maintains all the state for the deepnote app.
type App struct {
	*sc.Client
	*sc.Group

	Config Config

	synthdefs map[string]*sc.Synthdef
}

// NewApp creates a new deepnote app.
func NewApp(config Config) (*App, error) {
	// Create the app and the synthdefs.
	app := &App{
		Config:    config,
		synthdefs: map[string]*sc.Synthdef{},
	}
	app.loadTHX1()
	app.loadTHX2()
	app.loadTHX3()
	app.loadTHX4()

	// Initialize the client.
	client, err := sc.NewClient("udp", config.LocalAddr, config.ScsynthAddr)
	if err != nil {
		return nil, err
	}
	app.Client = client

	// Create the default group.
	group, err := client.AddDefaultGroup()
	if err != nil {
		return nil, err
	}
	app.Group = group

	return app, nil
}

// Run runs the deepnote app.
func (app *App) Run() error {
	// Send the synthdef.
	def := app.Config.Synthdef
	if err := app.SendDef(app.synthdefs[def]); err != nil {
		return err
	}

	// Write the synthdef to a file.
	defFile, err := os.Create(def + ".gosyndef")
	if err != nil {
		return err
	}
	if err := app.synthdefs[def].Write(defFile); err != nil {
		return err
	}

	// Create a synth node.
	var (
		sid    = app.NextSynthID()
		action = sc.AddToTail
		ctls   = map[string]float32{}
	)
	if _, err := app.Group.Synth(def, sid, action, ctls); err != nil {
		return err
	}
	return nil
}

// fundamentals returns the fundamental frequencies.
// This generates a sorted list of frequencies between freq-min and freq-max.
func (app *App) fundamentals() []float64 {
	frequencies := make([]float64, app.Config.NumVoices)
	for i := range frequencies {
		frequencies[i] = rrand(app.Config.FreqMin, app.Config.FreqMax)
	}
	sort.Float64s(frequencies)
	return frequencies
}
