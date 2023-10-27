package main

import (
	"fmt"
	"math"
	"time"
)

var A, B, C float64

var width, height float64 = 160, 44
var zBuffer [160 * 44]float64
var buffer [160 * 44]uint8
var backgroundASCIICode uint8 = '.'
var distanceFromCam float64 = 100
var horizontalOffset float64
var K1 float64 = 40
var x, y, z float64

var incrementSpeed float64 = 0.6

func sin(theta float64) float64 {
	return math.Sin(theta)
}

func cos(theta float64) float64 {
	return math.Cos(theta)
}

func calculateX(i float64, j float64, k float64) float64 {
	return j*sin(A)*sin(B)*cos(C) - k*cos(A)*sin(B)*cos(C) +
		j*cos(A)*sin(C) + k*sin(A)*sin(C) + i*cos(B)*cos(C)
}

func calculateY(i float64, j float64, k float64) float64 {
	return j*cos(A)*cos(C) + k*sin(A)*cos(C) -
		j*sin(A)*sin(B)*sin(C) + k*cos(A)*sin(B)*sin(C) -
		i*cos(B)*sin(C)
}

func calculateZ(i float64, j float64, k float64) float64 {
	return k*cos(A)*cos(B) - j*sin(A)*cos(B) + i*sin(B)
}

func calculateForSurface(cubeX float64, cubeY float64, cubeZ float64, ch uint8) {
	x = calculateX(cubeX, cubeY, cubeZ)
	y = calculateY(cubeX, cubeY, cubeZ)
	z = calculateZ(cubeX, cubeY, cubeZ) + distanceFromCam

	ooz := 1 / z

	xp := int(width/2 + horizontalOffset + K1*ooz*x*2)
	yp := int(height/2 + K1*ooz*y)

	idx := xp + yp*int(width)
	if idx >= 0 && float64(idx) < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
		}
	}
}

func refresh() {
	w, h := int(width), int(height)
	for i := 0; i < w*h; i++ {
		zBuffer[i] = 0
	}
	for i := 0; i < w*h; i++ {
		buffer[i] = backgroundASCIICode
	}
}

func main() {
	fmt.Printf("\x1b[2J")
	for {
		refresh()
		cubeWidth := 20.0
		horizontalOffset = -2 * cubeWidth
		for cubeX := -cubeWidth; cubeX < cubeWidth; cubeX += incrementSpeed {
			for cubeY := -cubeWidth; cubeY < cubeWidth; cubeY += incrementSpeed {
				calculateForSurface(cubeX, cubeY, -cubeWidth, '@')
				calculateForSurface(cubeWidth, cubeY, cubeX, '$')
				calculateForSurface(-cubeWidth, cubeY, -cubeX, '~')
				calculateForSurface(-cubeX, cubeY, cubeWidth, '#')
				calculateForSurface(cubeX, -cubeWidth, -cubeY, ';')
				calculateForSurface(cubeX, cubeWidth, cubeY, '+')
			}
		}
		w, h := int(width), int(height)
		fmt.Printf("\x1b[H")
		for i := 0; i < w*h; i++ {
			if i%w == 0 {
				fmt.Printf("\n")
			} else {
				fmt.Printf("%c", buffer[i])
			}
		}
		A += 0.05
		B += 0.05
		C += 0.01
		time.Sleep(time.Millisecond * 16)
	}
}
