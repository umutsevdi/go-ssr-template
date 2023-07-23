# Go Fullstack SSR Template
A template for a web server that features server side rendering and API capabilities
using go's built-in html template engine.

## Features
- [x] Compile time page generation
- [x] Page indexing
- [x] Caching
- [x] Optional periodic cache updating
- [x] Runtime sitemap.xml generation
- [ ] Easier API routing.
- [ ] A cache mechanism to store pages with special variables. Currently
caching only stores pages that are compiled from components.

## Project Structure
Project structure is displayed below. After compiling you can deploy your
application binary with `content` directory and `config.json`.

```sh
.
├── app
│   ├── api
│   │   ├── router.go
│   │   └── /* Your API code goes here */
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── pages
│   │   ├── router.go
│   │   └── /* Your page rendering code goes here */
│   └── util
│       ├── config.go
│       └── index.go
├── config.json
├── content
│   ├── components
│   ├── pages
│   └── static
│       ├── css
│       ├── img
│       ├── js
│       ├── other
│       └── robots.txt /* Should be here */
└── Makefile
```
- **api**: Files in this directory are for responding requests from clients
after the page is loaded.
- **pages**: This directory is for inserting custom variables to the pages
such as personal data to render the final form of pages.
If a page is static or ready after initial compilation from the template engine
additional code is not required.
- **util**: This directory is for the library. Includes configuration files
parsing and caching.
- **config.json**: Configuration file for the server.

## Configuration
Configuration is handled using `config.json` file that is under the same directory
with the binary.
```json
{
    "url": "string", /* default: localhost */
    "port": 80,      /* default: 8080 */
    "page": {
        "enabled": true,
        "ttl": -1
    },
    "memory": {
        "enabled": true,
        "ttl": -1
    },
    "content": "string" /* default: content */
}
```
- **url**: The URL of the page, will be used to generate `sitemap.xml`.
- **port**: Port(default:8080)
- **page**: How often pages generated from components should be refreshed.
Setting it to false will stop periodic indexing, and only generate once.
(default: {false, 0})
- **memory**: Whether specific final pages should be cached, if so how often
they should be cached. (default: {false, 0})

`WARNING:` *This feature is not implemented yet.*
- **content**: Path to the content. If it's null, the program expects `content`
in the same directory.

Content directory must be laid out like this:
```
└── content
    ├── components
    ├── pages
    └── static
```

# Building
```sh
make
```
