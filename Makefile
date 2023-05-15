build:
	podman build -t docker.io/byteplow/idd4 .

helmbuild:
	sed -i "s/gitshaplaceholder/$$(git rev-parse --short HEAD)/g" contrib/deployment/chart/values.yaml
	sed -i "s/gitshaplaceholder/$$(git rev-parse --short HEAD)/g" contrib/deployment/chart/Chart.yaml
	helm package contrib/deployment/chart
	sed -i "s/$$(git rev-parse --short HEAD)/gitshaplaceholder/g" contrib/deployment/chart/values.yaml
	sed -i "s/$$(git rev-parse --short HEAD)/gitshaplaceholder/g" contrib/deployment/chart/Chart.yaml

publisch: build helmbuild
	podman tag docker.io/byteplow/idd4 docker.io/byteplow/idd4:$$(git rev-parse --short HEAD)
	podman push docker.io/byteplow/idd4:$$(git rev-parse --short HEAD)

	helm push idd4-*-$$(git rev-parse --short HEAD).tgz oci://docker.io/byteplow