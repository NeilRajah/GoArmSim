//util
//Author: Neil Balaskandarajah
//Created on: 09/23/2019
//A collection of useful utility methods

package main

import (
	"github.com/faiface/pixel"
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

//Check if two points are within the bounds of one another
//Point target - target to be within
//Point current - current (x,y) coordinate
//tolerance float64 - radius of the circle around the target you can be in
func withinBounds(target, current Point, tolerance float64) bool {
	xDiff := math.Abs(target.x - current.x)
	yDiff := math.Abs(target.y - current.y)

	if xDiff < tolerance && yDiff < tolerance {
		return true
	}
	return false
} //end withinBounds

//Scale a point by a value
//Point p - point to scale
//float64 scale - value to scale the coordinates of the point by
//return - scaled point
func scalePoint(p Point, scale float64) Point {
	return Point{p.x * scale, p.y * scale}
} //end scalePoint

//Convert a mouse point to a cartesian point the arm can move to
func mouseToCartesian(m pixel.Vec) Point {
	return scalePoint(Point{m.X - float64(width)/2, m.Y}, 1.0/float64(pixelToMeters))
} //end MouseToCartesian

//ClampToCSpace keeps a point within a 2-jointed arm's configuration space
//Point p - point to clamp
//float64 l1 - first joint's length
//float64 l2 - second joint's length
//return - the clamped point
func ClampToCSpace(p Point, l1, l2 float64) Point {
	val := p.x*p.x + p.y*p.y //x^2 + y^2
	in := (l1 - l2) * (l1 - l2)
	out := (l1 + l2) * (l1 + l2)

	if val > in && val < out { //within extremes of configuration space
		return p
	}
	//not within space
	theta := math.Atan2(p.y, p.x) //angle from + x-axis to point

	r := 0.0      //length of arm for point on edge of c-space
	if val < in { //too small
		r = l1 - l2 //edge of inner
		r *= 1.001
	} else { //too big
		r = l1 + l2 //edge of outer
		r *= 0.999
	} //if
	//r is scaled up/down by small amount to ensure point is within c-space and not slightly outside due to rounding error

	return Point{r * math.Cos(theta), r * math.Sin(theta)}
} //end clampToCSpace
