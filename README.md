# File Collector

## Usage

```
Usage of ./filecollector:
  -c string
        config file path (default "filecollector.json")
  -v    show version
```

## Configuration

example:

```json
{
  "host": "127.0.0.1",
  "port": 3000,
  "title": "CSS:APP Homeworks Upload",
  "forms": [
    {
      "prefix": "HW01",
      "storage": "./files/hw_01",
      "title": "CSS:APP Lab 1",
      "inputs": [
        {
          "name": "name",
          "label": "Name",
          "pattern": "[A-Z]+"
        },
        {
          "name": "number",
          "label": "Student Number"
        }
      ],
      "filenameTemplate": "hw_01_{{number}}-{{name}}"
    },
    {
      "prefix": "HW02",
      "storage": "./files/hw_02",
      "title": "CSS:APP Lab 2",
      "inputs": [
        {
          "name": "name",
          "label": "Name",
          "pattern": "[A-Z]+"
        },
        {
          "name": "number",
          "label": "Student Number"
        }
      ],
      "filenameTemplate": "hw_02_{{number}}-{{name}}"
    }
  ]
}
```