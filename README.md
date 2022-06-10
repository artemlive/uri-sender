# uri-sender
A simple basic application that can send screenshots of URIs by schedule 

## Usage
Here is the example config for uri-sender. 

```{
  "notifiers": [
    {

      "type": "slack",
      "recipients": ["@test"],
      "message": "<https://example.com|hello>",
      "slack_api_key": "",
      "cron": "*/1 * * * *",
      "screenshot": {
        "url": "https://example.com",
        "htmlElement": "div#report",
        "wait": 10,
        "outPath": ".test/screenshots"
      }
    }
  ]
}
```

You can define *slack api_key* directly through the config or using the *SLACK AUTH TOKEN* ENV variable. If both of them are defined config-defined variable has more priority. 
