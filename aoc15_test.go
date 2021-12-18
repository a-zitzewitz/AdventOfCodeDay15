package main

import "testing"

func Test_findBestPath(t *testing.T) {
	doTest(t, "test-input.txt", 40, 315)
	doTest(t, "test-input2.txt", 21, 197)
}

func doTest(t *testing.T, fileName string, expectedRisk1 uint, expectedRisk2 uint) {
	lowestRisk1, lowestRisk2 := findBestPathForInput(fileName)
	if lowestRisk1 != expectedRisk1 {
		t.Errorf("Expected risk of %d for input '%s', but got %d", expectedRisk1, fileName, lowestRisk1)
	}
	if lowestRisk2 != expectedRisk2 {
		t.Errorf("Expected risk (enlarged map) of %d for input '%s', but got %d", expectedRisk2, fileName, lowestRisk2)
	}
}
