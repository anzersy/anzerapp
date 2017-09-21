package main

import (
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/pkg/errors"
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v2"

)

var VERSION = "v0.0.0-dev"


type Func func(str string)

func main() {
	app := cli.NewApp()
	app.Name = "anzerapp"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Action = run
	app.Run(os.Args)
}

func run(c *cli.Context) error {
	//identify the error channel
	exit := make(chan error)

	//use goroutine to kick off the communication with the cattle by event stream.
	go func(exit chan<- error) {
		err := TestEventStream("http://120.77.214.113:8080/v2-beta/projects/1a5","7127B26276D2D9B6518B","igc99qWdRcRRj3YCepgRYmj7fKts2tE8eZMqCLQZ")
		exit <- errors.Wrapf(err, "Closed the event stream.")
	}(exit)

	err := <-exit
	return err
}

//Identify the handler which watched on resource.change event.
func HandleResourceChange(event *events.Event, client *client.RancherClient) error {
	logrus.Info("Start to handle the resourcechange event.")
	logrus.Info(event)
	return nil
}



//Identify the process which contained the core logic of communication with cattle.
func TestEventStream(cattleURL, accessKey, secretKey string) error {
	logrus.Info("Start to connect the Cattle.")

	//initialize the event handler mapp
	eventhandlermap := map[string]events.EventHandler{
		"resource.change":	HandleResourceChange,
		"ping":	func(e *events.Event, c *client.RancherClient) error {logrus.Info("Here we ping.");return  nil},
	}


	//Use event handle map and subscribe to cattle, and then establish the event stream to process the event.
	router, err :=events.NewEventRouter("", 0, cattleURL, accessKey, secretKey, nil, eventhandlermap, "", 100, events.DefaultPingConfig)

	if err != nil {
		return err
	}

	err = router.StartWithoutCreate(nil)
	return err
}