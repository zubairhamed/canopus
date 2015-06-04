package tests
import (
	"testing"
	"github.com/zubairhamed/canopus"
	"github.com/stretchr/testify/assert"
)

func TestCoreResourceParsing(t *testing.T) {

	cases1 := []struct {
		in 			string
		elemCount 	int
		targets 	[]string
		attrCount 	[]int
	}{
		{
			"</1>,</2>,</3>,</4>,</5/0>,</5/1>,</5/2>,</5/3>",
			8,
			[]string{"/1","/2","/3","/4","/5/0","/5/1","/5/2","/5/3"},
			[]int{0,0,0,0,0,0,0,0,0},

		},
		{
			"</sensors>;ct=40;title=\"Sensor Index\",</sensors/temp>;rt=\"temperature-c\";if=\"sensor\",</sensors/light>;rt=\"light-lux\";if=\"sensor\",<http://www.example.com/sensors/t123>;anchor=\"/sensors/temp\";rel=\"describedby\",</t>;anchor=\"/sensors/temp\";rel=\"alternate\"",
			5,
			[]string{"/sensors","/sensors/temp","/sensors/light", "http://www.example.com/sensors/t123", "/t"},
			[]int{1, 2, 2, 2, 2},
		},
	}

	for _, c := range cases1 {
		resources := canopus.CoreResourcesFromString(c.in)
		assert.Equal(t, len(resources), c.elemCount)

		for i, o := range resources {
			assert.Equal(t, o.Target, c.targets[i])
		}
	}
}
/*
	cases1 := []struct {
		in LWM2MObjectType
	}{
		{oma.OBJECT_LWM2M_SERVER},
		{oma.OBJECT_LWM2M_ACCESS_CONTROL},
		{oma.OBJECT_LWM2M_DEVICE},
		{oma.OBJECT_LWM2M_CONNECTIVITY_MONITORING},
		{oma.OBJECT_LWM2M_FIRMWARE_UPDATE},
		{oma.OBJECT_LWM2M_LOCATION},
		{oma.OBJECT_LWM2M_CONNECTIVITY_STATISTICS},
	}

	for _, c := range cases1 {
		err := cli.EnableObject(c.in, nil)

		assert.Nil(t, err, "Error enabling object: ", c.in)
	}

 */