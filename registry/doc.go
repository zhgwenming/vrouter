//
// registry key hierarchy:
//
// v2/keys/_vrouter
//     ├── hosts
//     │     ├── hostname1
//     │     │     ├── active
//     │     │     ├── bridgeinfo
//     │     │     └── ifaceinfo
//     │     │
//     │     ├── ....
//     │     │
//     │     ├── hostnameN
//     │     │     ├── ...
//     │
//     ├── routes
//     │     ├── hostname1
//     │     ├── ...
//     │     ├── hostnameN
//     │     ├── ...
//     │
//     ├── services
//     │     ├── service1
//     │     ├── ...
//     │     ├── serviceN
//     │     ├── ...
//     │
//     ├── members	// for member discovery
//     │     ├── member1
//     │     ├── ...
//     │     ├── memberN
//     │     ├── ...
//     │
//     ├── heartbeats	// cluster member heartbeats
//     │     ├── member1
//     │     ├── ...
//     │     ├── memberN
//     │     ├── ...
//     │
package registry
