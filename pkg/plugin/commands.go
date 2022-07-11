package plugin

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var commandsShort = "Retrieves the command line and any arguments specified at the container level"

var commandsDescription = ` Prints command and arguments used to start each container (if specified), single pods and 
containers can be selected by name.  If no name is specified the container commands of all pods
in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers`

var commandsExample = `  # List containers command info from pods
  %[1]s command

  # List container command info from pods output in JSON format
  %[1]s command -o json

  # List container command info from a single pod
  %[1]s command my-pod-4jh36

  # List command info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s command -c web-container

  # List command info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s command -c web-container --sort '!CONTAINER'

  # List command info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s command -c web-container --sort PODNAME

  # List container command info from all pods where label app matches web
  %[1]s command -l app=web

  # List container command info from all pods where the pod label app is either web or mail
  %[1]s command -l "app in (web,mail)"`

type commandLine struct {
	cmd  []string
	args []string
}

func Commands(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var columnInfo containerInfomation
	// var tblHead []string
	var podname []string
	// var showPodName bool = true
	// var nodeLabels map[string]map[string]string
	// var podLabels map[string]map[string]string

	log := logger{location: "Commands"}
	log.Debug("Start")

	loopinfo := commands{}
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
	// columnInfo.treeView = commonFlagList.showTreeView
	builder.ShowTreeView = commonFlagList.showTreeView

	// tblHead = columnInfo.GetDefaultHead()
	// if commonFlagList.showTreeView {
	// 	// we have to control the name when displaying a tree view as the table
	// 	//  object dosent have the extra info to be able to process it
	// 	tblHead = append(tblHead, "NAME")
	// }

	// tblHead = append(tblHead, "COMMAND", "ARGUMENTS")
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

	// 	//do we need to show the pod line: Pod/foo-6f67dcc579-znb55
	// 	if columnInfo.treeView {
	// 		tblOut := podCommandsBuildRow(pod, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 	}

	// 	columnInfo.containerType = "S"
	// 	for _, container := range pod.Spec.Containers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		cmdLine := commandLine{
	// 			cmd:  container.Command,
	// 			args: container.Args,
	// 		}
	// 		tblOut := commandsBuildRow(cmdLine, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}

	// 	columnInfo.containerType = "I"
	// 	for _, container := range pod.Spec.InitContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		cmdLine := commandLine{
	// 			cmd:  container.Command,
	// 			args: container.Args,
	// 		}
	// 		tblOut := commandsBuildRow(cmdLine, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}

	// 	columnInfo.containerType = "E"
	// 	for _, container := range pod.Spec.EphemeralContainers {
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		cmdLine := commandLine{
	// 			cmd:  container.Command,
	// 			args: container.Args,
	// 		}
	// 		tblOut := commandsBuildRow(cmdLine, columnInfo)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}
	// }

	if err := table.SortByNames(commonFlagList.sortList...); err != nil {
		return err
	}

	outputTableAs(table, commonFlagList.outputAs)
	return nil

}

type commands struct {
}

func (s commands) Headers() []string {
	return []string{
		// "USED", "REQUEST", "LIMIT", "%REQ", "%LIMIT",
	}
}

func (s commands) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s commands) HideColumns(info BuilderInformation) []int {
	return []int{}
}

// func podStatsProcessBuildRow(pod v1.Pod, info containerInfomation) []Cell {
func (s commands) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	return []Cell{
		NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
		NewCellText(""),
		NewCellText(""),
	}, nil
}

func (s commands) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	cmdLine := commandLine{
		cmd:  container.Command,
		args: container.Args,
	}
	out := make([][]Cell, 1)
	out[0] = s.commandsBuildRow(cmdLine, info)
	return out, nil
}

func (s commands) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	cmdLine := commandLine{
		cmd:  container.Command,
		args: container.Args,
	}
	out := make([][]Cell, 1)
	out[0] = s.commandsBuildRow(cmdLine, info)
	return out, nil
}

func (s commands) commandsBuildRow(cmdLine commandLine, info BuilderInformation) []Cell {
	var cellList []Cell

	if info.TreeView {
		cellList = info.BuildTreeCell(cellList)
	}

	cellList = append(cellList,
		NewCellText(strings.Join(cmdLine.cmd, " ")),
		NewCellText(strings.Join(cmdLine.args, " ")),
	)

	return cellList
}
