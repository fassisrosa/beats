package beater

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/fassisrosa/beats/openstackbeat/config"
	"github.com/fassisrosa/beats/openstackbeat/openstackapi"
)

type Openstackbeat struct {
	Configuration *config.Config
	done          chan struct{}
	period        time.Duration
	openStackAPI  openstackapi.OpenStackAPI
	openStackUrl  string
}

// Creates beater
func New() *Openstackbeat {
	return &Openstackbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Openstackbeat) Config(b *beat.Beat) error {

	// Load beater configuration
	err := cfgfile.Read(&bt.Configuration, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Openstackbeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.Configuration.Openstackbeat.Period == "" {
		bt.Configuration.Openstackbeat.Period = "1s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.Configuration.Openstackbeat.Period)
	if err != nil {
		return err
	}

	bt.openStackAPI, err = openstackapi.NewOpenStackAPI(bt.Configuration.Openstackbeat.OpenStackVersion)
	if err != nil {
		return err
	}
	bt.openStackUrl = bt.Configuration.Openstackbeat.OpenStackUrl

	return nil
}

func (bt *Openstackbeat) Run(b *beat.Beat) error {
	logp.Info("openstackbeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		bt.captureCurrentOpenStackVMInfo(b)
	}
}

func (bt *Openstackbeat) captureCurrentOpenStackVMInfo (b *beat.Beat) {
	virtualMachineInfo, err := bt.openStackAPI.GetAllInfo(bt.openStackUrl)
	if err != nil {
		fmt.Printf("Error running: %s\n", err)
		os.Exit(1)
	}
	for _, oneVM := range virtualMachineInfo {
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
		}
		// collect structure info
		vmValue := reflect.ValueOf(oneVM)
		vmType := vmValue.Type()
    		for i := 0; i < vmValue.NumField(); i++ {
			valueInfo := vmType.Field(i)
        		event[valueInfo.Name] = vmValue.FieldByName(valueInfo.Name).Interface()
    		}
		logp.Debug("openstackbeat publishing event %s", event.String())
		b.Events.PublishEvent(event)
	}
	logp.Info("openstackbeat published %d events", len(virtualMachineInfo))
}

func (bt *Openstackbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Openstackbeat) Stop() {
	close(bt.done)
}
