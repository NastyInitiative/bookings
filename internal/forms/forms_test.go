package forms

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/pippo", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error(":: Got invalid when it should have been valid ::")
	}
}

func TestForm_Required(t *testing.T) {

	// First case - validated field does not exist
	r := httptest.NewRequest("POST", "/pippo", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error(":: Form shows valid when required fields missing")
	}

	// Second case - no data is passed inside required fields
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	// Thirs case - it has the required data
	form = New(postData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error(":: Shows it does not have required fields when it does ")
	}

}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	// First case - validated field does not exist
	form.MinLength("a", 3)
	if form.Valid() {
		t.Error(":: Length should be invalid because field is missing, got valid instead ::")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "aaa")

	// Second case - a has invalid length
	form = New(postData)
	form.MinLength("a", 3)
	l := len(postData.Get("a"))
	if form.Valid() {
		t.Errorf(fmt.Sprintf(":: Length of %s should be %d, got %d instead ::", "a", 3, l))
	}

	// Third case - b has valid length
	form = New(postData)
	form.MinLength("b", 3)
	if !form.Valid() {
		t.Errorf(":: Length of %s should be valid, got invalid instead ::", "b")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postData := url.Values{}
	form := New(postData)

	// First case - validated field does not exist
	form.IsEmail("a")
	if form.Valid() {
		t.Error(":: Email should be invalid since no value is passed, got valid instead ::")
	}

	postData = url.Values{}
	postData.Add("a", "a@a.com")
	postData.Add("b", "a@a")

	// Second case - email format is invalid
	form = New(postData)
	form.IsEmail("b")
	if form.Valid() {
		t.Error(":: Email should be invalid since the format doesn't match, got valid instead ::")
	}

	// Third case - email format is valid
	form = New(postData)
	form.IsEmail("a")
	if !form.Valid() {
		t.Error(":: Email should be valid, got invalid instead ::")
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	// First case - field does not exists
	form := New(postData)
	has := form.Has("whatever")
	if has {
		t.Error("form shows it has the field when in doesn't")
	}

	postData = url.Values{}
	postData.Add("a", "a")

	// Second case - field does exists
	form = New(postData)
	has = form.Has("a")
	if !has {
		t.Error(":: form shows it does not have field when it does")
	}
}
