package main

import (
	"bytes"
	"encoding/json"
	"machine"
)

type Config struct {
	LowThreshold  [4]float32 `json:"low"`
	HighThreshold [4]float32 `json:"high"`
	State         byte       `json:"state"`
}

func (c Config) Low(i byte) float32  { return c.LowThreshold[i] }
func (c Config) High(i byte) float32 { return c.HighThreshold[i] }
func (c Config) IsOn(i byte) bool {
	mask := byte(1 << i)
	return (c.State & mask) == mask
}

func (c *Config) SetLow(i byte, val float32)  { c.LowThreshold[i] = val }
func (c *Config) SetHigh(i byte, val float32) { c.HighThreshold[i] = val }
func (c *Config) SetOn(i byte, val bool) {
	if val {
		mask := byte(1 << i)
		c.State |= mask
	} else {
		mask := 0xFF ^ byte(1<<i)
		c.State &= mask
	}
}

func saveConfig(config Config) error {
	msg, err := json.Marshal(config)
	if err != nil {
		return err
	}

	needed := int64(len(msg)) / machine.Flash.EraseBlockSize()
	if needed == 0 {
		needed = 1
	}

	err = machine.Flash.EraseBlocks(0, needed)
	if err != nil {
		return err
	}
	_, err = machine.Flash.WriteAt([]byte(msg), 0)

	return err
}

func defaultConfig() Config {
	return Config{
		LowThreshold:  [4]float32{.5, .5, .5, .5},
		HighThreshold: [4]float32{.65, .65, .65, .65},
		State:         0xff,
	}
}

func loadConfig() (Config, error) {
	saved := make([]byte, 1024)
	_, err := machine.Flash.ReadAt(saved, 0)
	if err != nil {
		return defaultConfig(), err
	}

	var cfg Config
	buf := bytes.NewReader(saved)
	decoder := json.NewDecoder(buf)
	err = decoder.Decode(&cfg)
	if err != nil {
		return defaultConfig(), err
	}

	return cfg, nil
}
