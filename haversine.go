package main

import (
	"github.com/umahmood/haversine"
)

func Distance(here, there Location) float64 {
	d1 := haversine.Coord{Lat: here.Lat, Lon: here.Long}
	d2 := haversine.Coord{Lat: there.Lat, Lon: there.Long}

	_, km := haversine.Distance(d1, d2)

	return km * 1000
}
