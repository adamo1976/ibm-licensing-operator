
# License Service for stand-alone IBM Containerized Software without IBM Cloud Paks

<b>Scenario: Learn how to deploy License Service on Kubernetes clusters and track license usage of the stand-alone IBM Containerized Software</b>

You can use the `ibm-licensing-operator` to install License Service on Kubernetes clusters. License Service collects information about license usage of IBM containerized products. You can retrieve license usage data through a dedicated API call and generate an audit snapshot on demand.

**Note:** License Service is integrate into IBM Cloud Pak solutions. You do not have to deploy it to clusters where IBM Cloud Pak solutions are deployed. License Service should already be there and collect usage data for all IBM containerized products that are enabled for reporting.

Use the installation scenario that is outline below to deploy License Service to a cluster with IBM stand-alone containerized software, where IBM Cloud Pak solutions are note deployed.

## Overview

License Service collects and measures the license usage of your products at the cluster level. You can retrieve this data upon request for monitoring and compliance. You can also retrieve an audit snapshot of the data that is audit evidence.

License Service
- Collects and measures the license usage of Virtual Processor Core (VPC) and Managed Virtual Server (MVS) metrics at the cluster level for the IBM containerized products and IBM CLoud Pak solution that are enabled for reporting.
- Collects and measures the license usage of IBM containerized software that are enabled for reporting.
- Refreshes the data every 5 minutes. With this frequency, you can capture changes in a dynamic cloud environment.
- Provides the API that you can use to retrieve data that outlines the highest license usage on the cluster.
- Provides the API that you can use to retrieve an audit snapshot that lists the highest license usage values for the requested period for products that are deployed on a cluster.

## Audit and compliance

License Service collects data that is required for compliance and audit purposes. With License Service, you can retrieve of an audit snapshot per cluster without any configuration.

Audit snapshot needs to be generated at least once a quarter for the last 90 days, and stored for 2 years. It needs to be stored in a location from which it could be retrieved and delievered to auditors. For legal requirements, see 
[IBM Container Licenses on Passport Advantage](https://www.ibm.com/software/passportadvantage/containerlicenses.html).

## Backup

The license usage data that is collected by License Service is stored in the cluster memory. Nonetheless, it is a good practice to generate an audit snapshot periodically for backup purposes and store it in a safe location. You do not need to perform any other backup.

Before decommissioning a cluster, you need to generate an audit snapshot to record the cluster license usage until the day of decommissioning.

## Documentation

- [Preparing for installation](Preparing_for_installation.md)
  - [Supported platforms](Preparing_for_installation.md#supported-platforms)
  - [Operator versions](Preparing_for_installation.md#operator-versions)
- [Installing License Service](Installation_scenarios.md)
    - [Automatically installing ibm-licensing-operator with a stand-alone IBM Containerized Software using Operator Lifecycle Manager (OLM)](Automatic_installation.md)
    - [Manually installing License Service on OCP 4.2+](Install_on_OCP.md)
    - [Manually installing License Service on Kubernetes from scratch with `kubectl`](Install_from_scratch.md)
    - [Offline installation](Install_offline.md)
- [Configuration](Configuration.md)
  - [Configuring ingress](Configuration.md#configuring-ingress)
  - [Checking License Service components](Configuration.md#checking-license-service-components)
  - [Using custom certificates](Configuration.md#using-custom-certificates)
  - [Cleaning existing License Service dependencies](Configuration.md#cleaning-existing-license-service-dependencies)
- [Retrieving license usage data from the cluster](Retrieving_data.md)
  - [Available APIs](Retrieving_data.md#available-apis)
  - [Tracking license usage in multicluster environment](Retrieving_data.md#tracking-license-usage-in-multicluster-environment)
- [Uninstalling License Service from a Kubernetes cluster](Uninstalling.md)
- [Troubleshooting](Troubleshooting.md)
  - [Preparing resources for offline installation without git](Troubleshooting.md#prepareing-resources-for-offline-installation-without-git)