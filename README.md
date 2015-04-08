# hipchat-message
Small utility that allows you to send HipChat messages from terminal

hipchat-message reads message from standard input, so you can either use pipe
```
$ uname -a | hipchat-message @mkevac
```
or use Ctrl-D to finish writing
```
$ hipchat-message @mkevac
This is a test message
[Ctrl-D]
```

hipchat-message accepts recepients in one of the forms that you can see when running without any parameters
```
$ hipchat-message
Please enter receiver as a parameter in one of the forms:
  1. Name
  2. @username
  3. +room
```

Use `-c` to send message as a code:
```
$ hipchat-message -h
Usage of hipchat-message:
  -c=false: send message as a code
  -n=false: create new config file
```
