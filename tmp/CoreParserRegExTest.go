package main
import (
    "regexp"
    "strings"
    "github.com/zubairhamed/goap"
)

func main() {
    str := `</test>;rt="test";ct=0,</validate>;rt="validate";ct=0,</hello>;rt="Type1";ct=0;if="If1",</bl%C3%A5b%C3%A6rsyltet%C3%B8y>;rt="blåbærsyltetøy";ct=0,</sink>;rt="sink";ct=0,</separate>;rt="separate";ct=0,</large>;rt="Type1 Type2";ct=0;sz=1700;if="If2",</secret>;rt="secret";ct=0,</broken>;rt="Type2 Type1";ct=0;if="If2 If1",</weird33>;rt="weird33";ct=0,</weird44>;rt="weird44";ct=0,</weird55>;rt="weird55";ct=0,</weird333>;rt="weird333";ct=0,</weird3333>;rt="weird3333";ct=0,</weird33333>;rt="weird33333";ct=0,</123412341234123412341234>;rt="123412341234123412341234";ct=0,</location-query>;rt="location-query";ct=0,</create1>;rt="create1";ct=0,</large-update>;rt="large-update";ct=0,</large-create>;rt="large-create";ct=0,</query>;rt="query";ct=0,</seg1>;rt="seg1";ct=40,</path>;rt="path";ct=40,</location1>;rt="location1";ct=40,</multi-format>;rt="multi-format";ct=0,</3>;rt="3";ct=50,</4>;rt="4";ct=50,</5>;rt="5";ct=50`
    var re = regexp.MustCompile(`(<[^>]+>\s*(;\s*\w+\s*(=\s*(\w+|"([^"\\]*(\\.[^"\\]*)*)")\s*)?)*)`)

    var resources []*goap.CoreResource
    m := re.FindAllString(str, -1)
    for _, match := range m {
        var elemRe = regexp.MustCompile(`<\/[a-zA-Z0-9_%-]+>`)
        elemMatch := elemRe.FindString(match)
        target :=elemMatch[1:len(elemMatch)-1]

        resource := goap.NewCoreResource()
        resource.Target = target

        attrs := strings.Split(match[len(elemMatch)+1:], ";")

        for _, attr := range attrs {
            pair := strings.Split(attr, "=")

            resource.AddAttribute(pair[0], strings.Replace(pair[1], "\"", "", -1))
        }

        resources = append(resources, resource)
    }
}
