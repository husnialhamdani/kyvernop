/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup all resources created",
	Long:  `Cleanup all resources created`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cleanup called")
		cleanup(10, "default")
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cleanup(size int, namespace string) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	log.Print("Cleaning up resources...")
	for i := size - 1; i >= 0; i-- {
		counter := strconv.Itoa(i)
		objects.deleteNamespace(clientset, counter)
		objects.deleteDeployment(clientset, counter, namespace)
		objects.deleteConfigmap(clientset, counter, namespace)
		objects.deletePod(clientset, counter, namespace)
		objects.deleteSecret(clientset, counter, namespace)
		objects.deleteCronjob(clientset, counter, namespace)
	}
}
