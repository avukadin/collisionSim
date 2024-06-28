package main

import (
	"math"
)


type Cell struct {
	balls []*Ball
}

type UniformGrid struct {
	cellSize float64
	grid [][]Cell
}

func NewUniforGrid(balls []Ball, cellSize float64) UniformGrid {
	nVerticalCells := int(math.Ceil(float64(screenHeight) / cellSize))
	nHorizontalCells := int(math.Ceil(float64(screenWidth) / cellSize))
	
	grid := make([][]Cell, nVerticalCells)
	for i := range grid {
		grid[i] = make([]Cell, nHorizontalCells)
	}
	ug := UniformGrid{
		cellSize: cellSize,
		grid: grid,
	}
	for i := range balls {
		cell := ug.getBallCell(&balls[i])
		cell.balls = append(cell.balls, &balls[i])
	}
	return ug
}

func (ug *UniformGrid) getBallCell(b *Ball) *Cell {
	row := int(b.r.y / ug.cellSize)
	col := int(b.r.x / ug.cellSize)
	row = int(math.Min(float64(len(ug.grid)-1), math.Max(0, float64(row))))
	col = int(math.Min(float64(len(ug.grid[0])-1), math.Max(0, float64(col))))
	return &ug.grid[row][col]
}

func (ug *UniformGrid) GetNeighbors(b *Ball) []*Ball {
	row := int(b.r.y / ug.cellSize)
	col := int(b.r.x / ug.cellSize)
	neighbors := []*Ball{}
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if row+i >= 0 && row+i < len(ug.grid) && col+j >= 0 && col+j < len(ug.grid[0]) {
				neighbors = append(neighbors, ug.grid[row+i][col+j].balls...)
			}
		}
	}
	return neighbors
}
