# FileCollector

## Usage

```
Usage of ./filecollector:
  -c string
        config file path (default "filecollector.json")
```

## Configuration

example:

```json
{
  "host": "127.0.0.1",
  "port": 3000,
  "storage": "./files",
  "title": "CSS:APP Lab 1",
  "inputs": [
    {
      "name": "name",
      "label": "Name"
    },
    {
      "name": "number",
      "label": "Student Number"
    }
  ],
  "filename": "{{number}}-{{name}}"
}
```