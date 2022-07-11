package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var securityShort = "Shows details of configured container security settings"

var securityDescription = ` View SecurityContext configuration that has been applied to the containers. Shows 
runAsUser and runAsGroup fields among others.
`

var securityExample = `  # List container security info from pods
  %[1]s security

  # List container security info from pods output in JSON format
  %[1]s security -o json

  # List container security info from a single pod
  %[1]s security my-pod-4jh36

  # List security info for all containers named web-container searching all 
  # pods in the current namespace
  %[1]s security -c web-container

  # List security info for all containers called web-container searching all pods in current
  # namespace sorted by container name in descending order (notice the ! charator)
  %[1]s security -c web-container --sort '!CONTAINER'

  # List security info for all containers called web-container searching all pods in current
  # namespace sorted by pod name in ascending order
  %[1]s security -c web-container --sort PODNAME

  # List container security info from all pods where label app matches web
  %[1]s security -l app=web

  # List container security info from all pods where the pod label app is either web or mail
  %[1]s security -l "app in (web,mail)"`

//list details of configured liveness readiness and startup security
func Security(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var columnInfo containerInfomation
	// var tblHead []string
	var podname []string
	// var showPodName bool = true
	// var SELinuxOptions bool
	// var nodeLabels map[string]map[string]string
	// var podLabels map[string]map[string]string

	log := logger{location: "Security"}
	log.Debug("Start")

	loopinfo := security{}
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
	if cmd.Flag("selinux").Value.String() == "true" {
		log.Debug("loopinfo.ShowSELinuxOptions = true")
		loopinfo.ShowSELinuxOptions = true
	}

	// if cmd.Flag("selinux").Value.String() == "true" {
	// 	showSELinuxOptions = true
	// 	tblHead = append(tblHead, "USER", "ROLE", "TYPE", "LEVEL")
	// } else {
	// 	tblHead = append(tblHead, "ALLOW_PRIVILEGE_ESCALATION", "PRIVILEGED", "RO_ROOT_FS", "RUN_AS_NON_ROOT", "RUN_AS_USER", "RUN_AS_GROUP")
	// }
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
	// 		tblOut := podSecurityBuildRow(pod, columnInfo, showSELinuxOptions)
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 	}

	// 	columnInfo.containerType = "S"
	// 	for _, container := range pod.Spec.Containers {
	// 		var tblOut []Cell
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		if showSELinuxOptions {
	// 			tblOut = seLinuxBuildRow(columnInfo, container.SecurityContext, pod.Spec.SecurityContext)
	// 		} else {
	// 			tblOut = securityBuildRow(columnInfo, container.SecurityContext, pod.Spec.SecurityContext)
	// 		}
	// 		columnInfo.ApplyRow(&table, tblOut)
	// 		// tblFullRow := append(columnInfo.GetDefaultCells(), tblOut...)
	// 		// table.AddRow(tblFullRow...)
	// 	}

	// 	columnInfo.containerType = "I"
	// 	for _, container := range pod.Spec.InitContainers {
	// 		var tblOut []Cell
	// 		// should the container be processed
	// 		if skipContainerName(commonFlagList, container.Name) {
	// 			continue
	// 		}
	// 		columnInfo.containerName = container.Name
	// 		if showSELinuxOptions {
	// 			tblOut = seLinuxBuildRow(columnInfo, container.SecurityContext, pod.Spec.SecurityContext)
	// 		} else {
	// 			tblOut = securityBuildRow(columnInfo, container.SecurityContext, pod.Spec.SecurityContext)
	// 		}
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

type security struct {
	ShowSELinuxOptions bool
}

func (s security) Headers() []string {
	if s.ShowSELinuxOptions {
		return []string{
			"USER",
			"ROLE",
			"TYPE",
			"LEVEL",
		}
	} else {
		return []string{
			"ALLOW_PRIVILEGE_ESCALATION",
			"PRIVILEGED",
			"RO_ROOT_FS",
			"RUN_AS_NON_ROOT",
			"RUN_AS_USER",
			"RUN_AS_GROUP",
		}
	}
}

func (s security) BuildContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s security) BuildEphemeralContainerStatus(container v1.ContainerStatus, info BuilderInformation) ([][]Cell, error) {
	return [][]Cell{}, nil
}

func (s security) HideColumns(info BuilderInformation) []int {
	return []int{}
}

func (s security) BuildPod(pod v1.Pod, info BuilderInformation) ([]Cell, error) {
	if s.ShowSELinuxOptions {
		return []Cell{
			NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
		}, nil
	} else {
		return []Cell{
			NewCellText(fmt.Sprint("Pod/", info.PodName)), //name
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
		}, nil
	}
}

