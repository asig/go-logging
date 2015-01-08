package logging

import "testing"

func TestParseLogLevels(t *testing.T) {
	const levelSpec = "INFO,foo=DEBUG,foo/bar=SEVERE"
	parseLogLevelSpec(levelSpec)

	assert(len(levelPerComponent) == 3, "Wrong number of components", t)
	assert(levelPerComponent[""] == INFO, "Wrong level for ''", t)
	assert(levelPerComponent["foo"] == DEBUG, "Wrong level for 'foo'", t)
	assert(levelPerComponent["foo/bar"] == SEVERE, "Wrong level for 'foo/bar'", t)
}

func TestParseLogLevels_NoDefaultLevel(t *testing.T) {
	const levelSpec = "foo=DEBUG,foo/bar=SEVERE"
	parseLogLevelSpec(levelSpec)

	assert(len(levelPerComponent) == 3, "Wrong number of components", t)
	assert(levelPerComponent[""] == WARNING, "Wrong level for ''", t)
	assert(levelPerComponent["foo"] == DEBUG, "Wrong level for 'foo'", t)
	assert(levelPerComponent["foo/bar"] == SEVERE, "Wrong level for 'foo/bar'", t)
}

func TestGetLogLevels(t *testing.T) {
	const levelSpec = "WARNING,foo=DEBUG,foo/bar=SEVERE"
	parseLogLevelSpec(levelSpec)

	assert(getComponentLevel("foo/bar") == SEVERE, "Wrong level for 'foo/bar'", t)
	assert(getComponentLevel("foo/whatnot") == DEBUG, "Wrong level for 'foo/whatnot'", t)
	assert(getComponentLevel("foo") == DEBUG, "Wrong level for 'foo'", t)
	assert(getComponentLevel("bar") == WARNING, "Wrong level for 'bar'", t)
	assert(getComponentLevel("") == WARNING, "Wrong level for ''", t)
}

func assert(condition bool, msg string, t *testing.T) {
	if !condition {
		t.Error(msg)
	}
}
