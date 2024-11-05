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
5. Change Configurtation > Configuration window > Options > Change address
5. Start the add-on.

## Configuration

Edit the add-on configuration with your address query. For example:

```yaml
Change Configurtation : Addon > Configuration  sub window > Options > Change address
```
![ha-wast-collection-config](https://github.com/rozicdejan/ha-go-waste-collection/blob/main/pictures/ha-waste-collection-configuration.png?raw=true)

## Displaying in Home Assistant Dashboard
You can display the waste collection schedule using two types of visual elements in your dashboard:
### 1. Entities Card (Text-Based View)
To set up a basic text-based display, use the following configuration:
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
![ha-wast-collection-enteties](https://github.com/rozicdejan/ha-go-waste-collection/blob/main/pictures/ha-waste-collection-schedule-entities.png?raw=true)


### 2. Picture-Elements Card (Visual Display with Icons)
For a more visually engaging setup, you can use a picture-elements card. This configuration includes custom icons for each type of waste and dynamically displays the next pickup date beneath each icon.

Add the following to your Picture-Elements card configuration in your dashboard:





```yaml
type: picture-elements
elements:
  - type: image
    entity: sensor.waste_collection_ha
    title: Waste Collection Schedule
    image: /local/Mesani.svg
    style:
      top: 30%
      left: 20%
      width: 15%
    tap_action:
      action: more-info
  - type: state-label
    entity: sensor.waste_collection_ha
    attribute: next_mko
    style:
      top: 10%
      left: 20%
      color: white
  - type: image
    entity: sensor.waste_collection_ha
    image: /local/Embalaza.svg
    style:
      top: 30%
      left: 50%
      width: 15%
    tap_action:
      action: more-info
  - type: state-label
    entity: sensor.waste_collection_ha
    attribute: next_emb
    style:
      top: 10%
      left: 50%
      color: white
  - type: image
    entity: sensor.waste_collection_ha
    image: /local/Bioloski_odpadki.svg
    style:
      top: 30%
      left: 80%
      width: 15%
    tap_action:
      action: more-info
  - type: state-label
    entity: sensor.waste_collection_ha
    attribute: next_bio
    style:
      top: 10%
      left: 80%
      color: white
camera_view: auto
image: /api/image/serve/8ab0513c5235c08431f5f53b653b6849/512x512
dark_mode_image: /api/image/serve/1ee0e0a51a1d1c5b8d969b7c8266fd10/512x512
layout_options:
  grid_columns: 4
  grid_rows: 3
```

Note: Ensure that the images (Mesani.svg, Embalaza.svg, Bioloski_odpadki.svg) are saved in the /local/ directory in Home Assistant’s file structure.

![ha-wast-collection](https://github.com/rozicdejan/ha-go-waste-collection/blob/main/pictures/ha-waste-collection-schedule.png?raw=true)



## Troubleshooting
If you encounter issues, check the add-on logs for any errors.