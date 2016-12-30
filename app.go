package main

import (
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/scgolang/sc"
)

// App maintains all the state for the deepnote app.
type App struct {
	*sc.Client
	*sc.Group

	Config Config

	server    *sc.Server
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
	app.loadTHX5()
	app.loadTHX6()

	// Start scsynth.
	server := &sc.Server{Network: "udp", Port: 57120}
	_, _, err := server.Start(5 * time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "starting server")
	}
	app.server = server

	// Initialize the client.
	client, err := sc.NewClient("udp", config.LocalAddr, config.ScsynthAddr, 5*time.Second)
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
	defer func() { _ = app.server.Stop() }() // Best effort.

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
	return errors.Wrap(app.server.Wait(), "waiting for scsynth")
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
