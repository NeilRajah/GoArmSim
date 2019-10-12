//arm2
//Author: Neil Balaskandarajah
//Created on: 10/06/2019
//A two degree of freedom arm

package main

import (
	// "fmt"
	"math"
)

//A two degree of freedom arm is made of two Arm structs, with the elbow dynamically changing start position
type Arm2 struct {
	arm1 *Arm //the base joint (shoulder)
	arm2 *Arm //the second joint (elbow)
} //end struct

//Updates the position of the arms, translating the second joint start to the first joint end
func (a2 *Arm2) update() {
	a2.arm2.setStartPt(a2.arm1.getEndPtPxl())
} //end update

//updates the individual arms
func (a2 Arm2) rest() {
	a2.update()
	a2.arm1.update()
	a2.arm2.update()
} //end rest

func (a2 Arm2) moveToPoint(p Point, tolerance float64) {
	d := PointDistance(Point{0, 0}, p)
	// gamma := math.Atan2(p.y, p.x)
	// alpha := cosLawAngle(a2.arm1.length, d, a2.arm2.length)
	// beta := cosLawAngle(a2.arm1.length, a2.arm2.length, d)

	// ang1 := gamma - alpha
	// ang2 := math.Pi - beta

	// ang1 := gamma + alpha
	// ang2 := beta - math.Pi

	// fmt.Println(d, ToDegrees(gamma), ToDegrees(alpha), ToDegrees(beta), ToDegrees(ang1), ToDegrees(ang2))

	ang2 := math.Acos((p.x*p.x + p.y*p.y - d*d) / (2 * a2.arm1.length * a2.arm2.length))
	ang1 := math.Atan(p.y/p.x) - math.Atan((a2.arm2.length*math.Sin(ang2))/(a2.arm1.length+a2.arm2.length*math.Cos(ang2)))

	a2.arm1.movePIDFF(ang1, a2.arm1.angle, tolerance)
	a2.arm2.movePIDFF(ang2, a2.arm2.angle, tolerance)
}
