HealthCheck-Backend
======
The HealthCheck-Backend is a status tracker and manager that allows its users to keep track of devices (including their equipment) placed in given locations within different projects rewritten in GO.

The placed devices are capable of pinging their current status (ONLINE, OFFLINE, WARNING). Depinding on their equipment devices are also capable of sending images (jpeg format), which are only stored within the RAM of the backend.

As in right now the API supports all routes just as the original backend and has its features extended with a register system and support for DELETE operations.

The API is hosted under render.com and can be accessed through https://healthapi-u6xn.onrender.com. However because the free tier is used for hosting, the backend will go into hibernation after a certain idle time.
This causes images in the RAM to get lost and accessing the API during hibernation will take some time or outright reject the request until it is ready again.

For functionality demonstration purposes the app https://github.com/ManuelSperl/makia-ai-app was used. The modified .apk can found in the root of this repository.
