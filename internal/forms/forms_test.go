package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/make-reservation", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/make-reservation", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r = httptest.NewRequest("POST", "/make-reservation", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Não é válido quando era para ser válido")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/make-reservation", nil)
	form := New(r.PostForm)

	form.Has("a")
	if form.Valid() {
		t.Error("Deu válido quando não deveria ser")
	}

	postedData := url.Values{}
	postedData.Add("a", "s")

	r = httptest.NewRequest("POST", "/make-reservation", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Has("a")

	if !form.Valid() {
		t.Error("Está inválido quando deveria estar válido")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/make-reservation", nil)
	form := New(r.PostForm)

	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("Deu Válido quando era para ser inválido, formulário, vazio")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("a", "sasas")

	r = httptest.NewRequest("POST", "/make-reservation", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("Deu inválido quando era para ser válido")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_Email(t *testing.T) {
	r := httptest.NewRequest("POST", "/make-reservation", nil)
	form := New(r.PostForm)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Deu Válido quando era para ser inválido, formulário, vazio")
	}

	postedData := url.Values{}
	postedData.Add("email", "me@me.com")

	r = httptest.NewRequest("POST", "/make-reservation", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Deu inválido quando era para ser válido")
	}
}
