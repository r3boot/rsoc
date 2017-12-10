package jobrunner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestJobRunner_Worker(t *testing.T) {
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	jobRunner := NewJobRunner(log, config)

	go jobRunner.Worker(0)
	close(jobRunner.JobQueue)
	result := <-jobRunner.ResultQueue

	wantedResult := Result{
		Node: KILL_PILL,
	}

	if result != wantedResult {
		t.Errorf("Did not get a kill pill from worker")
	}

	jobRunner = NewJobRunner(log, config)
	testJob := NodeJob{
		Node:    "nonexisting.example.com",
		Command: "nonexisting",
	}

	go jobRunner.Worker(0)
	jobRunner.JobQueue <- testJob
	time.Sleep(1 * time.Second)
	close(jobRunner.JobQueue)

	result = <-jobRunner.ResultQueue
	if result.err != nil {
		if !strings.Contains(result.err.Error(), "failed to run ssh") {
			t.Errorf("Worker result.err != nil, but did not get ssh failure error")
		}
	} else {
		t.Errorf("Worker submitted invalid result, but did not get an error")
	}
}

func TestJobRunner_StartWorkers(t *testing.T) {
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	jobRunner := NewJobRunner(log, config)

	wantedWorkers := 3
	wantedResult := Result{
		Node: KILL_PILL,
	}

	jobRunner.StartWorkers(wantedWorkers)
	jobRunner.SubmitKillJob()

	foundKillPills := 0

	for i := 0; i < wantedWorkers; i++ {
		result := <-jobRunner.ResultQueue
		if result != wantedResult {
			t.Errorf("Did not get a kill pill from worker")
		}
		foundKillPills += 1
	}

	if foundKillPills != wantedWorkers {
		t.Errorf("Wanted %d kill pills, but got %d", wantedWorkers, foundKillPills)
	}
}

func TestJobRunner_GetNumWorkers(t *testing.T) {
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	jobRunner := NewJobRunner(log, config)

	wantedWorkers := 3
	wantedResult := Result{
		Node: KILL_PILL,
	}

	jobRunner.StartWorkers(wantedWorkers)
	runningWorkers := jobRunner.GetNumWorkers()
	if runningWorkers != wantedWorkers {
		t.Errorf("jobRunner.StartWorkers want %d workers, but %d started", wantedWorkers, runningWorkers)
	}

	jobRunner.SubmitKillJob()

	foundKillPills := 0

	for i := 0; i < wantedWorkers; i++ {
		result := <-jobRunner.ResultQueue
		if result != wantedResult {
			t.Errorf("Did not get a kill pill from worker")
		}
		foundKillPills += 1
	}

	if foundKillPills != wantedWorkers {
		t.Errorf("Wanted %d kill pills, but got %d", wantedWorkers, foundKillPills)
	}
}

func TestJobRunner_Submit(t *testing.T) {
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	jobRunner := NewJobRunner(log, config)

	wantedWorkers := 1
	wantedResult := Result{
		Node: KILL_PILL,
	}

	jobRunner.StartWorkers(wantedWorkers)

	testNonExistingClusterJob := Job{
		Cluster: "nonexisting",
	}

	err := jobRunner.Submit(testNonExistingClusterJob)
	if err != nil && !strings.Contains(err.Error(), "JobRunner.Submit: Unknown job") {
		fmt.Printf("err: %v\n", err)
		t.Errorf("jobrunner.Submit: submitted job with nonexisting cluster, but did not get the correct error")
	}

	testNonExistingCommandJob := Job{
		Cluster: "test",
		Command: "nonexisting",
	}
	err = jobRunner.Submit(testNonExistingCommandJob)
	if err != nil && !strings.Contains(err.Error(), "JobRunner.Submit: MainConfig.GetCommand: No such command") {
		fmt.Printf("err: %v\n", err)
		t.Errorf("jobrunner.Submit: submitted job with nonexisting command, but did not get the correct error")
	}

	testValidJob := Job{
		Cluster: "local",
		Command: "true",
	}
	err = jobRunner.Submit(testValidJob)
	if err != nil {
		t.Errorf("jobRunner.Submit: valid job but err != nil: %v", err)
	}

	for i := 0; i < wantedWorkers+1; i++ {
		_ = <-jobRunner.ResultQueue
	}

	testSshParametersJob := Job{
		Cluster: "local",
		Command: "sshtrigger",
	}
	err = jobRunner.Submit(testSshParametersJob)
	if err != nil {
		t.Error("jobRunner.Submit err != nil: %v", err)
	}

	for i := 0; i < wantedWorkers; i++ {
		response := <-jobRunner.ResultQueue
		fmt.Printf("debug: %s\n", response.err.Error())
		fmt.Printf("nil: %v\n", response.err == nil)
		if response.err != nil {
			if response.err.Error() != "failed to run ssh" {
				t.Errorf("jobRunner.Submit: Sent job with error, but got invalid error back")
			}
		} else {
			t.Errorf("jobRunner.Submit: Sent job with error, but got no error")
		}
	}

	nonExistingCommandJob := Job{
		Cluster: "local",
		Command: "nonexistingcommand",
	}

	err = jobRunner.Submit(nonExistingCommandJob)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "JobRunner.Submit: MainConfig.GetCommand: No such command ") {
			t.Errorf("jobRunner.Submit: Submitted job with nonexisting command, but got wrong error back")
		}
	} else {
		t.Errorf("jobRunner.Submit: Submitted job with nonexisting command but got no error")
	}

	jobRunner.SubmitKillJob()

	foundKillPills := 0

	for i := 0; i < wantedWorkers; i++ {
		result := <-jobRunner.ResultQueue
		if result != wantedResult {
			t.Errorf("Did not get a kill pill from worker")
		}
		foundKillPills += 1
	}

	if foundKillPills != wantedWorkers {
		t.Errorf("Wanted %d kill pills, but got %d", wantedWorkers, foundKillPills)
	}
}

