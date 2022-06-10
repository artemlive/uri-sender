# uri-sender
A simple basic application that can send screenshots of URIs by schedule 

## Usage
Here is the example config for uri-sender. 

```
{
  "notifiers": [
    {
      "type": "slack",
      "recipients": ["@artemlive"],
      "message": "<https://www.google.com/search?q=test|Google search test>",
      "slack_api_key": "",
      "cron": "*/1 * * * *",
      "screenshot": {
        "url": "https://www.google.com/search?q=test",
        "htmlElement": "",
        "wait": 10,
        "outPath": ".test/screenshots"
      }
    }
  ]
}
```

You can define *slack api_key* directly through the config or using the *SLACK AUTH TOKEN* ENV variable. If both of them are defined config-defined variable has more priority. 

Here is the example of the message generated by the above config:
<img width="852" alt="Screenshot" src="https://user-images.githubusercontent.com/3328394/173033926-4ffa78b6-e2cf-4966-abde-a20828b38aaa.png">
