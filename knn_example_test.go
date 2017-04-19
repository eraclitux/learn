package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleNewkNN() {
	rC, err := learn.LoadCSV("datasets/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	// Load train set in memory.
	trainSet, err := learn.Normalize(rC)
	if err != nil {
		log.Fatal(err)
	}
	clf := learn.NewkNN(trainSet, 5)
	// Categorize single sample.
	var tab learn.MemoryTable = make([][]interface{}, 1)
	// FIXME use Normalize after refactoring
	tab[0] = []interface{}{0.2, 0.62, 0.07, 0.04}
	prediction, err := clf.Predict(tab)
	if err != nil {
		log.Fatal(err)
	}
	r, err := prediction.Row(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("predicted category:", r[0])

	// OUTPUT:
	// predicted category: setosa
}
