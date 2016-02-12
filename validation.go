// Copyright (c) 2015 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"fmt"
)

type ConfMatrix struct {
	mm map[string]map[string]int
}

func computeAccuracy(cm ConfMatrix) float64 {
	correctN := 0.0
	total := 0.0
	for k, v := range cm.mm {
		for q, n := range v {
			// Sum all elements on diagonal
			total += float64(n)
			if q == k {
				correctN += float64(n)
			}
		}
	}
	return correctN / total
}

func (cm ConfMatrix) String() string {
	// FIXME deal with mm == nil
	labels := make([]string, 0, len(cm.mm))
	for k, _ := range cm.mm {
		labels = append(labels, k)

	}
	strMatrix := ""
	for i, l := range labels {
		strMatrix += fmt.Sprintf("%18s(%d):", l, i+1)
		for _, j := range labels {
			strMatrix += fmt.Sprintf("%12d", cm.mm[l][j])
		}
		strMatrix += "\n"

	}
	strMatrix += "\n"
	strMatrix += fmt.Sprintf("Overall accuracy: %f", computeAccuracy(cm))
	return strMatrix
}

type ValidationReport map[string]validation

type validation struct {
	Precision float64
	Recall    float64
}

func computePrecision(k string, m map[string]int) float64 {
	tp := float64(m[k])
	tpPlusfp := 0.0
	for _, v := range m {
		tpPlusfp += float64(v)
	}
	return tp / tpPlusfp

}

// more false negatives => bigger recall
//	tp	fp
//	fn	tn
func computeRecall(k string, cm ConfMatrix) float64 {
	tp := 0.0
	totalFn := 0.0
	for q, v := range cm.mm {
		// This is tp
		if q == k {
			tp = float64(v[k])
			continue
		} else {
			totalFn += float64(v[k])
		}
	}
	return tp / (tp + totalFn)

}

// Report calculates precision and recall for a classification
// result.
func Report(cm ConfMatrix) ValidationReport {
	var vr ValidationReport = make(map[string]validation)
	for k, v := range cm.mm {
		vr[k] = validation{
			Precision: computePrecision(k, v),
			Recall:    computeRecall(k, cm),
		}
	}
	return vr
}

// ConfMatrix computes confusion matrix. It expects predict Table
// to store predicted labels in single field rows.
func ConfusionMatrix(expect, predict Table) (ConfMatrix, error) {
	var confMatrix = ConfMatrix{mm: make(map[string]map[string]int)}
	for i := 0; i < expect.Len(); i++ {
		row, err := expect.Row(i)
		if err != nil {
			return ConfMatrix{}, err
		}
		pRow, err := predict.Row(i)
		if err != nil {
			return ConfMatrix{}, err
		}
		expectedLabel := row[len(row)-1].(string)
		predictedLabel := pRow[0].(string)
		if m, ok := confMatrix.mm[expectedLabel]; ok {
			m[predictedLabel]++
		} else {
			confMatrix.mm[expectedLabel] = make(map[string]int)
		}
	}
	// get all labels
	labels := make([]string, 0, len(confMatrix.mm))
	for k, _ := range confMatrix.mm {
		labels = append(labels, k)

	}
	// add 0 elements labels
	for k, m := range confMatrix.mm {
		for _, l := range labels {
			if _, ok := confMatrix.mm[k][l]; !ok {
				m[l] = 0
			}
		}
	}
	return confMatrix, nil
}
