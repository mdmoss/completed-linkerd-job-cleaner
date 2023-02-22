# completed-linkerd-job-cleaner

[<img src="https://img.shields.io/docker/v/mdmoss/completed-linkerd-job-cleaner?label=Docker%20image&sort=semver">](https://hub.docker.com/repository/docker/mdmoss/completed-linkerd-job-cleaner)

Automatically clean up completed Kubernetes Jobs with lingering Linkerd proxy sidecars.

## Why does this exist?

When Pods created by Jobs and CronJobs in a Kubernetes cluster use a Linkerd service mesh, the Linkerd proxy sidecar container lives on after the main container terminates. This prevents the pod terminating, and so it hangs around in a `NotReady` state indefinitely.

There's two popular ways to handle this:
- Make sure you don't put your Jobs and CronJobs in the service mesh, or
- Update your jobs to call the `/shutdown` endpoint using [linkerd-await](https://github.com/linkerd/linkerd-await) or something similar.

This tool provides a third option: deployed to a cluster, it will periodically scan running Pods and delete any that are complete (as well as the Job that owns them).

Pods will be deleted if all their containers have terminated, other than the linkerd-proxy sidecar container.

Jobs will be deleted if all the containers in a Pod they own other than the proxy have terminated with exit code 0.

Further reading: https://linkerd.io/2.12/tasks/graceful-shutdown/#graceful-shutdown-of-job-and-cronjob-resources

## Running locally against a Kubernetes cluster

No artifacts are published other than the Docker image, so to run locally you'll need to check out the repo and build yourself.

Once you have, you can run against a cluster by providing a path to a kubeconfig file.

```
completed-linkerd-job-cleaner -kubeconfig ~/.kube/config
```

### Options available

```
$ ./completed-linkerd-job-cleaner -help
Usage of ./completed-linkerd-job-cleaner:
  -kubeconfig string
    	Path to kubeconfig file, for running externally to a cluster
  -shutdown-self
    	Post to http://localhost:4191/shutdown to shut down our own proxy when finished (default if -kubeconfig is not provided)
  -verbose
    	Enable verbose logging
```

## Deploying to Kubernetes

Running completed-linkerd-job-cleaner in a cluster requires permissions to list pods and delete jobs. There's an example in `resources.yaml`.

If the required service account is attached, the published image should work straight out of the box.

## Building

The project should build easily using Go.

```
go build
```