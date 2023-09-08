package main

type Config struct {
	Worker    int     `yaml:"worker"`
	Pool      int     `yaml:"pool"`
	StatsDir  string  `yaml:"stats-dir"`
	AvgDelay  float64 `yaml:"avg-delay"`
	Failure   int     `yaml:"failure"`
}
