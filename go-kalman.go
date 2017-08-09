package main

import (
	"fmt"
	"os"
	"gocsv"
	"image"
	"image/color"
	"image/png"
)

var Q_angle = 150.0
var Q_gyro = 0.001118999
var R_angle = 100.0

var x_angle = 0.0
var x_bias = 0.0
var P_00 = 0.0
var P_01 = 0.0
var P_10 = 0.0
var P_11 = 0.0


func KalmanCalculate(newAngle, newRate, Looptime float64) (float64) {

	dtk := (Looptime / 1000)
	fmt.Printf("\nlooptime %f",dtk)
	x_angle += dtk * (newRate - x_bias)
	P_00 += - dtk * (P_10 + P_01) + Q_angle * dtk
	P_01 += - dtk * P_11
	P_10 += - dtk * P_11
	P_11 += + Q_gyro * dtk

	y := newAngle - x_angle
	S := P_00 + R_angle
	K_0 := P_00 / S
	K_1 := P_10 / S

	x_angle += K_0 * y
	x_bias += K_1 * y
	P_00 -= K_0 * P_00
	P_01 -= K_0 * P_01
	P_10 -= K_1 * P_00
	P_11 -= K_1 * P_01

	return (x_angle)
}

type fuel_data struct {
	Tis  uint64 `csv:"tis"`
	Fuel float64 `csv:"data"`
}

func main() {

	i:=0
	filename := "filename.csv"
	LocationsFile, err := os.OpenFile(filename, os.O_RDONLY | os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer LocationsFile.Close()

	FuelData := []*fuel_data{}

	if err := gocsv.UnmarshalFile(LocationsFile, &FuelData); err != nil {
		// Load Locations from file
		panic(err)
	}
	fmt.Printf("Done Marshalling Length of fuel data %d\n",len(FuelData))
	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{len(FuelData),500}})
	var temp_tis = uint64(0)
	for _, temp_pos := range FuelData {

		val := temp_pos.Tis - temp_tis
		plot := int(20 * KalmanCalculate(temp_pos.Fuel, 0, float64(val)))
		plot_raw := int(20 * temp_pos.Fuel)
		m.SetNRGBA(i,(plot*-1)+200,color.NRGBA{0, 0, 0, 255})
		m.SetNRGBA(i,(plot_raw*-1)+350,color.NRGBA{255,77,165, 255})
		temp_tis = temp_pos.Tis
		i++
		fmt.Printf("\nVal: %d Raw %d",plot,plot_raw)
	}
	LocationsFile.Close()
  //Creating Image
	g, err_png := os.OpenFile(filename + ".png", os.O_CREATE|os.O_WRONLY, 0666)
	if err_png != nil {
		fmt.Println(err_png)
		os.Exit(1)
	}
  //Render to image
	if err_png = png.Encode(g, m); err_png != nil {
		fmt.Printf("Error : %s",err_png)
		os.Exit(1)
	}
}
