# Building Confidence in Kubernetes Controllers: Lessons Learned from Using E2E-Framework

Companion repository for the presentation `Building Confidence in Kubernetes Controllers: Lessons Learned from Using E2E-Framework`.

## Running from a laptop

Running with `go test` allows to run the code against a target cluster, useful while developing a new test.

```shell
kind create cluster
go test -timeout 0 -race -v ./... -args --context kind-kind
?   	github.com/maruina/kubecon-2024-eu/pkg/conditions	[no test files]
?   	github.com/maruina/kubecon-2024-eu/pkg/config	[no test files]
?   	github.com/maruina/kubecon-2024-eu/pkg/resources	[no test files]
?   	github.com/maruina/kubecon-2024-eu/pkg/resources/container	[no test files]
?   	github.com/maruina/kubecon-2024-eu/pkg/resources/pod	[no test files]
=== RUN   TestPodNative
=== RUN   TestPodNative/basic_pod_native
=== RUN   TestPodNative/basic_pod_native/pod_is_ready
=== RUN   TestPodNative/basic_pod_native/environment_variables_are_injected
--- PASS: TestPodNative (6.20s)
    --- PASS: TestPodNative/basic_pod_native (5.13s)
        --- PASS: TestPodNative/basic_pod_native/pod_is_ready (5.01s)
        --- PASS: TestPodNative/basic_pod_native/environment_variables_are_injected (0.09s)
=== RUN   TestPodDns
=== RUN   TestPodDns/dns
=== RUN   TestPodDns/dns/pod_is_ready
=== RUN   TestPodDns/dns/dns_is_resolving
--- PASS: TestPodDns (6.19s)
    --- PASS: TestPodDns/dns (5.16s)
        --- PASS: TestPodDns/dns/pod_is_ready (5.01s)
        --- PASS: TestPodDns/dns/dns_is_resolving (0.14s)
PASS
ok  	github.com/maruina/kubecon-2024-eu	14.350s
```

## Running on a cluster

Run the docker image as a cronjob on a live cluster to run the tests suite.

```shell
kubectl apply -f deploy/e2e.yaml
```

## Integrating with CI Visibility

We can generate JUnit XML test reports to integrate with [Datadog Test Visibility](https://docs.datadoghq.com/tests/)
