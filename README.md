# Supermetrics-test

This application is Supermetrics homework resolvement. Application exposes /users endpoint which returns userlist.

Currently application uses simple key-method for authentication. 

Lot of the scripts are quickly done and are missing prober automation. Due to the lack of time these features are not developed.

<br>

## Development tools

This project is developed on Linux (Ubuntu wsl) and is guaranteed to work on Linux-system. If using Mac or Windows then
the application itself should work fine but some automation might require os specifig development.

Prefered IDE is Vscode as it has extension to support (Rest Client by Huachao Mao) [test-file](scripts/test.http file). Also curl 
command is supported for manual testing.
Developer is free to use anyother rest client for testing but they are not automated by the project.

### Folders

Application uses very simple folder-structure. 

`internals` - Application logic<br>
`tools` - Developer scripts etc.

<br>

## Development

Initialize local environment by running [init_localenv.go](./tools/init_localenv.go) (run in tools-path).


### To run application directly 

You can use script [run.sh](./tools/run.sh) but https://github.com/air-verse/air is required.
 
Run commands:<br> 
`go mod download`<br>
`go run .`

Now the server should be running. You can now check [Testing](#Testing) process to verify status.
<br>

### To run application on docker

You can use automated command: [run_docker.sh](./tools/run_docker.sh)<br>

### To run application on cloud/on-premises

To run application on non-local environment then add this project to CI/CD. Notice that secrets (ex. API_SECRET) are assumed to be handled by pipeline.
On kubernetes you can use kubernetes secrets or services like hashivault.

Dockerfile shows required envs in ARG-list but these can be passed to application in any other ways as sre sees best. Application only assumes to find these values as environment variables.

Do not use port which requires root-privileges as Dockerfile uses non-root user.


## Testing

Testing manually is still clumsy and lacking automation but because of the time these features are not developed further. 

### Manual test

Example http request is located on [here](./tools/test.http). On vscode install Rest Client extension and right click file content. Press 'Send Request' on drop-down menu.

Curl test is located [here](./tools/curl.test.sh). It uses test.http file for input params so verify that token is updated.

If you get error 401/403 you need to update Authorization token. To update token run [createTestToken](./tools/createTestToken.go) which uses .env file as source for the input.
Change this to new param for test.http. 
Notice that script must be runned in the root-directory with .env file described above.


### Unit tests

Because project has limited time therefore unit test do not cover 100% of the project. Tests are made only around the important business logic or security.

run command: `go test ./...`<br>
Or if you want to print debug messages: `go test -v ./...`