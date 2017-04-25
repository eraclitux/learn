// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"bytes"
	"fmt"
)

type ConfMatrix struct {
	mm map[string]map[string]int
}

// BUG(eraclitux): order labels to
// have a predictable output.
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
	return strMatrix
}

type ValidationReport struct {
	Labels map[string]Validation
	// Overall
	Accuracy float64
}

type Validation struct {
	Precision float64
	Recall    float64
}

// BUG(eraclitux): order labels to
// have a predictable output.
func (r ValidationReport) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%12s | %-6s | %-6s |\n", "feature", "precision", "recall")
	for k, v := range r.Labels {
		fmt.Fprintf(&buf, "%12s | %9.2f | %6.2f |\n", k, v.Precision, v.Recall)
	}
	fmt.Fprintf(&buf, "Overall accuracy: %.2f", r.Accuracy)
	return buf.String()
}

// Validate computes precision, recall and
// overall accuracy.
func Validate(cm ConfMatrix) ValidationReport {
	vr := ValidationReport{
		Labels:   make(map[string]Validation),
		Accuracy: computeAccuracy(cm),
	}
	for k, v := range cm.mm {
		vr.Labels[k] = Validation{
			Precision: computePrecision(k, v),
			Recall:    computeRecall(k, cm),
		}
	}
	return vr
}

// ConfusionM computes confusion matrix.
// predict Table must store
// labels in single field rows, expected labels
// are taken from the last field of expect's rows.
func ConfusionM(expect, predict Table) (ConfMatrix, error) {
	var confMatrix = ConfMatrix{mm: make(map[string]map[string]int)}
	nRows, _ := expect.Caps()
	for i := 0; i < nRows; i++ {
		row, err := expect.Row(i)
		if err != nil {
			return ConfMatrix{}, err
		}
		pRow, err := predict.Row(i)
		if err != nil {
			return ConfMatrix{}, err
		}
		expectedLabel, ok := row[len(row)-1].(string)
		if !ok {
			return ConfMatrix{}, fmt.Errorf("learn: %v is not a string", expectedLabel)
		}
		predictedLabel, ok := pRow[0].(string)
		if !ok {
			return ConfMatrix{}, fmt.Errorf("learn: %v is not a string", predictedLabel)
		}
		if m, ok := confMatrix.mm[expectedLabel]; ok {
			m[predictedLabel]++
		} else {
			confMatrix.mm[expectedLabel] = make(map[string]int)
			confMatrix.mm[expectedLabel][predictedLabel]++
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
