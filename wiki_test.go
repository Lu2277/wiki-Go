package main

import "testing"

var x = Page{
	Title: "hello",
	Body:  []byte("something"),
}

func TestLoad_save(t *testing.T) {
	type fields struct {
		Title string
		Body  []byte
	}
	tests := []struct {
		name    string
		filed   fields
		wantErr bool
	}{
		{"save", fields(x), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Page{
				Title: tt.filed.Title,
				Body:  tt.filed.Body,
			}
			if err := p.save(); (err != nil) != tt.wantErr {
				t.Errorf("page.save() error=%v,wantErr %v", err, tt.wantErr)

			}
		})
	}
}

func TestLoad_page(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		arg     args
		want    *Page
		wantErr bool
	}{
		{"load", args{title: x.Title}, &x, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get, err := loadPage(tt.arg.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadPage() error=%v,wantErr %v", get, tt.wantErr)

			}
			t.Logf("get body: %#v", string(get.Body))
		})
	}
}
