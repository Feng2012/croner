package croner

import (
	"testing"
	"time"
	"fmt"
)

var tmp = [5]string{}

func resetTmp() {
	tmp = [5]string{}
}

func resetRunner() {
	if DefaultManager.running{
		DefaultManager.RemoveAll()
		DefaultManager.Stop()
	}
}

type PanicJob struct {}

func (j PanicJob) Run() JobRunReturn{
	panic("hello, i am a panic job")
}

type GoodJob struct {}

func (j GoodJob) Run() JobRunReturn{
	return JobRunReturn{"Hello , I am a good job", nil }
}

type TimeOutJob struct {}

func (j TimeOutJob) Run() JobRunReturn {
	time.Sleep(5 * time.Second)
	return JobRunReturn{"Hello , I am a timeout job", nil }
}


func hookAppendResultToTmp(runReturn *JobRunReturn) {
	for i, v := range tmp{
		if v == ""{
			tmp[i] = fmt.Sprintf("%v", *runReturn)
			break
		}
	}
}

// test simple logic
func TestRunning(t *testing.T) {
	resetRunner()
	resetTmp()
	OnJobReturn(hookAppendResultToTmp)
	DefaultManager.Start()

	entryId, _ := DefaultManager.Add("@every 2s", GoodJob{})
	// sleep 3 second, tmp length should be 1
	time.Sleep(3 * time.Second)
	print("0000000")
	if tmp[1] != "" || tmp[0] == ""{
		t.FailNow()
	}
	// sleep 2 second again, tmp length should be 2
	time.Sleep(2 * time.Second)
	if tmp[2] != "" || tmp[1] == ""{
		t.FailNow()
	}
	// status should be "IDLE"
	if DefaultManager.JobMap[entryId].Status() != "IDLE" {
		t.FailNow()
	}

	// successTime should be 2, totalTime should be 2
	if DefaultManager.JobMap[entryId].successCount != 2 ||
		DefaultManager.JobMap[entryId].totalCount != 2 {
			t.FailNow()
	}
}

//  Test Ignore Panic = True
func TestIgnorePanic(t *testing.T) {
	resetRunner()
	resetTmp()
	OnJobReturn(hookAppendResultToTmp)
	DefaultManager.SetConfig(CronManagerConfig{true, false, 0, 0})
	DefaultManager.Start()

	entryId, _ := DefaultManager.Add("@every 2s", PanicJob{})
	// sleep 5 second, tmp length should be 2, Even this is a panic job
	time.Sleep(5 * time.Second)
	if tmp[2] != "" || tmp[1] == ""{
		t.FailNow()
	}

	// status should be "Fail"
	if DefaultManager.JobMap[entryId].Status() != "FAIL" {
		t.FailNow()
	}

	// successTime should be 0, totalTime should be 2
	if DefaultManager.JobMap[entryId].successCount != 0 ||
		DefaultManager.JobMap[entryId].totalCount != 2 {
		t.FailNow()
	}
}

// Test Ingore Panic = False
func TestNotIgnorePanic(t *testing.T) {
	resetRunner()
	resetTmp()
	OnJobReturn(hookAppendResultToTmp)
	DefaultManager.SetConfig(CronManagerConfig{false, false, 0, 0})
	DefaultManager.Start()

	entryId, _ := DefaultManager.Add("@every 2s", PanicJob{})
	// sleep 5 second, tmp length should be 1
	time.Sleep(5 * time.Second)
	if tmp[1] != "" || tmp[0] == ""{
		t.FailNow()
	}

	// status should be "STOP"
	if DefaultManager.JobMap[entryId].Status() != "STOP" {
		t.FailNow()
	}

	// successTime should be 0, totalTime should be 1
	if DefaultManager.JobMap[entryId].successCount != 0 ||
		DefaultManager.JobMap[entryId].totalCount != 1 {
		t.FailNow()
	}
}


// Test Only One = True
//func TestOnlyOne(t *testing.T) {
//	resetRunner()
//	resetTmp()
//	OnJobReturn(hookAppendResultToTmp)
//	DefaultManager.SetConfig(CronManagerConfig{false, true, 0, 0})
//	DefaultManager.Start()
//
//	entryId, _ :=DefaultManager.Add("@every 2s", TimeOutJob{})
//	// only one running even every 2s
//	time.Sleep(5500 * time.Millisecond)
//	if tmp[1] != "" || tmp[0] == ""{
//		print("Fail: Two job return")
//		t.FailNow()
//	}
//
//	// status should be "STOP"
//	if DefaultManager.JobMap[entryId].Status() != "RUNNING" {
//		print("Fail: Status should be running")
//		t.FailNow()
//	}
//
//	// successTime should be 0, totalTime should be 1
//	if DefaultManager.JobMap[entryId].successCount != 0 ||
//		DefaultManager.JobMap[entryId].totalCount != 0 {
//		t.FailNow()
//	}
//}
