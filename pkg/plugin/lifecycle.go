package plugin

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var lifecycleShort = "Show lifecycle actions for each container in a named pod"

var lifecycleDescription = ` Prints lifecycle actions for individual containers. If no name is specified the
configured actions of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard, I for init and E for Ephemerial containers`

var lifecycleExample = `  # List individual container lifecycle events from pods
  %[1]s lifecycle

  # List conttainers lifecycle events from pods output in JSON format
  %[1]s lifecycle -o json

  # List lifecycle events from all containers in a single pod
  %[1]s lifecycle my-pod-4jh36

  # List lifecycle events of all containers named web-container searching all 
  # pods in the current namespace
  %[1]s lifecycle -c web-container

  # List lifecycle events of containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s lifecycle -c web-container --sort '!CONTAINER'

  # List lifecycle events of containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s lifecycle -c web-container --sort PODNAME

  # List container lifecycle events from all pods where label app equals web
  %[1]s lifecycle -l app=web

  # List lifecycle events from all containers where the pod label app is either web or mail
  %[1]s lifecycle -l "app in (web,mail)"`

type lifecycleAction struct {
	action     string
	actionName string
}

func Lifecycle(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var columnInfo containerInfomation
	// var tblHead []string
	var podname []string
	// var showPodName bool = true
	// var nodeLabels map[string]map[string]string
	// var podLabels map[string]map[string]string

	log := logger{location: "Capabilities"}
	log.Debug("Start")

	loopinfo := lifecycle{}
	builder := RowBuilder{}
	builder.LoopSpec = true
	builder.ShowPodName = true
	builder.ShowInitContainers = true

	connect := Connector{}
	if err := connect.LoadConfig(kubeFlags); err != nil {
		return err
	}

	// if a single pod is selected we dont need to show its name
	if len(args) >= 1 {
		podname = args
		if len(podname[0]) >= 1 {
			log.Debug("builder.ShowPodName = false")
			builder.ShowPodName = false
		}
	}
	commonFlagList, err := processCommonFlags(cmd)
	if err != nil {
		return err
	}
	connect.Flags = commonFlagList
	builder.CommonFlags = commonFlagList
	builder.Connection = &connect

	// podList, err := connect.GetPods(podname)
	// if err != nil {
	// 	return err
	// }

	if cmd.Flag("node-label").Value.String() != "" {
		label := cmd.Flag("node-label").Value.String()
		log.Debug("builder.LabelNodeName =", label)
		builder.LabelNodeName = label
	}

	if cmd.Flag("pod-label").Value.String() != "" {
		label := cmd.Flag("pod-label").Value.String()
		log.Debug("builder.LabelPodName =", label)
		builder.LabelPodName = label
	}

	table := Table{}
	builder.Table = &table
	columnInfo.table = &table
	builder.ShowTreeView = commonFlagList.showTreeView

	// tblHead = columnInfo.GetDefaultHead()
	// if commonFlagList.showTreeView {
	// 	// we have to control the name when displaying a tree view as the table
	// 	//  object dosent have the extra info to be able to process it
	// 	tblHead = append(tblHead, "NAME")
	// }

	// tblHead = append(tblHead, "LIFECYCLE", "HANDLER", "ACTION")
	// table.SetHeader(tblHead...)

	// if len(commonFlagList.filterList) >= 1 {
	// 	err = table.SetFilter(commonFlagList.filterList)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// commonFlagList.showPodName = showPodName
	// columnInfo.SetVisibleColumns(table, commonFlagList)

	builder.BuildRows(loopinfo)

	// for _, pod := range podList {
	// 	columnInfo.LoadFromPod(pod)

	// 	if columnInfo.labelNodeName != "" {
	// 		columnInfo.labelNodeValue = nodeLabels[pod.Spec.NodeName][columnInfo.labelNodeName]
	// 	}
	// 	if columnInfo.labelPodName != "" {
	// 		columnInfo.labelPodValue = podLabels[pod.Name][columnInfo.labelPodName]
	// 	}

	// 	columnInfo.containerType = "S"
	// 	for _, container := range pod.Spec.Containers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		// add the probes to our map (if defined) so we can loop through each
	// 		lifecycleList := buildLifecycleList(container.Lifecycle)
	// 		// loop over all probes build the output table and add the podname if multipule pods will be output
	// 		for name, action := range lifecycleList {
	// 			tblOut := lifecycleBuildRow(columnInfo, name, action)
	// 			columnInfo.ApplyRow(&table, tblOut)
	// 			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 			// table.AddRow(tblFullRow...)
	// 		}
	// 	}

	// 	columnInfo.containerType = "I"
	// 	for _, container := range pod.Spec.InitContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		lifecycleList := buildLifecycleList(container.Lifecycle)
	// 		// loop over all probes build the output table and add the podname if multipule pods will be output
	// 		for name, action := range lifecycleList {
	// 			tblOut := lifecycleBuildRow(columnInfo, name, action)
	// 			columnInfo.ApplyRow(&table, tblOut)
	// 			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 			// table.AddRow(tblFullRow...)
	// 		}
	// 	}

	// 	columnInfo.containerType = "E"
	// 	for _, container := range pod.Spec.EphemeralContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		lifecycleList := buildLifecycleList(container.Lifecycle)
	// 		// loop over all probes build the output table and add the podname if multipule pods will be output
	// 		for name, action := range lifecycleList {
	// 			tblOut := lifecycleBuildRow(columnInfo, name, action)
	// 			columnInfo.ApplyRow(&table, tblOut)
	// 			// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 			// table.AddRow(tblFullRow...)
	// 		}
	// 	}
	// }

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type lifecycle struct {
}

func (s lifecycle) Headers() []string {
	return []string{
		"LIFECYCLE", "HANDLER", "ACTION",
	}
}

func (s lifecycle) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s lifecycle) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s lifecycle) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s lifecycle) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
		NewCellText(""),
	}, nil
}

