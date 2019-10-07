//arm2
//Author: Neil Balaskandarajah
//Created on: 10/06/2019
//A two degree of freedom arm

package main

import (
// "fmt"
)

//A two degree of freedom arm is made of two Arm structs, with the elbow dynamically changing start position
type Arm2 struct {
	arm1 *Arm //the base joint (shoulder)
	arm2 *Arm //the second joint (elbow)
} //end struct

//Updates the position of the arms, translating the second joint start to the first joint end
func (a2 Arm2) update() {
	a2.arm2.setStartPt(a2.arm1.getEndPtPxl())
} //end update
