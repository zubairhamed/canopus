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
		attributes  []map[string]string
	}{
		{
			"</1>,</2>,</3>,</4>,</5/0>,</5/1>,</5/2>,</5/3>",
			8,
			[]string{"/1","/2","/3","/4","/5/0","/5/1","/5/2","/5/3"},
			[]int{0,0,0,0,0,0,0,0,0},
			nil,
		},
		{
			"</sensors>;ct=40;title=\"Sensor Index\",</sensors/temp>;rt=\"temperature-c\";if=\"sensor\",</sensors/light>;rt=\"light-lux\";if=\"sensor\",<http://www.example.com/sensors/t123>;anchor=\"/sensors/temp\";rel=\"describedby\",</t>;anchor=\"/sensors/temp\";rel=\"alternate\"",
			5,
			[]string{"/sensors","/sensors/temp","/sensors/light", "http://www.example.com/sensors/t123", "/t"},
			[]int{1, 2, 2, 2, 2},
			[]map[string]string {
				map[string]string {
					"ct": "40",
					"title": "Sensor Index",
				},
				map[string]string {
					"rt": "temperature-c",
					"if": "sensor",
				},
				map[string]string {
					"rt": "light-lux",
					"if": "sensor",
				},
				map[string]string {
					"anchor": "/sensors/temp",
					"rel": "describedby",
				},
				map[string]string {
					"anchor": "/sensors/temp",
					"rel": "alternate",
				},
			},
		},
	}

	for _, c := range cases1 {
		resources := canopus.CoreResourcesFromString(c.in)
		assert.Equal(t, len(resources), c.elemCount)

		for i, o := range resources {
			assert.Equal(t, o.Target, c.targets[i])

			for _, a := range o.Attributes {
				key := a.Key
				val := a.Value

				assert.Equal(t, c.attributes[i][key], val)
			}
		}
	}
}
