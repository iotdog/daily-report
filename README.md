# 1. Introduction
A management tool for daily report

# 2. How to use it

* run mongodb from docker

```
mkdir db
sudo docker run --rm --name mongodb -v $(pwd)/db:/data/db -p 27017:27017 -d mongo
```

* run server app

```
go run server.go
```

* add worker

visit http://localhost:1024/html/worker.html in browser.

* submit report

visit http://localhost:1024/html/report.html in browser.

* view today's reports list

visit http://localhost:1024/html/report_view.html in browser.
