Smart Task Orchestrated 

I am building an application, which is main goal is build task orchestration. It is kind of sand box environment for the task. Main I idea is user will push - code/ script, data, or external connection to box it will execute is safely in system and return output to connection or to the place where it was triggered.Backend - go, gin, mongo, Kafka, docker
Frontend - ReactJs

Application go through -

1. User visit to our website -
    1. He will see the a landing page showing motive and description and a great marking UI will be shown.
    2. There will a navbar, having 3 things. Docs, Sign Up and Sign In
    3. There will a footer section where contacts will appearing. Contain my linkedIn, phone number and email.© 2026 Pritam Kumar Maurya. Crafted with ☕ and occasional debugging tears 🐛Powered by - one line about tech used.
2. User logged in and landed to Dashboard page.
    1. There is a navbar having docs, log out
    2. There is left side toggled menu list. 
        1. Dashboard(default landing when sign in happens. If user is logged in he should be cashed to page)
        2. Team management (This is for organisational thing once we build the MVP will implement it. Basically there will be organisation selling of sand box we have to manage team)
        3. Configuration page. (Create a schedule job or one time job) - this I am not sure as of now need to brainstorm 
        4. Profile - User detail, reset password, api key generation securely.
    3. Footer exist same
3. Dashboard
    1. There are multiple switchable tabs in dashboard main monitoring purposes
    2. What are the process running as history uses of the application. 
    3. Ui should have some visual show also show that interpretation is good

Backend Go-through

1. User login with userid - password. (If userId is already taken then through error) 
2. Reset is for old and new password should be given to reset the password.
3. User basic details needs to be captured. Name email phone number editable without any verification 
4. Dashboard we have to make it have very holistic approach to monitor there are some user in starting we know that could be useful like running of the job there status history and all but in future it might be required to have different kind of monitoring so make it extendable
5. There are three kinds of the jobs this execute. I believe based on this 3 we can handle most of the situation.
    1. Simple one time job - User has to docker container / server where codebase is available it trigger the job in the code base with some shell commands.
    2. Simple Cron time job - User has to docker container / server where codebase is available it trigger the job in the code base with some shell commands.
    3. Sandbox - User configures the docker image to configure container from UI and give min-max value. Min is minimum number of container he can always available and max it could scale to max number of docker containers. It will help to handle the load assume there is min - 100 max - 1000. If load is 50 req per sec and taking 5 sec to execute it can handle. Use case is AI sandbox running, leetcode codecef like application. It should support file / json /text injection and output data in file or json text output.
    4. These there should be exposable though api though users token.
    5. There should be visibility in ui for every execution capturing logs console or application logs to have clear visible. Data should be visual showing ui is good.

Please first write the architecture plan first in md then implement it, this should be high scalable and write code in industry Lebel code . 
Please try to use go most famous libraries so that interview i can explain better show up that i have multiple library experiences. If not then write it by your self. Please delete old docs or codes not not needed. Product is open source they have to run in there environment so make is such that scaling for them will be easy and secure. But I will host so that students can use it in general so scaling for me also should work correctly.