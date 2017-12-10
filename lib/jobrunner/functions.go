package jobrunner

import (
	"fmt"
	"strings"

	"time"

	"encoding/json"

	"os"

	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/transport/ssh"
)

var lastJob = 0

func (r *JobRunner) Worker(id int) {
	log.Debugf("JobRunner.Worker: spawned with id %d", id)
	for job := range r.JobQueue {
		result := Result{Node: job.Node}

		if job.Node == "test" {
			log.Debugf("job.Command: %s", job.Command)
			switch job.Command {
			case TEST_STDERR:
				result.Stderr = TEST_STDERR_MSG
			case TEST_ERR:
				result.err = fmt.Errorf(TEST_ERR_MSG)
			default:
				result.Stdout = TEST_STDOUT_MSG
			}
			r.ResultQueue <- result
			continue
		}

		cmd := ""
		cmd = fmt.Sprintf("%s", job.Command)
		params := []string{"-n", job.Node, cmd}
		log.Debugf("Worker(%d): running ssh %s", id, strings.Join(params, " "))
		stdout, stderr, err := ssh.Run(params)
		if err != nil {
			errMsg := fmt.Errorf("failed to run ssh")
			log.Warningf("Worker(%d): %v: %v", id, errMsg, err)
			result.err = errMsg
			log.Debugf("Worker(%d): submitted err result: %v", id, result.err)
			r.ResultQueue <- result
			continue
		}

		result.Stdout = stdout
		result.Stderr = stderr
		result.err = err

		r.ResultQueue <- result
	}

	result := Result{
		Node: KILL_PILL,
	}
	r.ResultQueue <- result
	log.Debugf("Worker(%d): submitted kill result: %v", id, result)
}

func (r *JobRunner) StartWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go r.Worker(i)
	}

	r.numWorkers = numWorkers
}

func (r *JobRunner) GetNumWorkers() int {
	return r.numWorkers
}

func (r *JobRunner) Submit(job Job) error {
	clusterData := config.ClusterConfig{}

	for _, cluster := range r.config.Clusters {
		if cluster.Name == job.Cluster {
			clusterData = cluster
			break
		}
	}

	if job.Cluster == KILL_PILL {
		for len(r.JobQueue) > 0 {
			time.Sleep(1 * time.Second)
		}
		close(r.JobQueue)

		return nil
	} else if job.Cluster == TEST_CLUSTER {
		for range clusterData.Hosts {
			nodeJob := NodeJob{
				Node:    "test",
				Command: job.Command,
			}
			r.JobQueue <- nodeJob
			log.Debugf("JobRunner.Submit: Succesfully submitted job with id %d", lastJob)
			lastJob += 1
		}
		return nil
	}

	if clusterData.Name == "" {
		return fmt.Errorf("JobRunner.Submit: Unknown job")
	}

	jobData, err := r.config.GetCommand(job.Command)
	if err != nil {
		return fmt.Errorf("JobRunner.Submit: %v", err)
	}

	for _, node := range clusterData.Hosts {
		nodeJob := NodeJob{
			Node:    node,
			Command: jobData.Command,
		}
		r.JobQueue <- nodeJob
		log.Debugf("JobRunner.Submit: Succesfully submitted job with id %d", lastJob)
		lastJob += 1
	}

	return nil
}

func (r *JobRunner) Start(outputModifier string) {
	numKillMessages := 0

	log.Debugf("JobRunner.Start: waiting for responses")
	allResponses := []Result{}

	output := os.Stdout
	if r.TestFd != nil {
		output = r.TestFd
	}

	for response := range r.ResultQueue {
		log.Debugf("JobRunner.Start: Got response: %v", response)
		if response.Node == KILL_PILL {
			numKillMessages += 1
			if numKillMessages == r.numWorkers {
				log.Debugf("JobRunner.Start: Shutting down")
				close(r.ResultQueue)
				break
			}
			continue
		}

		switch outputModifier {
		case MOD_JSON:
			{
				allResponses = append(allResponses, response)
			}
		default:
			{
				if len(response.Stdout) > 0 {
					stdout := strings.Split(response.Stdout, "\n")
					if len(stdout) > 1 {
						stdout = stdout[:len(stdout)-1]
					}
					for _, line := range stdout {
						fmt.Fprintf(output, "%s stdout: %s\n", response.Node, line)
					}
				}

				if len(response.Stderr) > 0 {
					stderr := strings.Split(response.Stderr, "\n")
					if len(stderr) > 1 {
						stderr = stderr[:len(stderr)-1]
					}
					for _, line := range stderr {
						fmt.Fprintf(output, "%s stderr: %s\n", response.Node, line)
					}
				}

				if response.err != nil {
					fmt.Fprintf(output, "%s err: %v\n", response.Node, response.err)
				}
			}
		}
	}

	switch outputModifier {
	case MOD_JSON:
		{
			data, err := json.Marshal(allResponses)
			if err != nil {
				log.Warningf("JobRunner.Start json.Marshal: %v", err)
				return
			}

			fmt.Fprintf(output, "%s\n", string(data))
		}
	}
}

func (r *JobRunner) SubmitKillJob() {
	job := Job{
		Cluster: KILL_PILL,
	}

	r.Submit(job)
}
