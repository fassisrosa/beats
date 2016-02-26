package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/fassisrosa/beats/openstackbeat/beater"
)

func main() {
	err := beat.Run("openstackbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
