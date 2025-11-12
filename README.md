# openmcp-testing

OpenMCP-testing helps to set up e2e test suites for openmcp applications. Like [xp-testing](https://github.com/crossplane-contrib/xp-testing) but for [openmcp](https://github.com/openmcp-project).

* [`pkg/conditions`](./pkg/conditions/) provides common pre/post condition checks
* [`pkg/providers`](./pkg/providers/) provides functionality to test cluster-providers, platform-services and service-providers
* [`pkg/resources`](./pkg/resources/) provides functionality to (batch) import and delete resources
* [`pkg/setup`](./pkg/setup/) provides functionality to bootstrap an openmcp environment

## References

* [e2e-framework](https://github.com/kubernetes-sigs/e2e-framework)
