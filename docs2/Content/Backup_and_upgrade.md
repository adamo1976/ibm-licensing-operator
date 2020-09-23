# Backup and upgrade

## Backup

The license usage data that is collected by License Service is stored in the cluster memory. Nonetheless, it is a good practice to generate an audit snapshot periodically for backup purposes and store it in a safe location. You do not need to perform any other backup.

Before decommissioning a cluster, you need to generate an audit snapshot to record the cluster license usage until the day of decommissioning.

## Upgrades

- For online environments, License Service is automatically upgraded with each new operator release.
- For offline environments, to upgrade License Service to a new version, first uninstall License Service from the cluster and redeploy it.

**Related links**
- [Go back to home page](../License_Service_main.md#documentation)
- [Retrieving license usage data from the cluster](Retrieving_data.md)
- [Uninstalling](Uninstalling.md)
- [Offline installation](Install_offline.md)