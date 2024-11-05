# Waste Collection Home Assistant Add-on

This add-on retrieves waste collection dates from the Simbio website for a specified address and sends this data to Home Assistant, where it can be displayed on the dashboard.

## Features

- Automatically fetches the next pickup dates for mixed, packaging, and biological waste.
- Pushes data to Home Assistant using the REST API.
- Configurable address query (set your own address for accurate data).

## Installation

1. In Home Assistant, go to **Supervisor** > **Add-on Store**.
2. Click on the **Repositories** menu (three dots in the top-right corner).
3. Add this repository's URL to the list: `https://github.com/rozicdejan/ha-go-waste-collection`.
4. Find the **Waste Collection** add-on in the list and install it.
5. Start the add-on.

## Configuration

Edit the add-on configuration with your address query. For example:

```yaml
address: "ZAČRET 69"
```
## Displaying in Home Assistant Dashboard
Use the following entities configuration. Add card > Entities card and paste the code:
```yaml
type: entities
title: Waste Collection Schedule
entities:
  - entity: sensor.waste_collection_ha
    name: Next Mixed Waste (MKO) Pickup
    icon: mdi:trash-can
    type: attribute
    attribute: next_mko
  - entity: sensor.waste_collection_ha
    name: Next Packaging Waste (EMB) Pickup
    icon: mdi:recycle
    type: attribute
    attribute: next_emb
  - entity: sensor.waste_collection_ha
    name: Next Bio Waste (BIO) Pickup
    icon: mdi:leaf
    type: attribute
    attribute: next_bio


```
## Troubleshooting
If you encounter issues, check the add-on logs for any errors.