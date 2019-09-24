//util
//Author: Neil Balaskandarajah
//Created on: 09/23/2019
//A collection of useful utility methods

package main

import (
	"math"
)

//convert an angle from radians to degrees
//float64 n - the angle in radians
//return - the angle in degrees
func ToDegrees(n float64) float64 {
	return (n * 180.0) / (math.Pi)
} //end toDegrees

//convert an angle from degrees to radians
//float64 n - the angle in degrees
//return - the angle in radians
func ToRadians(n float64) float64 {
	return (n * math.Pi) / 180.0
} //end toRadians
