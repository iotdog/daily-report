function send_post_request(req_url, request_body, callback) {
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

function send_get_request(req_url, callback) {
  $.ajax({
    url: req_url,
    type: "get",
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

var BASE_URL = "http://localhost:1024"

function add_worker(worker_name, worker_num, email, dept, group, callback) {
  var req_url = BASE_URL + "/api/1.0/add_worker"
  var request_body = {
    "name": worker_name,
    "number": worker_num,
    "email": email,
    "dept": dept,
    "group": group,
  }
  send_post_request(req_url, request_body, callback)
}

function submit_report(worker_name, mainline, subline, plan, callback) {
  var req_url = BASE_URL + "/api/1.0/submit_report"
  var request_body = {
    "name": worker_name,
    "mainLine": mainline,
    "subLine": subline,
    "plan": plan,
  }
  send_post_request(req_url, request_body, callback)
}

function get_daily_report(callback) {
  var req_url = BASE_URL + "/api/1.0/get_daily_report"
  send_get_request(req_url, callback)
}
