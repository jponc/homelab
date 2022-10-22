# Homelab

My homelab is composed of 3x Dell Optiplex 7050 machines with Kubernetes K3S on top of Ubuntu Server 22 operating system.


# References
- https://computingforgeeks.com/install-kubernetes-on-ubuntu-using-k3s/
- https://www.jeffgeerling.com/blog/2022/quick-hello-world-http-deployment-testing-k3s-and-traefik
- https://www.youtube.com/watch?v=BlzRx6ROiX0
- https://support.kublr.com/support/solutions/articles/33000257554-let-s-enrypt-dns-solver-with-aws-route53-service
- https://devopscube.com/setup-prometheus-monitoring-on-kubernetes/

# Contour as K3S Ingress Controller

https://github.com/projectcontour/contour/issues/4718

I've removed Traefik as the default ingress-controller to be able to use Contour (envoy).
I've installed it manually via helmchart and disabled the use of host port as it conflicts with k3s installation

```
helm install my-release bitnami/contour --namespace projectcontour --create-namespace --set envoy.useHostPort=false --kube-context=homelab
```

# Monitoring

I'm currently using https://www.statuscake.com/ to send me notifications if my healthcheck endpoint is up/down
