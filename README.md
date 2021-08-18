# IPGeoLocation
## Overview
An exercise in building REST-API service that returns IP to location (Country/City) mapping.

## Usage
- To build run `./build.sh`
- To run docker `./run.sh`
- Environment variables can be set in env file
- Example usage:  
```  
curl http://localhost:8000/v1/find-country?ip=10.10.11.11
{"Country":"Italy","City":"Milano"}
```

