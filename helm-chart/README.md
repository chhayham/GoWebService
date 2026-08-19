# GoWebService Helm Chart

A Helm chart for deploying the GoWebService application to Kubernetes.

## Chart Details

This chart deploys a Go web service with the following features:

- Configurable replica count and resource limits
- Health checks (liveness, readiness, startup) using the `/healthcheck` endpoint
- Horizontal Pod Autoscaling (HPA) support
- Ingress support with TLS
- Service Account management
- Customizable security contexts

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Docker image of your GoWebService application

## Installing the Chart

### Basic Installation

```bash
helm install my-release ./helm-chart
```

### Install with Custom Values

```bash
helm install my-release ./helm-chart --set replicaCount=3
```

### Install with Values File

```bash
helm install my-release ./helm-chart -f custom-values.yaml
```

### Install with Ingress Enabled

```bash
helm install my-release ./helm-chart \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=gowebservice.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

## Uninstalling the Chart

```bash
helm uninstall my-release
```

## Configuration

The following table lists the configurable parameters of the GoWebService chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Docker image repository | `ghcr.io/chhayham/gowebservice` |
| `image.tag` | Docker image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `service.type` | Kubernetes service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `256Mi` |
| `resources.requests.cpu` | CPU request | `100m` |
| `resources.requests.memory` | Memory request | `128Mi` |
| `autoscaling.enabled` | Enable HPA | `false` |
| `livenessProbe.httpGet.path` | Liveness probe path | `/healthcheck` |
| `readinessProbe.httpGet.path` | Readiness probe path | `/healthcheck` |

### Example Custom Values File

```yaml
# custom-values.yaml
replicaCount: 3

image:
  repository: ghcr.io/chhayham/gowebservice
  tag: "latest"
  pullPolicy: Always

service:
  type: LoadBalancer
  port: 80

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: gowebservice.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: gowebservice-tls
      hosts:
        - gowebservice.example.com

resources:
  limits:
    cpu: 1000m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
```

## Health Checks

This chart configures health checks using the application's built-in `/healthcheck` endpoint:

- **Liveness Probe**: Checks every 10 seconds after initial 30 second delay
- **Readiness Probe**: Checks every 5 seconds after initial 10 second delay
- **Startup Probe**: Allows up to 100 seconds for startup (10 retries × 10 seconds)

## Application Endpoints

After deployment, the following endpoints are available:

| Endpoint | Description |
|----------|-------------|
| `/` | Homepage |
| `/healthcheck` | Health check endpoint |
| `/api` | API endpoint (GET only) |

## Building the Docker Image

Before deploying, make sure to build and push your Docker image:

```bash
# Build the image
docker build -t ghcr.io/chhayham/gowebservice:latest .

# Push to your registry
docker push ghcr.io/chhayham/gowebservice:latest
```

## Directory Structure

```
helm-chart/
├── Chart.yaml           # Chart metadata
├── values.yaml          # Default configuration values
├── templates/
│   ├── _helpers.tpl     # Helper templates
│   ├── deployment.yaml  # Deployment manifest
│   ├── service.yaml     # Service manifest
│   ├── ingress.yaml     # Ingress manifest (conditional)
│   ├── hpa.yaml         # HPA manifest (conditional)
│   ├── serviceaccount.yaml  # Service account (conditional)
│   ├── namespace.yaml   # Namespace manifest
│   └── NOTES.txt        # Post-install notes
```

## Supported Kubernetes Versions

This chart is compatible with Kubernetes 1.19 and above.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see the [LICENSE](../../LICENSE) file for details.

## Author

**Chhay Ham**
