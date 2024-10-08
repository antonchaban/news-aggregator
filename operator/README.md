# operator
Operator for managing the lifecycle of a custom resource in a Kubernetes cluster. 

Source resource is a custom resource that is created by the user. The operator watches for the creation/update/deletion of the custom resource and 
takes appropriate actions in news aggregator service.

HotNews resource is a custom resource that is created by the user. The operator watches for the creation/update/deletion of the custom resource and 
fetches the latest news from the news aggregator service and stores it in the HotNews resource due to specified criteria.

## Description
Source CRD is used for creating sources which will be maintained by news aggregator service.
if a user creates Source — it creates new sources in the news-aggregator;
if a user updates Source - it updates corresponding source in the news-aggregator;
if a user deletes Source - it removes the source from the news-aggregator.
Current Source statuses are displayed in the Source CRD Status field with last changes timestamp.

HotNews CRD is used for creating hot news which will be maintained by news aggregator service.
It uses such fields for defining the criteria for fetching news:
```go
type HotNewsSpec struct {
	// - Keywords: A list of keywords to filter news, must be always required.
	Keywords []string `json:"keywords"`
	// - DateStart: The start date for the news filter, can be empty.
	DateStart string `json:"date_start,omitempty"`
	// - DateEnd: The end date for the news filter, can be empty.
	DateEnd string `json:"date_end,omitempty"`
	// - Sources: All source names in the current namespace, if empty, will watch ALL available feeds. This should be names of Source resources.
	Sources []string `json:"sources,omitempty"`
	// - FeedGroups: Available sections of feeds from the 'feed-group-source' ConfigMap.
	FeedGroups []string `json:"feed_groups,omitempty"`
	// - SummaryConfig: Configuration for how the status will show the summary of observed hot news.
	SummaryConfig SummaryConfig `json:"summary_config"`
}
```
Then all the news that matches the criteria will be stored in the HotNews status field.

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 24.0.0+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.30+ cluster.

### To Deploy on the cluster
**Build and push your image:**

```sh
docker build . -f Dockerfile -t image-name
docker push image-name
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```


## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

