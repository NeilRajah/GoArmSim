//util
//Author: Neil Balaskandarajah
//Created on: 09/23/2019
//A collection of useful utility methods

package main

import (
	"math"
)

//Constants

const g = 9.81            //gravitational constant
const pixelToMeters = 500 //the number of pixels per 'meter' in the GUI

//ToDegrees converts an angle from radians to degrees
//float64 n - the angle in radians
//return - the angle in degrees
func ToDegrees(n float64) float64 {
	return (n * 180.0) / (math.Pi)
} //end toDegrees

//ToRadians converts an angle from degrees to radians
//float64 n - the angle in degrees
//return - the angle in radians
func ToRadians(n float64) float64 {
	return (n * math.Pi) / 180.0
} //end toRadians

//OutputClamp clamps a value between a min and a max
//float64 n - output to clamp
//float64 min - the bottom value of the clamp
//float64 max - the top value of the clamp
func OutputClamp(n, min, max float64) float64 {
	if math.Signbit(n) { //if negative
		return math.Max(n, min)
	}
	return math.Min(n, max) //if positive
} //end outputClamp
