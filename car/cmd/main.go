package main

import (
	"flag"

	"github.com/freeconf/examples/car"

	"github.com/freeconf/gconf/c2"
	"github.com/freeconf/gconf/device"
	"github.com/freeconf/gconf/meta"
	"github.com/freeconf/gconf/restconf"
)

// Initialize and start our Car micro-service application with C2Stack for
// RESTful based management
//
// To run:
//    export GOPATH=`pwd`
//    cd ./src/vendor/github.com/freeconf/gconf/examples/car/cmd
//    go run ./main.go -startup startup.json
//
// Then open web browser to
//   http://localhost:8080
//
var startup = flag.String("startup", "startup.json", "start-up configuration file.")
var verbose = flag.Bool("verbose", false, "verbose")

func main() {
	flag.Parse()
	c2.DebugLog(*verbose)

	// Any existing application
	app := car.New()

	// Where to looks for yang files, this tells library to use these
	// two relative paths.  StreamSource is an abstraction to data sources
	// that might be local or remote or combinations of all the above.
	yangPath := meta.PathStreamSource("..:../../../gconf/yang")

	// Every management has a "device" container. A device can have many "modules"
	// installed which are really microservices.
	//   carPath - where UI files are located
	//   ypath - where *.yang files are located
	d := device.New(yangPath)

	// Here we are installing the "car" module which is our main application.
	//   "car" - the name of the module causing car.yang to load from yang path
	//   car.Node(app) - we are linking our car application with driver to handle
	//           car management requests.
	chkErr(d.Add("car", car.Manage(app)))

	// Adding RESTCONF protocol support.  Should you want an alternate protocol,
	// you could do that here. restconf will add whatever modules it needs to your
	// device object.
	restconf.NewServer(d)

	// Even though the main configuration comes from the application management
	// system after call-home has registered this system it's often useful/neccessary
	// to bootstrap config for some of the local modules
	chkErr(d.ApplyStartupConfigFile(*startup))

	// wait for cntrl-c...
	select {}
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}