func TestJobRunner_Start(t *testing.T) {
	// Test grep output
	log, config, tmpFile := CreateTestConfig(t)
	defer CleanupTestconfig(t, tmpFile)

	wantedWorkers := 1

	jobRunner := NewJobRunner(log, config)

	grepOutFile := CreateTempFile(t)
	defer CleanupTempFile(t, grepOutFile)

	grepFd, err := os.Create(grepOutFile)
	if err != nil {
		t.Fatalf("os.Open err != nil: %v", err)
	}

	jobRunner.TestFd = grepFd

	grepStdoutJob := Job{Cluster: "test", Command: TEST_STDOUT}
	grepStderrJob := Job{Cluster: "test", Command: TEST_STDERR}
	grepErrJob := Job{Cluster: "test", Command: TEST_ERR}

	jobRunner.StartWorkers(wantedWorkers)

	jobRunner.Submit(grepStdoutJob)
	jobRunner.Submit(grepStderrJob)
	jobRunner.Submit(grepErrJob)

	jobRunner.SubmitKillJob()

	jobRunner.Start(MOD_GREP)

	grepFd.Close()

	data, err := ioutil.ReadFile(grepOutFile)
	if err != nil {
		t.Fatalf("ioutil.ReadFile err != nil: %v", err)
	}

	stdoutMsg := strings.Split(TEST_STDOUT_MSG, "\n")[0]
	stderrMsg := strings.Split(TEST_STDERR_MSG, "\n")[0]
	stdoutFilter := fmt.Sprintf("%s stdout: %s", grepStdoutJob.Cluster, stdoutMsg)
	stderrFilter := fmt.Sprintf("%s stderr: %s", grepStderrJob.Cluster, stderrMsg)
	errFilter := fmt.Sprintf("%s err: %s", grepErrJob.Cluster, TEST_ERR_MSG)

	foundStdout := false
	foundStderr := false
	foundErr := false

	for _, line := range strings.Split(string(data), "\n") {
		if line == stdoutFilter {
			foundStdout = true
		} else if line == stderrFilter {
			foundStderr = true
		} else if line == errFilter {
			foundErr = true
		}
	}

	if !foundStdout {
		t.Errorf("Did not find %s in output", stdoutFilter)
	}
	if !foundStderr {
		t.Errorf("Did not find %s in output", stderrFilter)
	}
	if !foundErr {
		t.Errorf("Did not find %s in output", errFilter)
	}

	// Test json stdout output
	jobRunner = NewJobRunner(log, config)

	jsonOutFile := CreateTempFile(t)
	defer CleanupTempFile(t, jsonOutFile)

	jsonFd, err := os.Create(jsonOutFile)
	if err != nil {
		t.Fatalf("os.Open err != nil: %v", err)
	}

	jobRunner.TestFd = jsonFd

	jsonStdoutJob := Job{Cluster: "test", Command: TEST_STDOUT}

	jobRunner.StartWorkers(wantedWorkers)

	jobRunner.Submit(jsonStdoutJob)

	jobRunner.SubmitKillJob()

	jobRunner.Start(MOD_JSON)

	jsonFd.Close()

	data, err = ioutil.ReadFile(jsonOutFile)
	if err != nil {
		t.Fatalf("ioutil.ReadFile err != nil: %v", err)
	}

	response := []Result{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Fatalf("json.Unmarshal err != nil: %v", err)
	}

	for _, entry := range response {
		if entry.Stdout != "" && entry.Stdout != TEST_STDOUT_MSG {
			t.Errorf("Sent TEST_STDOUT_MSG, but received: %v", entry.Stdout)
		} else if entry.Stderr != "" && entry.Stderr != TEST_STDERR_MSG {
			t.Errorf("Sent TEST_STDERR_MSG, but received: %v", entry.Stderr)
		} else if entry.err != nil && entry.err.Error() != TEST_ERR_MSG {
			t.Errorf("Sent TEST_ERR_MSG, but received: %v", entry.err)
		}
	}
}
