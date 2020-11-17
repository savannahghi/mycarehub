package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// can do setup here
	code := m.Run()
	// can do cleanup here
	os.Exit(code)
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "check that things are wired up properly",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				log.Printf("about to start main()...")
				main()
			}()
			time.Sleep(time.Second * 10) //  the 10s delay allows time for things to work or blow  up
		})
	}
}
