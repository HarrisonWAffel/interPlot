package Networking

import (
	"fmt"
	"github.com/fiorix/freegeoip"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"image/color"
)

//plotPoints uses static maps to create an output image, it takes  an array of Ip structs as its only parameter
func plotPoints(ips []freegeoip.DefaultQuery) {

	ctx := sm.NewContext()
	ctx.SetSize(4000, 3000)

	for i, _ := range ips {
		ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(ips[i].Location.Latitude, ips[i].Location.Longitude), color.RGBA{0xff, 0, 0, 0xff}, 16.0))
		fmt.Printf("\r%s  %d %s %d %s", "Plotted", i, "out of ", len(ips), "IPS")
	}

	fmt.Println("Rendering Image...")
	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	fmt.Println("Saving Image...")
	if err := gg.SavePNG("templates/output.png", img); err != nil {
		panic(err)
	}

}
