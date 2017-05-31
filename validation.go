// Copyright (c) 2017 Andrea Masi. All rights reserved.
// Use of this source code is governed by MIT license
// which that can be found in the LICENSE.txt file.

package learn

import (
	"bytes"
	"fmt"
	"sort"
)

// ConfMatrix stores confusion matrix
// needed to calculate precision and recall
// for classified labels.
type ConfMatrix struct {
	mm map[string]map[string]int
}

func (cm ConfMatrix) String() string {
	// FIXME deal with mm == nil
	labels := make([]string, 0, len(cm.mm))
	for k := range cm.mm {
		labels = append(labels, k)
	}
	// Sort labels to have
	// a predictable output.
	sort.Sort(sort.StringSlice(labels))
	var buf bytes.Buffer
	for i, l := range labels {
		fmt.Fprintf(&buf, "%18s(%d):", l, i+1)
		for _, j := range labels {
			fmt.Fprintf(&buf, "%12d", cm.mm[l][j])
		}
		fmt.Fprintf(&buf, "\n")
	}
	return buf.String()
}

// ValidationReport stores
// precision and recall for all the labels
// and the overall accuracy.
type ValidationReport struct {
	Labels   map[string]Validation
	Accuracy float64
}

// Validation stores validation data
// for a single label.
type Validation struct {
	Precision float64
	Recall    float64
}

func (r ValidationReport) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%12s | %-6s | %-6s |\n", "feature", "precision", "recall")
	// Sort labels to have
	// a predictable output.
	labels := make([]string, 0, len(r.Labels))
	for k := range r.Labels {
		labels = append(labels, k)
	}
	sort.Sort(sort.StringSlice(labels))
	for _, l := range labels {
		fmt.Fprintf(&buf, "%12s | %9.2f | %6.2f |\n", l, r.Labels[l].Precision, r.Labels[l].Recall)
	}
	fmt.Fprintf(&buf, "Overall accuracy: %.2f", r.Accuracy)
	return buf.String()
}

// Validate computes precision, recall and
// overall accuracy. Used for cross-validating
// Classifier.
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
		expectedLabel, ok := row[len(row)-1].(*category)
		if !ok {
			return ConfMatrix{}, fmt.Errorf("learn: %v is not a category", row[len(row)-1])
		}
		predictedLabel, ok := pRow[0].(string)
		if !ok {
			return ConfMatrix{}, fmt.Errorf("learn: %v is not a string", row[len(row)-1])
		}
		if m, ok := confMatrix.mm[expectedLabel.label]; ok {
			m[predictedLabel]++
		} else {
			confMatrix.mm[expectedLabel.label] = make(map[string]int)
			confMatrix.mm[expectedLabel.label][predictedLabel]++
		}
	}
	// get all labels
	labels := make([]string, 0, len(confMatrix.mm))
	for k := range confMatrix.mm {
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
