function send_request(req_url, request_body, callback) {
  $.ajax({
    url: req_url,
    type: "post",
    data: JSON.stringify(request_body),
    async: false,
    dataType: "json",
    crossDomain: true,
    contentType: "application/json",
    success: function(data) {
      callback(data)
    },
    error: function() {
      callback("error")
    }
  })
}

function add_worker(worker_name, worker_num, dept, group, callback) {
  var req_url = "http://localhost:1024/api/1.0/add_worker"
  var request_body = {
    "name": worker_name,
    "number": worker_num,
    "dept": dept,
    "group": group,
  }
  send_request(req_url, request_body, callback)
}

function submit_report(worker_name, mainline, subline, plan, callback) {
  var req_url = "http://localhost:1024/api/1.0/submit_report"
  var request_body = {
    "name": worker_name,
    "mainLine": mainline,
    "subLine": subline,
    "plan": plan,
  }
  send_request(req_url, request_body, callback)
}
