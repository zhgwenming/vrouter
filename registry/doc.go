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
//     ├── members
//     │     ├── member1
//     │     ├── ...
//     │     ├── memberN
//     │     ├── ...
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

package registry
