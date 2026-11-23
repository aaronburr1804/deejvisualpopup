package app

import (
	"testing"
)

func TestHUDLabel(t *testing.T) {
	tests := []struct {
		name    string
		targets []string
		rename  string
		want    string
	}{
		// device targets
		{"master", []string{"master"}, "", "Master"},
		{"master any case", []string{"MASTER"}, "", "Master"},
		{"mic", []string{"mic"}, "", "Mic"},
		// apps
		{"strip exe", []string{"discord.exe"}, "", "discord"},
		{"strip EXE", []string{"Discord.EXE"}, "", "Discord"},
		{"linux app", []string{"firefox"}, "", "firefox"},
		// unmapped
		{"unmapped raw", []string{"deej.unmapped"}, "", "deej.unmapped"},
		{"unmapped rename", []string{"deej.unmapped"}, "game", "game"},
		{"rename only unmapped", []string{"firefox"}, "game", "firefox"},
		// integrations: raw
		{"obs raw", []string{"deej.obs:Mic/Aux"}, "", "deej.obs:Mic/Aux"},
		// multiple targets, join
		{"joined", []string{"spotify.exe", "firefox.exe"}, "", "spotify + firefox"},
		{"joined mixed", []string{"spotify.exe", "deej.unmapped"}, "game", "spotify + game"},
		// short name
		{"short name", []string{"vlc"}, "", "vlc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hudLabel(tt.targets, tt.rename)
			if got != tt.want {
				t.Errorf("hudLabel(%q, %q) = %q, want %q", tt.targets, tt.rename, got, tt.want)
			}
		})
	}
}