func (s security) BuildContainerSpec(container v1.Container, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	if s.ShowSELinuxOptions {
		out[0] = seLinuxBuildRow(info, container.SecurityContext, info.Pod.Spec.SecurityContext)
	} else {
		out[0] = securityBuildRow(info, container.SecurityContext, info.Pod.Spec.SecurityContext)
	}
	return out, nil
}

func (s security) BuildEphemeralContainerSpec(container v1.EphemeralContainer, info BuilderInformation) ([][]Cell, error) {
	out := make([][]Cell, 1)
	if s.ShowSELinuxOptions {
		out[0] = seLinuxBuildRow(info, container.SecurityContext, info.Pod.Spec.SecurityContext)
	} else {
		out[0] = securityBuildRow(info, container.SecurityContext, info.Pod.Spec.SecurityContext)
	}
	return out, nil
}

func podSecurityBuildRow(pod v1.Pod, info containerInfomation, showSELinuxOptions bool) []Cell {
	if showSELinuxOptions {
		return []Cell{
			NewCellText(fmt.Sprint("Pod/", info.podName)), //name
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
		}
	} else {
		return []Cell{
			NewCellText(fmt.Sprint("Pod/", info.podName)), //name
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
			NewCellText(""),
		}
	}
}

func securityBuildRow(info BuilderInformation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {
	var cellList []Cell
	ape := Cell{}
	p := Cell{}
	rorfs := Cell{}
	ranr := Cell{}
	rau := Cell{}
	rag := Cell{}

	if psc != nil {
		if psc.RunAsNonRoot != nil {
			ranr = NewCellText(fmt.Sprintf("%t", *psc.RunAsNonRoot))
		}

		if psc.RunAsUser != nil {
			rau = NewCellInt(fmt.Sprintf("%d", *psc.RunAsUser), *psc.RunAsUser)
		}

		if psc.RunAsGroup != nil {
			rag = NewCellInt(fmt.Sprintf("%d", *psc.RunAsGroup), *psc.RunAsGroup)
		}
	}

	if csc != nil {
		if csc.AllowPrivilegeEscalation != nil {
			ape = NewCellText(fmt.Sprintf("%t", *csc.AllowPrivilegeEscalation))
		}

		if csc.Privileged != nil {
			p = NewCellText(fmt.Sprintf("%t", *csc.Privileged))
		}

		if csc.ReadOnlyRootFilesystem != nil {
			rorfs = NewCellText(fmt.Sprintf("%t", *csc.ReadOnlyRootFilesystem))
		}

		if csc.RunAsNonRoot != nil {
			ranr = NewCellText(fmt.Sprintf("%t", *csc.RunAsNonRoot))
		}

		if csc.RunAsUser != nil {
			rau = NewCellInt(fmt.Sprintf("%d", *csc.RunAsUser), *csc.RunAsUser)
		}

		if csc.RunAsGroup != nil {
			rag = NewCellInt(fmt.Sprintf("%d", *psc.RunAsGroup), *csc.RunAsGroup)
		}
	}

	if info.TreeView {
		cellList = info.BuildTreeCell(cellList)
	}

	cellList = append(cellList,
		ape,
		p,
		rorfs,
		ranr,
		rau,
		rag,
	)

	return cellList

}

func seLinuxBuildRow(info BuilderInformation, csc *v1.SecurityContext, psc *v1.PodSecurityContext) []Cell {
	var cellList []Cell
	seLevel := Cell{}
	seRole := Cell{}
	seType := Cell{}
	seUser := Cell{}

	if psc != nil {
		if psc.SELinuxOptions != nil {
			pselinux := psc.SELinuxOptions
			if len(pselinux.Level) > 0 {
				seLevel = NewCellText(pselinux.Level)
			}

			if len(pselinux.Role) > 0 {
				seRole = NewCellText(pselinux.Role)
			}

			if len(pselinux.Type) > 0 {
				seType = NewCellText(pselinux.Type)
			}

			if len(pselinux.User) > 0 {
				seUser = NewCellText(pselinux.User)
			}
		}
	}

	if csc != nil {
		if csc.SELinuxOptions != nil {
			cselinux := psc.SELinuxOptions
			if len(cselinux.Level) > 0 {
				seLevel = NewCellText(cselinux.Level)
			}

			if len(cselinux.Role) > 0 {
				seRole = NewCellText(cselinux.Role)
			}

			if len(cselinux.Type) > 0 {
				seType = NewCellText(cselinux.Type)
			}

			if len(cselinux.User) > 0 {
				seUser = NewCellText(cselinux.User)
			}
		}
	}

	if info.TreeView {
		cellList = info.BuildTreeCell(cellList)
	}

	cellList = append(cellList,
		seUser,
		seRole,
		seType,
		seLevel,
	)

	return cellList
}
