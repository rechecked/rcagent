# Manager

If you are using ReChecked Manager, you will use this section to tell the agent how to connect to the manager along with some options that tell it how to interact and behave with the manager.

Configuration for this section should be under the top level `manager` section.

## Config Options

Options with a * next to them are **required**.

### `url`

Optional URL value to the manager.

**Default:** `https://manage.rechecked.io/api`

### `apikey` *

The API key for the organization that the agent should be connected to in the manager.

### `ignoreCert`

By default the certificate for the manager is verified. If you'd like to ignore it for some reason, set this value to `true`.

**Default:** `false`

## Example Manager Config

```
manager:
	apikey: 5E51781D7FFF4E0AA757132B017AAC80
```
