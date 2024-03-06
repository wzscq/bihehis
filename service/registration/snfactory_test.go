package registration

import (
	"log/slog"
	"os"
	"testing"
)

func TestGetSN(t *testing.T) {
	factory := &DefaultSNFactory{
		SNPool: make(map[string]SNItem),
	}

	sn := factory.GetSN("test")
	slog.Info("GetSN", "result", sn)
	if sn != "0001" {
		t.Error("get sn failed")
		return
	}

	sn = factory.GetSN("test")
	slog.Info("GetSN", "result", sn)
	if sn != "0002" {
		t.Error("get sn failed")
		return
	}
}

func TestSavePool(t *testing.T) {
	factory := &DefaultSNFactory{
		SNPool: make(map[string]SNItem),
	}

	factory.GetSN("test")
	factory.GetSN("test")
	factory.SavePool("test.json")

	factory.LoadPool("test.json")
	sn := factory.GetSN("test")
	slog.Info("GetSN", "result", sn)
	if sn != "0003" {
		t.Error("get sn failed")
		return
	}

	//remove the file
	err := os.Remove("test.json")
	if err != nil {
		slog.Error("remove file failed", "error", err.Error())
		t.Error("remove file failed")
	}
}
