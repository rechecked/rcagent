# Passive Checks

Passive checks are a type of check that run on the agent side and the results are sent to the monitotring system. In our case, we normally send passive checks to Nagios via NRDP. You can use NRDP with both Nagios Core and Nagios XI which includes it pre-configured by default.

## Why Use Passive Checks?

One of the biggest reasons for using passive checks is that it will offload some of the resource requirements from the monitoring system. Active checks are ran through a subprocess by Nagios Core and will use various amounts of system resources depending on the kind of check.

Passive checks are more like a distributed system; running specific checks for the Nagios Core instance and returning just the output, result code, and perfdata back to Nagios Core. The resources required to read that data is much less than used to do the checks themselves.

!!! note

	There is one difference with passive checks compared to active checks that you should remember. You will need to set up [freshness checks](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/4/en/freshness.html) for passive checks to ensure that if you are no longer recieving a check from rcagent, you know there is a problem you need to deal with!

## Setup NRDP Sender(s)

## Adding Passive Checks