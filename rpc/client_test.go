package rpc

import "testing"

func TestFormatBytes(t *testing.T) {
	if got := FormatBytes(512); got != "512 B" {
		t.Fatalf("got %q", got)
	}
	if got := FormatBytes(1536); got != "1.5 KB" {
		t.Fatalf("got %q", got)
	}
}

func TestTaskProgress(t *testing.T) {
	task := Task{TotalLength: "100", CompletedLength: "25"}
	if p := task.Progress(); p < 0.24 || p > 0.26 {
		t.Fatalf("progress=%v", p)
	}
	complete := Task{Status: "complete"}
	if complete.Progress() != 1 {
		t.Fatalf("complete progress=%v", complete.Progress())
	}
}
