package effects

import "testing"

type TestStruct struct {
	Abc int
}

func TestSimpleApply(t *testing.T) {
	control := &BaseControl{
		Field:   "Abc",
		Initial: 123,
	}

	control.Init()

	testStruct := &TestStruct{}
	control.Apply(testStruct)

	if testStruct.Abc != 123 {
		t.Error("Abc != 123", testStruct.Abc)
	}
}

func TestSimpleTransform(t *testing.T) {
	control := &BaseControl{
		Field:     "Abc",
		Initial:   123,
		Transform: "v+50",
	}

	control.Init()

	testStruct := &TestStruct{Abc: 0}
	control.Apply(testStruct)

	if testStruct.Abc != 173 {
		t.Error("Abc != 53", testStruct.Abc)
	}
}

func BenchmarkTransform(b *testing.B) {
	control := &BaseControl{
		Field:     "Abc",
		Initial:   1,
		Transform: "v+123",
	}

	control.Init()

	testStruct := &TestStruct{Abc: 5}

	control.Apply(testStruct)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		control.Apply(testStruct)
	}

	if testStruct.Abc == 5 {
		b.Error("wrong value", testStruct.Abc)
	}
}
