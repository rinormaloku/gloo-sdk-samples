package main

import (
	"context"
	"flag"
	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwclient "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	matchers "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := gwclient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	vsClient := clientset.GatewayV1().VirtualServices("gloo-system")

	virtualService, err := vsClient.Create(context.Background(), GetDummyVirtualService(), metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	println(virtualService.GetName() + " virtual service created")
}

func GetDummyVirtualService() *v1.VirtualService {
	return &v1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "gloo-system",
			Name:      "echo-vs",
		},
		Spec: api.VirtualService{
			VirtualHost: &api.VirtualHost{
				Domains: []string{"*"},
				Routes: []*api.Route{
					{
						Matchers: []*matchers.Matcher{{
							PathSpecifier: &matchers.Matcher_Prefix{
								Prefix: "/",
							},
						}},
						Action: &api.Route_RouteAction{
							RouteAction: &gloov1.RouteAction{
								Destination: &gloov1.RouteAction_Single{
									Single: &gloov1.Destination{
										DestinationType: &gloov1.Destination_Upstream{
											Upstream: &core.ResourceRef{
												Namespace: "gloo-system",
												Name:      "default-echo-8080",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
