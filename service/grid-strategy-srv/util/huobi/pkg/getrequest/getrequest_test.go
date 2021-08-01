package getrequest

import (
	"testing"
)

func TestAddParam_PropertyEmpty_ReturnEmpty(t *testing.T) {
	request := new(GetRequest).Init()

	request.AddParam("", "value")
	result := request.BuildParams()

	expected := ""
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestAddParam_ValueEmpty_ReturnEmpty(t *testing.T) {
	request := new(GetRequest).Init()

	request.AddParam("property", "")
	result := request.BuildParams()

	expected := ""
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestBuildParam_NullParam_ReturnEmpty(t *testing.T) {
	request := new(GetRequest).Init()

	result := request.BuildParams()

	expected := ""
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestBuildParam_OneParam_ReturnOnePair(t *testing.T) {
	request := new(GetRequest).Init()
	request.AddParam("key", "value")

	result := request.BuildParams()

	expected := "key=value"

	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestBuildParams_UnEscapedParam_ReturnEscapedParam(t *testing.T) {
	request := new(GetRequest).Init()
	request.AddParam("key", "valueA:valueB/valueC=")

	result := request.BuildParams()

	expected := "key=valueA%3AvalueB%2FvalueC%3D"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}

func TestBuildParam_TwoParams_ReturnOrderedTwoPairs(t *testing.T) {
	request := new(GetRequest).Init()
	request.AddParam("id", "123")
	request.AddParam("Year", "2020")

	result := request.BuildParams()

	expected := "Year=2020&id=123"
	if result != expected {
		t.Errorf("expected: %s, actual: %s", expected, result)
	}
}
