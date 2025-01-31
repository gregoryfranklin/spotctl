package ofas

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"

	"github.com/spotinst/spotctl/internal/log"
	"github.com/spotinst/spotctl/internal/ocean/ofas/config"
	"github.com/spotinst/spotctl/internal/uuid"
)

const (
	spotConfigMapNamespace        = metav1.NamespaceSystem
	spotConfigMapName             = "spotinst-kubernetes-cluster-controller-config"
	clusterIdentifierConfigMapKey = "spotinst.cluster-identifier"

	pollInterval = 5 * time.Second
	pollTimeout  = 5 * time.Minute
)

func ValidateClusterContext(ctx context.Context, client kubernetes.Interface, clusterIdentifier string) error {
	cm, err := client.CoreV1().ConfigMaps(spotConfigMapNamespace).Get(ctx, spotConfigMapName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("could not get ocean configuration, %w", err)
	}

	id := cm.Data[clusterIdentifierConfigMapKey]
	if id != clusterIdentifier {
		return fmt.Errorf("current cluster identifier is %q, expected %q", id, clusterIdentifier)
	}

	return nil
}

func CreateDeployerRBAC(ctx context.Context, client kubernetes.Interface, namespace string) error {
	sa, crb, err := config.GetDeployerRBAC(namespace)
	if err != nil {
		return fmt.Errorf("could not get deployer rbac objects, %w", err)
	}

	_, err = client.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("could not create deployer service account, %w", err)
	}

	_, err = client.RbacV1().ClusterRoleBindings().Create(ctx, crb, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("could not create deployer cluster role binding, %w", err)
	}

	return nil
}

type jobValues struct {
	Name            string
	Namespace       string
	ImagePullSecret string
	ImageDeployer   string
	ImageOperator   string
	ServiceAccount  string
}

func Deploy(ctx context.Context, client kubernetes.Interface, namespace string) error {
	values := jobValues{
		Name:           fmt.Sprintf("ofas-deploy-%s", uuid.NewV4().Short()),
		Namespace:      namespace,
		ImageDeployer:  "public.ecr.aws/f4k1p1n4/bigdata-deployer:main", // TODO(thorsteinn) temporary, will be done from backend
		ServiceAccount: config.ServiceAccountName,
	}

	jobTemplate, err := template.New("deployJob").Parse(deployJobTemplate)
	if err != nil {
		return fmt.Errorf("could not parse job template, %w", err)
	}

	jobManifestBytes := new(bytes.Buffer)
	err = jobTemplate.Execute(jobManifestBytes, values)
	if err != nil {
		return fmt.Errorf("could not execute job template, %w", err)
	}

	jobManifest := jobManifestBytes.String()

	job := &batchv1.Job{}
	err = yamlutil.NewYAMLOrJSONDecoder(strings.NewReader(jobManifest), len(jobManifest)).Decode(job)
	if err != nil {
		return fmt.Errorf("could not decode job manifest, %w", err)
	}

	log.Debugf("Creating deploy job %s/%s", job.Namespace, job.Name)
	createdJob, err := client.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not create deploy job, %w", err)
	}

	jobSucceeded := false
	err = wait.Poll(pollInterval, pollTimeout, func() (bool, error) {
		job, err := client.BatchV1().Jobs(createdJob.Namespace).Get(ctx, createdJob.Name, metav1.GetOptions{})
		if err != nil {
			log.Debugf("Could not get deploy job, err: %s", err.Error())
			return false, nil
		}

		done, succeeded := checkJobStatus(job)
		if done {
			jobSucceeded = succeeded
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		return fmt.Errorf("wait for deploy job completion failed, %w", err)
	}

	if !jobSucceeded {
		return fmt.Errorf("deploy job failure")
	}

	return nil
}

func checkJobStatus(job *batchv1.Job) (done bool, succeeded bool) {
	log.Debugf("Deploy job conditions: %v", job.Status.Conditions)
	for _, condition := range job.Status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			if condition.Type == batchv1.JobComplete {
				done = true
				succeeded = true
			} else if condition.Type == batchv1.JobFailed {
				done = true
				succeeded = false
			}
		}
	}
	return done, succeeded
}

const deployJobTemplate = `apiVersion: batch/v1
kind: Job
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  ttlSecondsAfterFinished: 300
  template:
    spec: 
      containers:
        - image:
            {{.ImageDeployer}}
          name: deployer
          args:
            - install
            - --create-bootstrap-environment 
          resources: { }
          imagePullPolicy: Always
      serviceAccountName: {{.ServiceAccount}}
      restartPolicy: Never
`