func (s lifecycle) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	lifecycleList := s.buildLifecycleList(container.Lifecycle)
	for name, action := range lifecycleList {
		out = append(out, s.lifecycleBuildRow(info, name, action))
	}
	return out, nil
}

func (s lifecycle) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := [][]Cell{}
	return out, nil

}

func (s lifecycle) lifecycleBuildRow(info BuilderInformation, handlerName string, lifecycles lifecycleAction) []Cell {

	return []Cell{
		NewCellText(handlerName),
		NewCellText(lifecycles.actionName),
		NewCellText(lifecycles.action),
	}
}

//check each type of probe and return a list
func (s lifecycle) buildLifecycleList(lifecycle *v1.Lifecycle) map[string]lifecycleAction {
	lifeCycleList := make(map[string]lifecycleAction)
	if lifecycle == nil {
		return lifeCycleList
	}

	if lifecycle.PostStart != nil {
		lifeCycleList["preStop"] = s.buildLifecycleAction(lifecycle.PostStart)
	}

	if lifecycle.PreStop != nil {
		lifeCycleList["preStop"] = s.buildLifecycleAction(lifecycle.PostStart)
	}

	return lifeCycleList
}

//given a lifecycle handler return a lifecycle action with the action translated to a string
func (s lifecycle) buildLifecycleAction(lifecycle *v1.LifecycleHandler) lifecycleAction {
	item := lifecycleAction{}

	//translate Exec action
	if lifecycle.Exec != nil {
		item.actionName = "Exec"
		item.action = strings.Join(lifecycle.Exec.Command, " ")
		return item
	}

	//translate HTTP action
	if lifecycle.HTTPGet != nil {
		item.actionName = "HTTPGet"
		actionStr := ""
		p := lifecycle.HTTPGet
		if len(p.Scheme) > 0 {
			actionStr = strings.ToLower(string(p.Scheme)) + "://"
		}

		if len(p.Host) > 0 {
			actionStr += p.Host
		}

		actionStr += portAsString(p.Port)

		if len(p.Path) > 0 {
			actionStr += p.Path
		}
		item.action = actionStr
		return item
	}

	//translate TCPSocket action
	if lifecycle.TCPSocket != nil {
		item.actionName = "TCPSocket"
		actionStr := ""
		item.action = lifecycle.TCPSocket.String()
		if len(lifecycle.TCPSocket.Host) > 0 {
			actionStr += lifecycle.TCPSocket.Host
		}
		actionStr += portAsString(lifecycle.TCPSocket.Port)
		item.action = actionStr
		return item
	}

	return lifecycleAction{}
}
