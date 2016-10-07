package main

import (
	"flag"
	"os"
)

// Config contains all the configuration parameters for the deepnote app.
type Config struct {
	LocalAddr   string  `json:"local_addr"`
	ScsynthAddr string  `json:"scsynth_addr"`
	Synthdef    string  `json:"synthdef"`
	NumVoices   int     `json:"num_voices"`
	FreqMin     float64 `json:"freq_min"`
	FreqMax     float64 `json:"freq_max"`
}

func NewConfig() (Config, error) {
	var (
		config = Config{}
		fs     = flag.NewFlagSet("deepnote", flag.ExitOnError)
	)
	fs.StringVar(&config.LocalAddr, "local-addr", "127.0.0.1:0", "local UDP listening address")
	fs.StringVar(&config.ScsynthAddr, "scsynth-addr", "127.0.0.1:57120", "scsynth UDP listening address")
	fs.StringVar(&config.Synthdef, "synthdef", DefaultSynthdef, "synthdef (format is THX<VERSION> where VERSION is one of the versions of the thx deep note from http://www.earslap.com/article/recreating-the-thx-deep-note.html)")
	fs.IntVar(&config.NumVoices, "num-voices", 30, "number of voices")
	fs.Float64Var(&config.FreqMin, "freq-min", 200, "lower bound used when generating random fundamental frequencies")
	fs.Float64Var(&config.FreqMax, "freq-max", 400, "upper bound used when generating random fundamental frequencies")
	return config, fs.Parse(os.Args[1:])
}
