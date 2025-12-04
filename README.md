## What should you do to be able to obtain a chart with memory?
a) Add to your watch side app so that your watch sends memory information to your app side
 getPerformance('memory')
 setup this.call 
 (Example below)
 
b) On the side-app side, send data to the application that will collect data
c) Turn on the server that will collect data and display the data by visiting the website in a browser
http://localhost:8081/<YOUR_APP_NAME>


## How to add Data

To add data, you need to send it to the server

API Method: POST
URL: http://localhost:8081/add-data/<YOUR_APP_NAME>
Body: Data as JSON
memory object is 1:1 from  getPerformance('memory', 'perf');
This can be done by adding the field in the app-side:

==== app-service/index.js=====
example:

```javascript
    onCall(data) {
      if (data.method == 'MEMORY') {
        const dataMemory = JSON.parse(data.data);
        dataMemory.description = new Date().toLocaleTimeString();
        fetch({
          url: 'http://localhost:8081/add-data/my_app',
          method: 'POST',
          body: JSON.stringify(dataMemory),
        });
      }
    },
```

example data:

```json
{
    "memory": {
        "system": {
            "used": 1546320,
            "total": 3145728
        },
        "framework": {
            "used": 1036308,
            "peak": 1331228
        },
        "app": [
            {
                "appid": 1034052,
                "used": 439828,
                "peak": 651732,
                "modules": [
                    {
                        "file": "app",
                        "used": 277580,
                        "peak": 531292
                    },
                    {
                        "file": "pages/main_index",
                        "used": 162248,
                        "peak": 369564
                    }
                ]
            }
        ]
    },
    "description": "1:56:03 PM"
}
```


Watch app side
==== app.js =====

```javascript
import { getPerformance } from '@zos/app'
(...)

      setInterval(() => {
        const resp = getPerformance('memory');
        this.call({
          method: 'MEMORY',
          data: JSON.stringify(resp),
        });
      }, 1000);
      
      ! Don't forget to remove the interval in onDestroy
      clearInterval(intervalID)
```


To clear data,. You can start the server, but I know it is inconvenient, so you can run trigger API to clear the cache

API Method: DELETE
URL: http://localhost:8081/add-data/<YOUR_APP_NAME>
Body: null

example:

```javascript
    onCall(data) {
      if (data.method == 'MEMORY_DELETE') {
        fetch({
          url: 'http://localhost:8081/add-data/my_app',
          method: 'DELETE',
        });
      }
    },
```







A description can be attached to each package with a memory dump, as a string, to make it easier to determine the moment when we want to check what is happening with the memory.


```json
{
    "memory": {
        "system": {
            "used": 1546320,
            "total": 3145728
        },
        "framework": {
            "used": 1036308,
            "peak": 1331228
        },
        "app": [
            {
                "appid": 1034052,
                "used": 439828,
                "peak": 651732,
                "modules": [
                    {
                        "file": "app",
                        "used": 277580,
                        "peak": 531292
                    },
                    {
                        "file": "pages/main_index",
                        "used": 162248,
                        "peak": 369564
                    }
                ]
            }
        ]
    },
    "description": "memory value after turning on function A"
}

```




### Build 
Require go 1.24+ 

https://go.dev/doc/install

The backend of giu depends on OpenGL 3.3


# MacOS
xcode-select --install

# Windows
Install mingw download here. Thanks @alchem1ster!
Add the binaries folder of mingw to the path (usually is \mingw64\bin).

# Linux
sudo apt install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev libglx-dev libgl1-mesa-dev libxxf86vm-dev

Redhat:

sudo dnf install libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libGL-devel libXxf86vm-devel

you may also need to install C/C++ compiler (like g++) if it isn't already installed. Follow go compiler prompts.

sudo dpkg --add-architecture arm64
sudo apt update
sudo apt install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu \
    libx11-dev:arm64 libxcursor-dev:arm64 libxrandr-dev:arm64 libxinerama-dev:arm64 libxi-dev:arm64 libglx-dev:arm64 libgl1-mesa-dev:arm64 libxxf86vm-dev:arm64
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ HOST=aarch64-linux-gnu go build -v



