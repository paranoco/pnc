constant "listen_addr" {
    set = "127.0.0.1:8000"
}

constant "server" {
    set = "server"
}

service "http" "web_proxy" {
    listen_addr = "${const.listen_addr}"

    process "main" {
        command = ["/usr/local/bin/awesome-app", "${const.server}"]
    }

    process "mgmt" {
        command = ["/usr/local/bin/awesome-app", "mgmt"]
    }
}