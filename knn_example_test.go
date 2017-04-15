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
	// Load all data in memory.
	trainData, err := learn.Normalize(rC)
	if err != nil {
		log.Fatal(err)
	}
	cf := learn.NewkNN(trainData, 5)
	var tab learn.MemoryTable = make([][]interface{}, 1)
	// FIXME use Normalize after refactoring
	tab[0] = []interface{}{0.1, 0.1, 0.1, 0.1}
	prediction, err := cf.Predict(tab)
	if err != nil {
		log.Fatal(err)
	}
	r, err := prediction.Row(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r[0])

	// OUTPUT:
	// setosa
}
