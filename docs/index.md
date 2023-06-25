# Introduction

This is the technical documentation for rcagent. Here you'll find information about setup, configuration, and management of ReChecked Agent (rcagent).

## What is ReChecked Agent?

ReChecked Agent is a light-weight, open source monitoring agent that is compatible with Nagios-stlye monitoring systems. We built rcagent with Nagios XI and Nagios Core in mind, but other systems with similar check structure should be able to use rcagent as well.

You install rcagent on the system you want to monitor or on a system you will use to monitor other things and return them back to the main monitoring application. You can use both active and passive checks with rcagent.

## Why Another Agent?

Some people might be wondering why this agent exits. The main reason is to create a more usable, stable, open source agent that can be used by anyone. There are also a lot of features that agents like NCPA and NSClient++ do not include that could improve how agents work.

One thing that is different is that rcagent is headless. The agent itself does not come with a user interface, instead you can use the open source [ReChecked Viewer](https://view.rechecked.io) to see an agents status and test check configurations.

There is also the fact that distributed systems are common, and this agent is built to be light weight and provide a way for users to distribute the monitoring load from the monitoring system itself onto other systems as well. We have plans to build out the way we manage agents with ReChecked Manager as well.

## Why Choose Go?

There was some back and forth on deciding what language to use, however we ultimately chose to use Go for rcagent for three main reasons:

- Go is a compiled language, soy should be easy to manage different builds on different systems. A problem with other agents is how complicated the build process is.
- Go is popular and easy to get into. Unlike C or even Rust, Go is much easier to handle for newer users. We want to create an agent that can have community interaction and is easy for other users to get into the code and provide PRs.
- Go is built with concurrency in mind. It is incredibly easy to build solutions that require multiple concurrent functions running at once.

## How to Get Support

You can use the [discussions section](https://github.com/rechecked/rcagent/discussions) or [create an issue](https://github.com/rechecked/rcagent/issues) on GitHub if you find a bug.

If you need quick response one on one technical support, you can get an email based support package from [rechecked.io](https://rechecked.io) by emailing connect@rechecked.io.

## License

[We are using the standard GPL-v3 license](https://github.com/rechecked/rcagent/blob/main/LICENSE), so the agent is completely open source. You can feel comfortable contributing to the agent and using it in any environment.
