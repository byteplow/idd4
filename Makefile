version=$(shell git describe --tags | sed 's/-/-alpha/' | sed 's/-/+/2')$(shell test -z "$$(git status --porcelain)" || echo -dirty)
imageVersion=$(shell echo $(version) | cut -c2-)

build:
	podman build -t docker.io/byteplow/idd4 .

package:
	sed -i "s/imageVersionPlaceholder/$(imageVersion)/g" contrib/deployment/chart/values.yaml
	sed -i "s/v0.0.0-versionPlaceholder/$(version)/g" contrib/deployment/chart/Chart.yaml
	helm package contrib/deployment/chart
	sed -i "s/$(imageVersion)/imageVersionPlaceholder/g" contrib/deployment/chart/values.yaml
	sed -i "s/$(version)/v0.0.0-versionPlaceholder/g" contrib/deployment/chart/Chart.yaml

publish: package build
	podman push docker.io/byteplow/idd4
	podman tag docker.io/byteplow/idd4 docker.io/byteplow/idd4:$(imageVersion)
	podman push docker.io/byteplow/idd4:$(imageVersion)

	helm push idd4-$(version).tgz oci://docker.io/byteplow

install:
	helm upgrade --install --set ui.image.tag=latest -f values.yaml $(shell cat releasename) ./contrib/deployment/chart/