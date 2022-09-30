package multichecker_test

import (
	"testing"

	mainmultichecker "github.com/borisbbtest/go_home_work/internal/multichecker"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), mainmultichecker.Analyzer, "./...")
}
