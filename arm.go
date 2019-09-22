//arm
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a single-jointed arm in cartesian space

package main

import (
	"math"
)

//A single-jointed arm in cartesian space
type arm struct {
	start  point   //the base point of the arm
	length float64 //the length of the arm
	angle  float64 //the angle of the arm from the horizontal measured CCW
} //end struct

//Get the end point of the arm
func (a arm) getEndPt() point {
	endX := a.length*math.Cos(a.angle) + a.start.x
	endY := a.length*math.Sin(a.angle) + a.start.y
	return point{endX, endY}
} //end GetEndPt
