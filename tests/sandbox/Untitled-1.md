### Requirements

1. Create a Dockerfile for a test sandbox environment.
2. The container should be a Python-based image with some initial common dependencies pre-installed (for testing purposes).
3. The sandbox environment must allow installation of additional dependencies at runtime.

---

### Sandbox Execution Flow

4. Define a data format for execution:

   * Check if required dependencies are available.
   * If not, install them first.
   * Then write the file.
   * Execute the file.
   * Return the response.

5. This format will be used when AI-generated code is executed inside the sandbox.

---

### Deployment & Scaling

6. The Docker image will be registered via the UI.

7. From the UI, set:

   * Minimum containers: 1
   * Maximum containers: 3

8. On deployment confirmation:

   * Containers should be started.
   * Health checks should run continuously (1 min interval).
   * If a container fails, maintain at least the minimum number of running containers.

9. For security:

   * The system should not function if the token is missing or misconfigured.

---

### Monitoring & Dashboard

10. The UI dashboard should display:

* All configurations
* Running containers for each configuration

11. On selecting a configuration:

* Show configuration details
* Show currently running containers

12. The monitoring system should track:

* Number of active containers
* Number of busy vs idle containers
* Resource usage (for every)

---

### Scaling Behavior

13. If all containers are busy:

* Start a new container (if below max limit)
* Otherwise, return "resource not available"

---

### Updates & Lifecycle

14. When updating a configuration:

* New containers should be added without affecting running ones

15. Provide an option to:

* Forcefully stop running containers

16. Image update rules:

* Stop all running containers first
* Ui should have provision of multi select and select all and do force stop or view detail (resource used by container)
* Then update the image
* Restart with the minimum required containers

---

### Naming Convention

17. Images should follow a pattern including a `config_id` for easy identification using regex. While developing make sure that this is exlucise of the other two times one time and sheduled. Configuration, data management, resouce mangament, monitory should afftect sandbox. so develope or design data and code in such way that they are independet of it  and can run in any enviroment saparatly. and should not have any conflict with other two types (for this if u need to changes in current code base it fine but make sure that it is woking end to end). Once exclusivity achived then focus on developement of sand box letter will see the two. Follow the UI coding accoding to current UI code style.


---

### Constraints & Details

18. Use only shell script. Use some file transfer protocal you can not push 50 mb csv via a string. Either it is csv, xlsx, json and python file push direcly in file way not in string or object. if simple object is there then we can pass it but complext files and thing we can not pass there. check that also
19. Token: `68708e8e-aff3-4428-b1ce-2113ab247748`
20. Email: `h@gmail.com`
