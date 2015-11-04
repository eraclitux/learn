package sml

import (
	"fmt"

	"github.com/eraclitux/stracer"
)

type ConfMatrix map[string]map[string]int

func (cm ConfMatrix) String() string {
	labels := make([]string, 0, len(cm))
	for k, _ := range cm {
		labels = append(labels, k)

	}
	strMatrix := ""
	for i, l := range labels {
		strMatrix += fmt.Sprintf("%18s(%d):", l, i+1)
		for _, j := range labels {
			strMatrix += fmt.Sprintf("%12d", cm[l][j])
		}
		strMatrix += "\n"

	}
	return strMatrix
}

// FIXME
func Report() {}

// ConfMatrix computes confusion matrix. It expects predict Table
// to store predicted labels in single field rows.
func ConfusionMatrix(expect, predict Table) (ConfMatrix, error) {
	var confMatrix ConfMatrix = make(map[string]map[string]int)
	for i := 0; i < expect.Len(); i++ {
		row, err := expect.Row(i)
		if err != nil {
			return nil, err
		}
		pRow, err := predict.Row(i)
		if err != nil {
			return nil, err
		}
		expectedLabel := row[len(row)-1].(string)
		predictedLabel := pRow[0].(string)
		if m, ok := confMatrix[expectedLabel]; ok {
			m[predictedLabel]++
		} else {
			confMatrix[expectedLabel] = make(map[string]int)
		}
	}
	// get all labels
	labels := make([]string, 0, len(confMatrix))
	for k, _ := range confMatrix {
		labels = append(labels, k)

	}
	// add 0 elements labels
	for k, m := range confMatrix {
		for _, l := range labels {
			if _, ok := confMatrix[k][l]; !ok {
				m[l] = 0
			}
		}
	}
	stracer.Traceln("confMatrix:", map[string]map[string]int(confMatrix))
	return confMatrix, nil
}
