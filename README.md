faas-rancher

===========

![OpenFaaS on Rancher 1.x](https://pbs.twimg.com/media/DI-IU-1UIAACfYe.png)

This is a plugin to enable Rancher 1.x as an [OpenFaaS](https://www.openfaas.com/) backend.

[OpenFaaS](https://www.openfaas.com/) is an event-driven serverless framework for containers. Any container for Windows or Linux can be leveraged as a serverless function. FaaS is quick and easy to deploy (less than 60 secs) and lets you avoid writing boiler-plate code.

If you'd like to know more about the OpenFaaS project head over to - https://www.openfaas.com/

The code in this repository is a daemon or micro-service which can provide the basic functionality the FaaS Gateway requires:

* List functions
* Deploy function
* Delete function
* Invoke function synchronously

Any other metrics or UI components will be maintained separately in the main FaaS project.

### QuickStart

For now, [this blog post](https://medium.com/cloud-academy-inc/openfaas-on-rancher-684650cc078e) shows how you can deploy OpenFaaS on Rancher via the Catalog.

### Status

This provider targets Rancher 1.x. Since Rancher 1.x [is being deprecated](https://rancher.com/docs/rancher/v2.x/en/faq/) this repository is now in maintenance mode. Please see the [OpenFaaS provider for Kubernetes](https://github.com/openfaas/faas-netes) which works with Rancher 2.x.
