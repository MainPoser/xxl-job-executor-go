package xxl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"
)

type testLogger struct{}

func (testLogger) Info(format string, a ...interface{})  {}
func (testLogger) Error(format string, a ...interface{}) {}

func newTestExecutor(serverAddr string) *executor {
	opts := newOptions(
		ServerAddr(serverAddr),
		ExecutorIp("127.0.0.1"),
		ExecutorPort("9999"),
		SetLogger(testLogger{}),
	)

	return &executor{
		opts: opts,
		regList: &taskList{
			data: make(map[string]*Task),
		},
		runList: &taskList{
			data: make(map[string]*Task),
		},
		log: opts.l,
	}
}

func TestRunTaskUsesIndependentTaskStateForConcurrentRuns(t *testing.T) {
	callbacks := make(chan int64, 2)
	callbackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/callback" {
			http.NotFound(w, r)
			return
		}

		var body call
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode callback body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body) != 1 {
			t.Errorf("callback body length = %d, want 1", len(body))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		callbacks <- body[0].LogID
		_, _ = w.Write(returnGeneral())
	}))
	defer callbackServer.Close()

	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})

	exec := newTestExecutor(callbackServer.URL)
	exec.RegTask("task.concurrent", func(ctx context.Context, param *RunReq) string {
		if param.JobID == 1 {
			close(firstStarted)
			<-releaseFirst
		}
		return param.ExecutorParams
	})

	runTask(t, exec, &RunReq{
		JobID:                 1,
		ExecutorHandler:       "task.concurrent",
		ExecutorParams:        "first",
		ExecutorBlockStrategy: serialExecution,
		LogID:                 101,
		LogDateTime:           1001,
	})

	select {
	case <-firstStarted:
	case <-time.After(time.Second):
		t.Fatal("first task did not start")
	}

	runTask(t, exec, &RunReq{
		JobID:                 2,
		ExecutorHandler:       "task.concurrent",
		ExecutorParams:        "second",
		ExecutorBlockStrategy: serialExecution,
		LogID:                 202,
		LogDateTime:           2002,
	})

	close(releaseFirst)

	got := []int64{
		readCallbackLogID(t, callbacks),
		readCallbackLogID(t, callbacks),
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })

	want := []int64{101, 202}
	if got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("callback log ids = %v, want %v", got, want)
	}
	if exec.runList.Exists("1") {
		t.Fatal("run list still contains job 1")
	}
	if exec.runList.Exists("2") {
		t.Fatal("run list still contains job 2")
	}
}

func runTask(t *testing.T, exec *executor, req *RunReq) {
	t.Helper()

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal run request: %v", err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/run", bytes.NewReader(body))

	exec.RunTask(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("run task status = %d, want %d", w.Code, http.StatusOK)
	}

	var res res
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode run task response: %v", err)
	}
	if res.Code != SuccessCode {
		t.Fatalf("run task response code = %d, want %d, body = %s", res.Code, SuccessCode, w.Body.String())
	}
}

func readCallbackLogID(t *testing.T, callbacks <-chan int64) int64 {
	t.Helper()

	select {
	case logID := <-callbacks:
		return logID
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for callback")
		return 0
	}
}
