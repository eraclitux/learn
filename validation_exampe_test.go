package learn_test

import (
	"fmt"
	"log"

	"github.com/eraclitux/learn"
)

func ExampleValidate() {
	rC, err := learn.LoadCSV("datasets/iris_train.csv")
	if err != nil {
		log.Fatal(err)
	}
	// Load train set in memory.
	trainSet, err := learn.Normalize(rC)
	if err != nil {
		log.Fatal(err)
	}
	rC, err = learn.LoadCSV("datasets/iris_test.csv")
	if err != nil {
		log.Fatal(err)
	}
	testSet, err := learn.Normalize(rC)
	if err != nil {
		log.Fatal(err)
	}
	clf := learn.NewkNN(trainSet, 5)
	//
	// Cross validation
	//
	predictedLabels, err := clf.Predict(testSet)
	if err != nil {
		log.Fatal(err)
	}
	confMatrix, err := learn.ConfusionM(testSet, predictedLabels)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(confMatrix)
	report := learn.Validate(confMatrix)
	fmt.Println(report)

	// OUTPUT:
	//        versicolor(1):           5           0           2
	//            setosa(2):           0           5           0
	//	virginica(3):           0           0           3
	//     feature | precision | recall |
	//   virginica |       1.0 |    0.6 |
	//  versicolor |       0.7 |    1.0 |
	//      setosa |       1.0 |    1.0 |
}
