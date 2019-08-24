package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"time"
)

type XYs []XY

type XY struct {
	X, Y float64
}

type xyAxis struct {
	x float64
	y float64
}


func MakeTemperatureGraph(ts *[]Temperatures) {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	p.Title.Text = "Changes in room temperature"
	p.X.Tick.Marker = plot.TimeTicks{Format: "15:04"}
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Temperature"
	p.Add(plotter.NewGrid())

	var nums []xyAxis
	for _, v := range *ts {
		date, err := time.Parse("2006-01-02 15:04", v.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, xyAxis{x: float64(date.Unix()), y: v.Temperature})
	}

	pts := make(plotter.XYs, len(nums))

	for i, axis := range nums {
		pts[i].X = axis.x
		pts[i].Y = axis.y
	}

	err = plotutil.AddLinePoints(p, pts)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Save(15*vg.Inch, 5*vg.Inch, imgDir + "graph.png"); err != nil {
		log.Fatal(err)
	}
}
