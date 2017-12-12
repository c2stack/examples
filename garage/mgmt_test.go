package garage

import (
	"testing"

	"github.com/freeconf/examples/car"

	"github.com/freeconf/c2g/device"
	"github.com/freeconf/c2g/meta"
)

var devices = make(map[string]device.Device)

func TestManagement(t *testing.T) {
	ypath := meta.PathStreamSource("../vendor/github.com/freeconf/c2g/yang:../car:.")

	dm := device.NewMap()

	car0 := car.New()
	dev0 := device.New(ypath)
	chkErr(dev0.Add("car", car.Manage(car0)))
	dm.Add("dev0", dev0)
	car0.Speed = 10
	car0.Start()

	g := NewGarage()
	dev1 := device.New(ypath)
	chkErr(dev1.Add("garage", Manage(g)))
	dm.Add("dev1", dev1)
	o := g.Options()
	o.TireRotateMiles = 100
	o.PollTimeMs = 100
	g.ApplyOptions(o)

	ManageCars(g, dm)

	wait := make(chan struct{})
	g.OnUpdate(func(c Car, w WorkType) {
		t.Logf("car %s did %s", c.Id(), w)
		wait <- struct{}{}
	})
	t.Log("waiting for maintenance...")
	<-wait
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}
