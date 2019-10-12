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
} //end OutputClamp

//CenterOfMass calculates the CoM between two points
//point p1 - first point
//float64 m1 - first point's mass
//point p2 - second point
//float64 m2 - second point's mass
func CenterOfMass(p1 Point, m1 float64, p2 Point, m2 float64) Point {
	x := (m1*p1.x + m2*p2.x) / (m1 + m2)
	y := (m1*p1.y + m2*p2.y) / (m1 + m2)
	return Point{x, y}
} //end CenterOfMass

//PointDistance calculates the length of the hypotenuse formed by two points
//Point p1 - first point
//Point p2 - second point
func PointDistance(p1, p2 Point) float64 {
	return math.Hypot(p2.x-p1.x, p2.y-p1.y)
} //end PointDistance

//Calculate the angle formed between two lines using cosine law
//float64 a - first length of triangle
//float64 b - second length of triangle
//float64 c - length opposite to angle formed by first two lengths
func cosLawAngle(a, b, c float64) float64 {
	return math.Acos((a*a + b*b - c*c) / (2 * a * b))
} //end cosLawAngle
