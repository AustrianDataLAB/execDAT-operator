package controllers

import (
	"context"
	"fmt"

	kbatch "k8s.io/api/batch/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckJobCompletion checks if a Job has completed
func CheckJobCompletion(job *kbatch.Job, r client.Client) (bool, error) {
	ctx := context.Background()

	err := r.Get(ctx, client.ObjectKeyFromObject(job), job)
	if err != nil {
		return false, fmt.Errorf("failed to get job: %w", err)
	}

	if job.Status.CompletionTime != nil {
		// Job has completed
		return true, nil
	}

	// Job is still running, wait and check again
	return false, nil
}
