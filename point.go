//point
//Author: Neil Balaskandarajah
//Created on: 09/21/2019
//struct to represent a point in cartesian space

package main

//A point in cartesian space
type point struct {
	x float64 //x coordinate of the point
	y float64 //y coordinate of the point
} //end struct

//set the x position of the point
//float64 newX - the new x coordinate of the point in cartesian space
func (p *point) setX(newX float64) {
	p.x = newX
} //end setX

//set the y position of the point
//float64 newY - the new y coordinate of the point in cartesian space
func (p *point) setY(newY float64) {
	p.y = newY
} //end setY
