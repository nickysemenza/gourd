package app

import (
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	h "github.com/nickysemenza/food/backend/app/handler"
)

func TestApp_Initialize(t *testing.T) {
	type fields struct {
		R *mux.Router
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *h.Env
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				R: tt.fields.R,
			}
			if got := a.Initialize(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_buildRoutes(t *testing.T) {
	type fields struct {
		R *mux.Router
	}
	type args struct {
		env *h.Env
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				R: tt.fields.R,
			}
			a.buildRoutes(tt.args.env)
		})
	}
}

func TestApp_RunServer(t *testing.T) {
	type fields struct {
		R *mux.Router
	}
	type args struct {
		host string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				R: tt.fields.R,
			}
			a.RunServer(tt.args.host)
		})
	}
}
