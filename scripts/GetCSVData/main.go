package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("./data1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all the rows in the Sheet1.
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range cols[2] {
		fmt.Println(v)
	}
}


