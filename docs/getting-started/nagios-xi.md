# Nagios XI

We've built rcagent to be able to easily integrate into Nagios XI using the ReChecked Configuration Wizard. Config Wizards in Nagios XI allow you to easily run through the steps to set up active checks for your systems.

## Installing the Config Wizard


## Running the Wizard


## Editing Hosts and Services

If you've already ran the configuration wizard, one of the easiest ways to make changes is to re-run the wizard. Re-running the wizard will overwrite any host/service that is named the same. Hosts and services that do not already exist will be created. However, hosts and services that may not be working, and were not edited by the wizard running again, will need to be removed from the Core Config Manager (CCM) in Nagios XI.

To change just a specific Host or Service, you'll have to edit them inside the Core Config Manager (CCM) and then apply the new configuration manually.