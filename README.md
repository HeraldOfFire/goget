# goget
A very simple multi-threaded file downloader in Go 🤖

## Usage
Edit config.json file to define synchronous download groups.  
Each group defines a URL template (the static part of the URL) and an array of variables.  
goget makes asynchronous requests for each variable in a group.  
If a variable contains "->" it's treated as an integer range of variables.

### Example config.json file
```json
{
    "basePath": "downloaded/",
    "groups": [
        {
            "path": "GenericSeries",
            "format": ".mp4",
            "urlTemplate": "https://site.com/downloads/GenericSeries/GenericSeries_Ep_<<variable>>_SUB_ITA.mp4",
            "urlVariables": [
                "01",
                "02",
                "03",
                "04",
                "05",
                "06",
                "07",
                "08",
                "09",
                "10->20"
            ]
        }
    ]
}
```
