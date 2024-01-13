# echo
echo framework project template 
## usage
### install tools
**gonew**
```bash
 go install golang.org/x/tools/cmd/gonew@latest
```
**wire**
```bash
go install github.com/google/wire/cmd/wire@latest
```
### new project
use command
```bash
gonew github.com/tpl-x/echo example.com/foo
```
> if you modified wire code, your should use wire to generate code for later usage.
