/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/husnialhamdani/kyvernop/objects"
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Run loads of resources creation",
	Long:  `Run loadsd of resources creation`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("execute called")
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
		config, err := kubeconfig.ClientConfig()
		if err != nil {
			panic(err)
		}
		clientset := kubernetes.NewForConfigOrDie(config)

		scales, _ := cmd.Flags().GetString("scales")
		quantityMap := map[string]int{"xs": 100, "small": 500, "medium": 1000, "large": 2000, "xl": 3000}
		size := quantityMap[scales] / 5

		//Get usage
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go getMetrics(wg, *clientset, 30, 10, "kyverno")

		//dependencies
		label := map[string]string{"app": "web"}
		namespace := "default"

		for i := 0; i < size; i++ {
			counter := strconv.Itoa(i)
			objects.CreateNamespace(*clientset, counter)
			objects.CreateDeployment(*clientset, counter, namespace, "nginx:latest", label)
			objects.CreateConfigmap(*clientset, counter, namespace)
			objects.CreatePod(*clientset, counter, namespace, "nginx")
			objects.CreateSecret(*clientset, counter, namespace)
		}

		time.Sleep(time.Duration(10) * time.Minute)

		//Delete resources - steps down
		fmt.Println("Deleting resource..")
		for i := size - 1; i >= size/2; i-- {
			counter := strconv.Itoa(i)
			objects.DeleteNamespace(*clientset, counter)
			objects.DeleteDeployment(*clientset, counter, namespace)
			objects.DeleteConfigmap(*clientset, counter, namespace)
			objects.DeletePod(*clientset, counter, namespace)
			objects.DeleteSecret(*clientset, counter, namespace)
			objects.DeleteCronjob(*clientset, counter, namespace)
		}

		wg.Wait()
		visualizeAnomaly()

		sendReport(os.Getenv("EMAILFROM"), os.Getenv("EMAILTO"), "Kyverno Automation Performance Testing report")
	},
}

func init() {
	executeCmd.Flags().StringP("scales", "s", "xs", "choose the scale size (small/medium/large/xl) default: xs")
	rootCmd.AddCommand(executeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getMetrics(wg *sync.WaitGroup, clientset kubernetes.Clientset, duration int, interval int, namespace string) {
	defer wg.Done()
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	var memoryUsage [][]int
	for len(memoryUsage) < (int(duration) * 60 / interval) {

		kyvernoPod, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/instance=kyverno"})
		if err != nil {
			fmt.Println(err)
		}
		name := kyvernoPod.Items[0].GetName()
		podmetricGet, err := mc.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), name, metav1.GetOptions{})

		if err != nil {
			fmt.Println(err)
		}
		memQuantity, ok := podmetricGet.Containers[0].Usage.Memory().AsInt64()
		if !ok {
			fmt.Println(!ok)
		}
		memoryUsage = append(memoryUsage, []int{len(memoryUsage), int(memQuantity) / 1000000})
		fmt.Println(memoryUsage)
		time.Sleep(time.Duration(interval) * time.Second)
	}

	csvfile, err := os.Create("usage.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvfile)
	for _, row := range memoryUsage {
		st := strings.Fields(strings.Trim(fmt.Sprint(row), "[]"))
		_ = csvwriter.Write(st)
	}
	csvwriter.Flush()
	csvfile.Close()
}

func visualizeAnomaly() {
	cmd := exec.Command("python3", "anomalydetection.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	cmd.Wait()
	fmt.Println("report generated in report.png")
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func sendReport(from string, to string, subject string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", "Kyverno Automation Performance Testing result:")
	m.Attach("report.png")

	d := gomail.NewDialer("smtp.gmail.com", 587, from, os.Getenv("EMAILPASS"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	} else {
		fmt.Println("email sent")
	}
}
