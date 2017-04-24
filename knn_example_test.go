package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleNewkNN() {
	trainSet, err := learn.ReadAllCSV("datasets/iris.csv")
	if err != nil {
		log.Fatal(err)
	}
	mu, sigma, err := learn.Normalize(trainSet, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	clf := learn.NewkNN(trainSet, 3)
	// Categorize single sample.
	var testSet learn.MemoryTable = make([][]interface{}, 1)
	testSet[0] = []interface{}{5.2, 3.4, 1.3, 0.1}
	_, _, err = learn.Normalize(testSet, mu, sigma)
	prediction, err := clf.Predict(testSet)
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
