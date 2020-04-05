# hue-controller
Home-automation project: Middleware for controlling hue lights. This code interacts with the Philips hue API and exposes APIs for the rest of the ecosystem to use

`hue-controller/api`        : API that is exposed to the user 
`hue-controller/config`     : Obtain and update the local config from the config manager service
`hue-controller/controller` : Controller that talks to the Hue Bridge API
`hue-controller/devices`    : Package to hande specifics with Hue Lights

## API Documentation

Responses and Requests will be JSON encoded

### List all active controllers

**Defintion**

`GET /api/controllers`

**Response**

- `500 Internal Server Error` 
- `404 Not Found` if the controllers not set yet
- `200, OK` on Success
```json
[
    {
        "connection": {
            "bridge": "" `string`,
            "username": "" `string`
        },
        "lights": [
            {
                "state": {
                    "On": `bool`,
                    "Brightness": `uint8`
                },
                "id": "" `string`
            },
            .
            .
            .
        ]
    },
    {
        .
        .
        .
    },
    .
    .
    .
]
```

### List the lights of a given controller

**Definition**

`GET /api/{bridge: string}/lights`
variable `bridge` is the human name of the HueBridge

**Responses**

- `400, Bad Request` on Bridge not found
- `200, OK` on success
```json
[
    {
        "state": {
            "On": `bool`,
            "Brightness": `uint8`
        },
        "id": "" `string`
    },
    .
    .
    .
]
```

### Get the details of a light

**Definition**

`GET /api/{bridge: string}/lights/{light_id: string}`
variable `bridge` is the human name of the HueBridge
varaiable `light_id` is the id of the light associated with the HueBridge

**Responses**

- `400, Bad Request` on Bridge/Light not found
- `200, OK` on success
```json
{
    "state": {
        "On": `bool`,
        "Brightness": `uint8`
    },
    "id": "" `string`
}
```

## ToDO
- http Get with request body?
