package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/rancher/support-bundle-kit/pkg/manager"
)

var (
	sbm = &manager.SupportBundleManager{}
)

// managerCmd represents the manager command
var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Support Bundle Kit manager",
	Long: `Support Bundle Kit manager

The manager collects following items:
- Cluster level bundle. Including resource manifests and pod logs.
- Any external bundles. e.g., Longhorn support bundle.

And it also waits for reports from support bundle agents. The reports contain:
- Logs of each node.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := sbm.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	},
}

func getEnvStringSlice(key string) []string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return []string{}
	}
	return strings.Split(value, ",")
}

func init() {
	rootCmd.AddCommand(managerCmd)
	managerCmd.PersistentFlags().StringSliceVar(&sbm.Namespaces, "namespaces", getEnvStringSlice("SUPPORT_BUNDLE_TARGET_NAMESPACES"), "List of namespaces delimited by ,")
	managerCmd.PersistentFlags().StringVar(&sbm.BundleName, "bundlename", os.Getenv("SUPPORT_BUNDLE_NAME"), "The support bundle name")
	managerCmd.PersistentFlags().StringVar(&sbm.OutputDir, "outdir", os.Getenv("SUPPORT_BUNDLE_OUTPUT_DIR"), "The directory to store the bundle")
	managerCmd.PersistentFlags().StringVar(&sbm.ManagerPodIP, "manager-pod-ip", os.Getenv("SUPPORT_BUNDLE_MANAGER_POD_IP"), "The support bundle manager's IP (pod runs this app)")
	managerCmd.PersistentFlags().StringVar(&sbm.ImageName, "image-name", os.Getenv("SUPPORT_BUNDLE_IMAGE"), "The support bundle image")
	managerCmd.PersistentFlags().StringVar(&sbm.ImagePullPolicy, "image-pull-policy", os.Getenv("SUPPORT_BUNDLE_IMAGE_PULL_POLICY"), "Pull policy of the support bundle image")
	managerCmd.PersistentFlags().StringVar(&sbm.NodeSelector, "node-selector", os.Getenv("SUPPORT_BUNDLE_NODE_SELECTOR"), "NodeSelector of agent DaemonSet. e.g., key1=value1,key2=value2")
	managerCmd.PersistentFlags().StringVar(&sbm.TaintToleration, "taint-toleration", os.Getenv("SUPPORT_BUNDLE_TAINT_TOLERATION"), "Toleration of agent DaemonSet. e.g., key1=value1:NoSchedule,key2=value2:NoSchedule")
	managerCmd.PersistentFlags().StringVar(&sbm.RegistrySecret, "registry-secret", os.Getenv("SUPPORT_BUNDLE_REGISTRY_SECRET"), "The registry secret for image pull")
	managerCmd.PersistentFlags().StringVar(&sbm.SpecifyCollector, "specify-collector", os.Getenv("SUPPORT_BUNDLE_COLLECTOR"), "Execute specify collector script. e.g., longhorn")
	managerCmd.PersistentFlags().StringSliceVar(&sbm.ExcludeResourceList, "exclude-resources", getEnvStringSlice("SUPPORT_BUNDLE_EXCLUDE_RESOURCES"), "List of resources to exclude. e.g., settings.harvesterhci.io,secrets")
	managerCmd.PersistentFlags().StringSliceVar(&sbm.BundleCollectors, "extra-collectors", getEnvStringSlice("SUPPORT_BUNDLE_EXTRA_COLLECTORS"), "Get extra resource for the specific components e.g., harvester")
	managerCmd.PersistentFlags().StringVar(&sbm.Description, "description", os.Getenv("SUPPORT_BUNDLE_DESCRIPTION"), "The support bundle description")
	managerCmd.PersistentFlags().StringVar(&sbm.IssueURL, "issue-url", os.Getenv("SUPPORT_BUNDLE_ISSUE_URL"), "The support bundle issue url")
	managerCmd.PersistentFlags().DurationVar(&sbm.NodeTimeout, "node-timeout", parseDurationString(os.Getenv("SUPPORT_BUNDLE_NODE_TIMEOUT")), "The support bundle node collection time out")
}

// parseDurationString could parse `1s` and `10m` duration string.
func parseDurationString(value string) time.Duration {
	if value == "" {
		return 0
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}

	return d
}
