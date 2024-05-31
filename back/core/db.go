package core

import (
	"assign/common"
	"assign/config"
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	_ "github.com/go-sql-driver/mysql"
)

func MakeDB() {
	var err error
	var clusterIP string

	fmt.Printf("start\n")
	makeSecret()
	fmt.Printf("makeSecret sucess\n")
	makePV()
	fmt.Printf("makePV sucess\n")
	makePVC(config.GlobalConfig.KubernetesClient)
	fmt.Printf("makePVC sucess\n")
	createDBContainer(config.GlobalConfig.KubernetesClient)
	fmt.Printf("createDBContainer sucess\n")
	clusterIP = createDBService(config.GlobalConfig.KubernetesClient)
	fmt.Printf("createDBService sucess\n")

	if err = CheckDeploymentStatus("mariadb", "default"); err != nil {
		common.StopProgram(err)
	}
	if err = CheckServiceStatus("mariadb", "default"); err != nil {
		common.StopProgram(err)
	}
	if config.GlobalConfig.DB, err = setUpDB(clusterIP); err != nil {
		common.StopProgram(err)
	}
}

func makeSecret() {
	username := os.Getenv("SUDO_USER")
	u, err := user.Lookup(username)
	if err != nil {
		log.Fatal(err)
	}
	abspath := filepath.Join(u.HomeDir, "/cloudcomputing/back/db/env.db")
	file, err := os.Open(abspath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mariadb-bdg",
		},
		StringData: data,
	}
	_, err = config.GlobalConfig.KubernetesClient.CoreV1().Secrets("default").Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

type PersistentVolumeSpec struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   metav1.ObjectMeta      `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

func makePV() {
	username := os.Getenv("SUDO_USER")
	u, err := user.Lookup(username)
	if err != nil {
		log.Fatal(err)
	}
	abspath := filepath.Join(u.HomeDir + "/cloudcomputing/back/db/pv.yaml")

	yamlFile, err := os.ReadFile(abspath)
	if err != nil {
		panic(err)
	}

	var pvSpec PersistentVolumeSpec
	if err := yaml.Unmarshal(yamlFile, &pvSpec); err != nil {
		panic(err)
	}

	capacity, ok := pvSpec.Spec["capacity"].(map[interface{}]interface{})
	if !ok {
		panic("capacity is not a map")
	}

	storage, ok := capacity["storage"].(string)
	if !ok {
		panic("storage is not a string")
	}
	hostPath, ok := pvSpec.Spec["hostPath"].(map[interface{}]interface{})
	if !ok {
		panic("hostPath is not a map")
	}

	path, ok := hostPath["path"].(string)
	if !ok {
		panic("path is not a string")
	}

	_, err = config.GlobalConfig.KubernetesClient.CoreV1().PersistentVolumes().Create(context.Background(), &v1.PersistentVolume{
		ObjectMeta: pvSpec.Metadata,
		Spec: v1.PersistentVolumeSpec{
			Capacity: v1.ResourceList{
				v1.ResourceStorage: resource.MustParse(storage),
			},
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			StorageClassName: "manual",
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: path,
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

type PersistentVolumeClaimSpec struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   metav1.ObjectMeta      `yaml:"metadata"`
	Spec       map[string]interface{} `yaml:"spec"`
}

func makePVC(clientset *kubernetes.Clientset) {
	username := os.Getenv("SUDO_USER")
	u, err := user.Lookup(username)
	if err != nil {
		log.Fatal(err)
	}
	abspath := filepath.Join(u.HomeDir, "/cloudcomputing/back/db/pvc.yaml")

	yamlFile, err := os.ReadFile(abspath)
	if err != nil {
		panic(err)
	}

	var pvcSpec PersistentVolumeClaimSpec
	if err := yaml.Unmarshal(yamlFile, &pvcSpec); err != nil {
		panic(err)
	}

	storageClassName, ok := pvcSpec.Spec["storageClassName"].(string)
	if !ok {
		panic("storage is not a string")
	}
	resources, ok := pvcSpec.Spec["resources"].(map[interface{}]interface{})
	if !ok {
		panic("capacity is not a map")
	}
	requests, ok := resources["requests"].(map[interface{}]interface{})
	if !ok {
		panic("storage is not a string")
	}
	storage, ok := requests["storage"].(string)

	_, err = clientset.CoreV1().PersistentVolumeClaims("default").Create(context.Background(), &v1.PersistentVolumeClaim{
		ObjectMeta: pvcSpec.Metadata,
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(storage),
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func createDBContainer(clientset *kubernetes.Clientset) error {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mariadb",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mariadb",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mariadb",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "mariadb",
							Image: "jitria/selfmariadb:latest",
							EnvFrom: []apiv1.EnvFromSource{
								{
									SecretRef: &apiv1.SecretEnvSource{
										LocalObjectReference: apiv1.LocalObjectReference{Name: "mariadb-bdg"},
									},
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 3306,
									Name:          "mariadb",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "mariadb-persistent-storage",
									MountPath: "/var/lib/mysql",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "mariadb-persistent-storage",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "mariadb-pv-claim",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createDBService(clientset *kubernetes.Clientset) string {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mariadb",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": "mariadb",
			},
			Ports: []v1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
			Type: v1.ServiceTypeLoadBalancer,
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		common.StopProgram(err)
	}

	serviceinfo, err := clientset.CoreV1().Services("default").Get(context.TODO(), "mariadb", metav1.GetOptions{})
	if err != nil {
		common.StopProgram(err)
	}

	return serviceinfo.Spec.ClusterIP
}

func setUpDB(ip string) (*sql.DB, error) {
	dsn := "root:qwer1234@tcp(" + ip + ":3306)/assign"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Printf("success connect db\n")
	return db, nil
}
