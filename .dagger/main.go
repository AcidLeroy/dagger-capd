// A generated module for DaggerCapd functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dagger-capd/internal/dagger"
)

type DaggerCapd struct {
	// +private
	Socket *dagger.Socket
}

type KubernetesService struct {
	*dagger.Service
	KubeConfig *dagger.File
}

func New(ctx context.Context, socket dagger.Socket) (*DaggerCapd, error) {

	return &DaggerCapd{
		Socket: &socket,
	}, nil
}

func (d *DaggerCapd) DevContainer(ctx context.Context, dockerFile *dagger.Directory) (*dagger.Container, error) {
	k8s, err := d.newKubernetesService("capd-test")
	if err != nil {
		return nil, err
	}

	_, err = k8s.Start(ctx)
	if err != nil {
		return nil, err
	}

	// get build context with dockerfile added
	container := dockerFile.DockerBuild().
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithMountedFile("/.kube/config", k8s.KubeConfig).
		With(d.waitForKubeAPI).
		With(d.setupCAPI)

	return container, nil
}

func (m *DaggerCapd) waitForKubeAPI(ctr *dagger.Container) *dagger.Container {
	return ctr.WithExec([]string{
		"sh", "-c",
		`until kubectl get --raw /apis >/dev/null 2>&1; do
		   echo "⏳ Waiting for Kubernetes API...";
		   sleep 2;
		 done`,
	})
}

func (m *DaggerCapd) setupCAPI(ctr *dagger.Container) *dagger.Container {
	return ctr.
		WithEnvVariable("CLUSTER_NAME", "test-cluster").
		WithEnvVariable("NAMESPACE", "test-cluster").
		WithEnvVariable("CLUSTER_TOPOLOGY", "true").
		WithEnvVariable("POD_SECURITY_STANDARD_ENABLED", "false").
		WithExec([]string{"kubectl", "create", "ns", "test-cluster"}).
		WithExec([]string{"clusterctl", "init", "--infrastructure", "docker"}).
		WithExec([]string{
			"kubectl", "wait", "--for=condition=Available", "deployment", "--all",
			"-n", "capi-system", "--timeout=300s",
		}).
		WithExec([]string{
			"kubectl", "wait", "--for=condition=Available", "deployment", "--all",
			"-n", "capi-kubeadm-bootstrap-system", "--timeout=300s",
		}).
		WithExec([]string{
			"kubectl", "wait", "--for=condition=Available", "deployment", "--all",
			"-n", "capi-kubeadm-control-plane-system", "--timeout=300s",
		}).
		WithExec([]string{
			"kubectl", "wait", "--for=condition=Available", "deployment", "--all",
			"-n", "capd-system", "--timeout=300s",
		}).
		WithExec([]string{
			"sh", "-c",
			`clusterctl generate cluster ${CLUSTER_NAME} --flavor development \
			  --target-namespace ${NAMESPACE} \
			  --kubernetes-version v1.30.0 \
			  --control-plane-machine-count=3 \
			  --worker-machine-count=3 \
			  --infrastructure docker > cluster-class.yaml`,
		}).WithExec([]string{
		"sh", "-c",
		`kubectl apply -f cluster-class.yaml`,
	})
}

func (d *DaggerCapd) newKubernetesService(name string) (*KubernetesService, error) {
	k3s := dag.K3S(name)

	// Need to mount the Docker socket to allow CAPD to communicate over the docker socket
	base := k3s.Container().WithUnixSocket("/var/run/docker.sock", d.Socket)
	server := base.AsService(dagger.ContainerAsServiceOpts{
		Args: []string{
			"sh", "-c",
			"k3s server --cluster-init --bind-address $(ip route | grep src | awk '{print $NF}') --disable traefik --disable metrics-server --egress-selector-mode=disabled > /dev/null 2>&1",
		},
		InsecureRootCapabilities: true,
		UseEntrypoint:            true,
	})

	return &KubernetesService{
		Service:    server,
		KubeConfig: k3s.Config(),
	}, nil
}
