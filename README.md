HealthCheck-Backend
======
Bachelor Project from Heinrich Chan

The HealthCheck-Backend is a status tracker and manager that allows its users to keep track of devices (including their equipment) placed in given locations within different projects rewritten in GO.

The placed devices are capable of pinging their current status (ONLINE, OFFLINE, WARNING). Depinding on their equipment devices are also capable of sending images (jpeg/png format), which are only stored within the RAM of the backend.

As in right now the API supports all routes just as the original backend.

Roadmap (2023)
------
* (Maybe) Add an option for graceful shutdown which has been introduced with GO 1.18
* Resolve post and update functionality depending on empty or missing parameters (possible solution with optional nil)
* Implement DELETE request method for models /project, /locations (admin only), /devices, /inventory (all users)
