package main

const descriptorTemplate = `{{define "config"}}{
    "key": "hipchat-karma-cop",
    "name": "Karma Cop",
    "description": "HipChat Karma Cop",
    "vendor": {
        "name": "Eric Daugherty",
        "url": "http://www.ericdaugherty.com"
    },
    "links": {
        "self": "{{.LocalBaseUrl}}/atlassian-connect.json",
        "homepage": "{{.LocalBaseUrl}}/atlassian-connect.json"
    },
    "capabilities": {
        "hipchatApiConsumer": {
            "scopes": [
                "send_notification"
            ]
        },
        "installable": {
            "callbackUrl": "{{.LocalBaseUrl}}/installable"
        },
        "webhook": [
          {
            "url": "{{.LocalBaseUrl}}/test",
            "pattern": "^/test_cop",
            "event": "room_message",
            "name": "Test Hook"
          },
          {
            "url": "{{.LocalBaseUrl}}/ninja",
            "pattern": "^s\/.*@([^\/])+(--|\\Q++\\E)",
            "event": "room_message",
            "name": "Ninja Hook"
          }

        ]
    }
}
{{end}}`
