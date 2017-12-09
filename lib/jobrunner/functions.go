package jobrunner

import (
	"fmt"
	"strings"

	"time"

	"encoding/json"

	"github.com/r3boot/rsoc/lib/config"
	"github.com/r3boot/rsoc/lib/transport/ssh"
)

var lastJob = 0

func (r *JobRunner) Worker(id int) {
	log.Debugf("JobRunner.Worker: spawned with id %d", id)
	for job := range r.JobQueue {
		cmd := ""
		cmd = fmt.Sprintf("%s", job.Command)
		params := []string{"-n", job.Node, cmd}
		stdout, stderr, err := ssh.Run(params)
		if err != nil {
			log.Warningf("JobRunner.Worker(%d): failed to run ssh: %v", id, err)
			continue
		}

		result := Result{
			Node:   job.Node,
			Stdout: stdout,
			Stderr: stderr,
			err:    err,
		}
		r.ResultQueue <- result
	}

	result := Result{
		Node: KILL_PILL,
	}
	r.ResultQueue <- result
}

func (r *JobRunner) StartWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go r.Worker(i)
	}

	r.numWorkers = numWorkers
}

func (r *JobRunner) Submit(job Job) error {
	if job.Cluster == KILL_PILL {
		for len(r.JobQueue) > 0 {
			time.Sleep(1 * time.Second)
		}
		close(r.JobQueue)

		return nil
	}

	clusterData := config.ClusterConfig{}

	for _, cluster := range r.config.Clusters {
		if cluster.Name == job.Cluster {
			clusterData = cluster
			break
		}
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

	for response := range r.ResultQueue {
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
					for _, line := range strings.Split(response.Stdout, "\n") {
						fmt.Printf("%s stdout: %s\n", response.Node, line)
					}
				}

				if len(response.Stderr) > 0 {
					for _, line := range strings.Split(response.Stderr, "\n") {
						fmt.Printf("%s stderr: %s\n", response.Node, line)
					}
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

			fmt.Printf("%s\n", string(data))
		}
	}
}

func (r *JobRunner) SubmitKillJob() {
	job := Job{
		Cluster: KILL_PILL,
	}

	r.Submit(job)
}
