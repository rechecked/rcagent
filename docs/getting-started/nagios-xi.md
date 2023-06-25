# Nagios XI

We've built rcagent to be able to easily integrate into Nagios XI using the ReChecked Configuration Wizard. Config Wizards in Nagios XI allow you to easily run through the steps to set up active checks for your systems.

## Downloading the Config Wizard

You can [download the latest version](https://github.com/rechecked/rcagent-nagiosxi/releases/latest/download/rcagent.zip).

There is a [GitHub rcagent-nagiosxi repo](https://github.com/rechecked/rcagent-nagiosxi) specifically for the Config Wizard, so if you have any problems with the wizard, feel free to mention them in the issues.

## Installing the Config Wizard

Once you've downloaded the wizard, you can go to your Nagios XI system, under the `Admin` tab. You can then go to `Manage Config Wizards`. On the top of the page there is a form selection box. Seelect the downloaded `rcagent.zip` file downloaded from the link above and submit with `Upload & Install`.

Now you can go to the `Configure` tab and select `Configuration Wizards` and see the ReChecked Agent configuration wizard.

## Running the Wizard

When you're on the `Configure` tab and selected the `Configuration Wizards` section. Find and click the ReChecked Agent configuration wizard.

From here you'll see step 1 is just a collection of connection data. You can fill in all the details for the agent you want to run the wizard against.

Step 2 is going to give you options for what you can configure. The ReChecked Config Wizard is considered a smart wizard, so it will grab some data from the rcagent and populate some defaults. They may not be what you want though, so be sure to check each of the defaults set before adding those checks. Just click the checkbox to add the checks you want.

You can now continue normally through the rest of the wizard and apply the configuration! You should see your new host and services appear shortly after.

## Editing Hosts and Services

If you've already ran the configuration wizard, one of the easiest ways to make changes is to re-run the wizard. Re-running the wizard will overwrite any host/service that is named the same. Hosts and services that do not already exist will be created. However, hosts and services that may not be working, and were not edited by the wizard running again, will need to be removed from the Core Config Manager (CCM) in Nagios XI.

To change just a specific Host or Service, you'll have to edit them inside the Core Config Manager (CCM) and then apply the new configuration manually.