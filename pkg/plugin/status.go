package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func Status(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var podname []string
	var showPodName bool = true
	var showPrevious bool

	clientset, err := loadConfig(kubeFlags)
	if err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			showPodName = false
		}
	}

	commonFlagList := processCommonFlags(cmd)

	podList, err := getPods(clientset, kubeFlags, podname, commonFlagList)
	if err != nil {
		return err
	}

	if cmd.Flag("previous").Value.String() == "true" {
		showPrevious = true
	}

	table := Table{}
	if !showPrevious {
		table.SetHeader(
			"T", "PODNAME", "CONTAINER", "READY", "STARTED", "RESTARTS", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE",
		)
	} else {
		table.SetHeader(
			"T", "PODNAME", "CONTAINER", "STATE", "REASON", "EXIT-CODE", "SIGNAL", "TIMESTAMP", "MESSAGE",
		)
	}

	if !showPodName {
		// we need to hide the pod name in the table
		table.HideColumn(1)
	}

	for _, pod := range podList {
		for _, container := range pod.Status.ContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statusBuildRow(container, pod.Name, "S", showPrevious)
			table.AddRow(tblOut...)
		}
		for _, container := range pod.Status.InitContainerStatuses {
			// should the container be processed
			if skipContainerName(commonFlagList, container.Name) {
				continue
			}
			tblOut := statusBuildRow(container, pod.Name, "I", showPrevious)
			table.AddRow(tblOut...)
		}
	}
	table.Print()
	return nil

}

func statusBuildRow(container v1.ContainerStatus, podName string, containerType string, showPrevious bool) []string {
	var reason string
	var exitCode string
	var signal string
	var message string
	var startedAt string
	var started string
	var strState string
	var state v1.ContainerState

	// fmt.Println("F:statusBuildRow:Name=", container.Name)

	if showPrevious {
		state = container.LastTerminationState
	} else {
		state = container.State
	}

	if state.Waiting != nil {
		strState = "Waiting"
		reason = state.Waiting.Reason
		message = state.Waiting.Message
	}

	if state.Terminated != nil {
		strState = "Terminated"
		exitCode = fmt.Sprintf("%d", state.Terminated.ExitCode)
		signal = fmt.Sprintf("%d", state.Terminated.Signal)
		startedAt = state.Terminated.StartedAt.String()
		reason = state.Terminated.Reason
		message = state.Terminated.Message
	}

	if state.Running != nil {
		strState = "Running"
		startedAt = state.Running.StartedAt.String()
	}

	if container.Started != nil {
		started = fmt.Sprintf("%t", *container.Started)
	}
	ready := fmt.Sprintf("%t", container.Ready)
	restarts := fmt.Sprintf("%d", container.RestartCount)

	if showPrevious {
		return []string{
			containerType,
			podName,
			container.Name,
			strState,
			reason,
			exitCode,
			signal,
			startedAt,
			message,
		}
	} else {
		return []string{
			containerType,
			podName,
			container.Name,
			ready,
			started,
			restarts,
			strState,
			reason,
			exitCode,
			signal,
			startedAt,
			message,
		}
	}

}
