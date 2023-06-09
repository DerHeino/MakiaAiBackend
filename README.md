HealthCheck-Backend
======
The HealthCheck-Backend is a status tracker and manager that allows its users to keep track of devices (including their equipment) placed in given locations within different projects rewritten in GO.

The placed devices are capable of pinging their current status (ONLINE, OFFLINE, WARNING). Depinding on their equipment devices are also capable of sending images (jpeg format), which are only stored within the RAM of the backend.

As in right now the API supports all routes just as the original backend and has its features extended with a register system and support for DELETE operations.

For functionality demonstration purposes the app https://github.com/ManuelSperl/makia-ai-app was used. The modified .apk can found under ./_frontend in main-branch
