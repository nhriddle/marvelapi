# Marvel API

This is a Web Application that caches and returns values from the Marvel Hero Api (https://developer.marvel.com/docs)

## Requirements

To use this application the following go-lang libraries are needed.

    - github.com/gorilla/mux 
    - github.com/ilyakaznacheev/cleanenv
    - github.com/davecgh/go-spew/spew
    - github.com/go-redis/redis/v8

To install dependecies 

```bash

make install

```


## Configuration

To create a configuration file copy the config-sample.yml to config.yaml

Edit the following :

```bash

# Server configurations
server:
  host: "localhost"
  port: ":8080"

# Marvel API configurations
marvelapi:
  publickey : ""
  privatekey: ""

# Redis Config for caching
redis:
  host : "localhost"
  port : ":6379"

```

## Usage

Make sure that you have installed all the needed libraries.

To build use :

```bash
go build
```

To run use :

```bash
.marvelapi
```

## License
[MIT](https://choosealicense.com/licenses/mit/)



## Caching Plan Explanation
The developer uses Redis as cache for the Marvel hero details. As planned it keeps 2 types of keys.

- MARVELAPI:IDS : keeps an array of all the IDs of Marvel superheroes.

- MARVELAPI:<ID> : Keeps a JSON string with the details of the superhero

### Caching Schedule
- At first there will be nothing within the Redis cache, that's why on the first "/characters" enpoint call we traveres all the heroes by calling the "/v1/public/characters" endpoint its response includes a field named "total" which is then used to compute if we've completed caching all the heroes on Marvel's database.


Marvel's API can only display 100 results so we have to traverse it (Total hero count) / 100 considering the remainder (example : currently Marvel's total heroes are 1400+ so we need to call Marvel's API 15 times) Using the "offset" parameter we add 100 until we have the same number of heroes in our cache. While doing this we are also adding the hero details on the "MARVELAPI:<ID>" key. 

We expire the "MARVELAPI:IDS" 24 hours after it was created using the redis SETEX function. This reinitializes and refreshes the values for the "MARVELAPI:IDS" and adds additional heroes on the "MARVELAPI:<ID>" key. 

We don't need to expire thei records for  hero details because they will be updated when the "MARVELAPI:IDS" key is updated. This will also minimize us calling the "/v1/public/characters/<ID>" endpoint from Marvel's developer portal. This will also put us within the limit for the current "3000" request limit per month.



## Project Status
This project is just the beginning. The developer learned Go Lang while doing this project. This obviously needs sooooooooo much refactoring.



## Testing Progress
Currently needing troubleshooting for my tests because of the usage of redis. (needs confirmation if it's my problem lol!). But created a unit testing document from postmant for web calls. Just make sure that the web application is running in your localhost before testing.


[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/7c79462483588adc89ee)








 
